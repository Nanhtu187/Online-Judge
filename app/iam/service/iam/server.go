package iam

import (
	"context"
	"github.com/Nanhtu187/Online-Judge/app/iam/config"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/logger"
	"github.com/Nanhtu187/Online-Judge/app/iam/repo"
	"github.com/Nanhtu187/Online-Judge/proto/rpc/iam/v1"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
	"strings"
)

type Server struct {
	iam.UnimplementedIamServiceServer
	service IService
}

func (s *Server) UpsertUser(ctx context.Context, request *iam.UpsertUserRequest) (*iam.UpsertUserResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}
	resp, err := s.service.UpsertUser(ctx, UpsertUserRequest{
		Username: request.Username,
		Name:     request.Name,
		School:   request.School,
		Class:    request.Class,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	return &iam.UpsertUserResponse{
		Code:    200,
		Message: "success",
		Data: &iam.CreateUserData{
			UserId: int32(resp.UserId),
		},
	}, nil

}

func (s *Server) GetUser(ctx context.Context, request *iam.GetUserRequest) (*iam.GetUserResponse, error) {
	req, err := vaildateGetUserRequest(request)
	if err != nil {
		return nil, err
	}
	resp, err := s.service.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &iam.GetUserResponse{
		Code:    200,
		Message: "success",
		Data: &iam.UserData{
			UserId: int32(resp.UserId),
			Name:   resp.Name,
			School: resp.School,
			Class:  resp.Class,
		},
	}, nil
}

func (s *Server) Login(ctx context.Context, request *iam.LoginRequest) (*iam.LoginResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	resp, err := s.service.Login(ctx, LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	return &iam.LoginResponse{
		Code:    200,
		Message: "success",
		Data: &iam.LoginData{
			AccessToken: resp.AccessToken,
		},
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, request *iam.RefreshTokenRequest) (*iam.RefreshTokenResponse, error) {
	token, err := getTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := s.service.RefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &iam.RefreshTokenResponse{
		Code:    200,
		Message: "success",
		Data: &iam.RefreshTokenData{
			AccessToken: resp.AccessToken,
		},
	}, nil

}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Extract(ctx).Error("error when get metadata from context")
		return "", ErrInvalidToken
	}

	const prefix = "Bearer "

	list := md.Get("authorization")
	if len(list) == 0 {
		return "", ErrInvalidToken
	}

	s := list[0]
	if !strings.HasPrefix(s, prefix) {
		return "", ErrInvalidToken
	}
	s = strings.TrimPrefix(s, prefix)
	return s, nil
}

func InitServer(db *gorm.DB, conf *config.Config) *Server {
	userRepo := repo.NewUserRepo(db)
	userPasswordRepo := repo.NewUserPasswordRepo(db)
	encryptService := NewEncryptoService(conf.PasswordEncryptKey)
	jwtService := NewTokenService([]byte(conf.TokenEncryptKey))
	svc := NewService(userRepo, userPasswordRepo, encryptService, jwtService)
	return NewServer(svc)
}

func NewServer(svc IService) *Server {
	return &Server{service: svc}
}
