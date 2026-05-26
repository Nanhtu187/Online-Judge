package service

import (
	"context"
	"errors"
	"time"

	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/models"
	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/repository"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/iam"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"golang.org/x/crypto/bcrypt"
)

type IamService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	GetUserInfo(ctx context.Context, userID string) (*pb.GetUserInfoResponse, error)
	VerifyToken(ctx context.Context, token string) (*pb.VerifyTokenResponse, error)
	CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error)
}

type iamService struct {
	repo      repository.IamRepository
	provider  database.IProvider
	logger    *zap.Logger
	jwtSecret []byte
}

func NewIamService(repo repository.IamRepository, provider database.IProvider, logger *zap.Logger, jwtSecret string) IamService {
	return &iamService{repo: repo, provider: provider, logger: logger, jwtSecret: []byte(jwtSecret)}
}

func (s *iamService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
	}

	err = s.provider.Transact(ctx, func(ctx context.Context) error {
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return err
		}

		// Assign default STUDENT role globally
		role, err := s.repo.GetRoleByName(ctx, "STUDENT")
		if err != nil {
			return err
		}

		return s.repo.AssignRole(ctx, &models.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
			Scope:  models.ScopeGlobal,
		})
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{UserId: user.ID}, nil
}

func (s *iamService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx = s.provider.Readonly(ctx)
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("user not found", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn("password mismatch", zap.String("email", req.Email))
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *iamService) GetUserInfo(ctx context.Context, userID string) (*pb.GetUserInfoResponse, error) {
	ctx = s.provider.Readonly(ctx)
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	perms, err := s.repo.GetGlobalPermissions(ctx, userID)
	if err != nil {
		s.logger.Error("failed to fetch global permissions", zap.String("user_id", userID), zap.Error(err))
	}

	return &pb.GetUserInfoResponse{
		User: &pb.User{
			Id:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Permissions: perms,
		},
	}, nil
}

func (s *iamService) VerifyToken(ctx context.Context, tokenString string) (*pb.VerifyTokenResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id not found in token")
	}

	return &pb.VerifyTokenResponse{UserId: userID}, nil
}

func (s *iamService) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	ctx = s.provider.Readonly(ctx)
	allowed, err := s.repo.CheckPermission(ctx, req.UserId, req.Permission, models.Scope(req.Scope.String()[len("SCOPE_"):]), req.ScopeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CheckPermissionResponse{Allowed: allowed}, nil
}
