package iam

import (
	"context"
	"errors"
	"github.com/Nanhtu187/Online-Judge/app/iam/model"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/logger"
	"github.com/Nanhtu187/Online-Judge/app/iam/repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type IService interface {
	UpsertUser(ctx context.Context, req UpsertUserRequest) (UpsertUserResponse, error)
	GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error)
	DeleteUser(ctx context.Context, req DeleteUserRequest) error
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
	RefreshToken(ctx context.Context, token string) (RefreshTokenResponse, error)
	GetListUser(ctx context.Context, req GetListUserRequest) (GetListUserResponse, error)
	ValidateTokenToDeleteUser(ctx context.Context, token string, userId int) error
}

// UpsertUser Upsert user info and password
func (s *service) UpsertUser(ctx context.Context, req UpsertUserRequest) (UpsertUserResponse, error) {
	userId, err := s.upsertUserInfo(ctx, req)
	if err != nil {
		return UpsertUserResponse{}, err
	}

	if req.Password != "" {
		err = s.upsertUserPassword(ctx, userId, req.Username, req.Password)
		if err != nil {
			logger.Extract(ctx).Error("Error when upsert user", zap.Error(err))
			return UpsertUserResponse{}, err
		}
	}

	return UpsertUserResponse{
		UserId: userId,
	}, nil
}

// Upsert user info to user table
func (s *service) upsertUserInfo(ctx context.Context, req UpsertUserRequest) (int, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if ok, err := s.userRepo.ExistedUser(ctx, req.Username); err != nil {
			logger.Extract(ctx).Error("Error when check existed user", zap.Error(err))
			return 0, err
		} else if ok {
			return 0, ErrUserExisted
		}

		if validateCreateUserErr := validateCreateUserRequest(req); validateCreateUserErr != nil {
			return 0, validateCreateUserErr
		}

		user = model.User{
			Username: req.Username,
			Name:     req.Name,
			School:   req.School,
			Class:    req.Class,
			Role:     model.UserRoleContestant,
		}
	} else if err == nil {
		if req.Name != "" {
			user.Name = req.Name
		}

		if req.School != "" {
			user.School = req.School
		}

		if req.Class != "" {
			user.Class = req.Class
		}
	} else {
		logger.Extract(ctx).Error("Error when get user", zap.Error(err))
		return 0, err
	}
	logger.Extract(ctx).Info("Upsert user", zap.Any("user", user))
	userId, err := s.userRepo.UpsertUser(ctx, user)
	return userId, nil
}

// Encrypt password and upsert to user_password table
func (s *service) upsertUserPassword(ctx context.Context, userId int, username string, password string) error {
	password, err := s.encryptService.EncryptPassword(ctx, password)
	if err != nil {
		logger.Extract(ctx).Error("Error when encrypt password", zap.Error(err))
		return err
	}
	userPassword, err := s.userPasswordRepo.GetUserPasswordByUserId(ctx, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userPassword = model.UserPassword{
			UserId:   userId,
			Username: username,
			Password: password,
		}
	} else if err == nil {
		userPassword.Password = password
	} else {
		logger.Extract(ctx).Error("Error when get user password", zap.Error(err))
		return err
	}
	if err = s.userPasswordRepo.UpsertUserPassword(ctx, userPassword); err != nil {
		logger.Extract(ctx).Error("Error when upsert user password", zap.Error(err))
		return err
	}
	return nil
}

func (s *service) GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	if req.UserId != 0 {
		user, err := s.userRepo.GetUserByUserId(ctx, req.UserId)
		if err != nil {
			logger.Extract(ctx).Error("Error when get user by user id", zap.Error(err))
			if err == gorm.ErrRecordNotFound {
				return GetUserResponse{}, ErrUserNotFound
			} else {
				return GetUserResponse{}, err
			}
		}
		return GetUserResponse{
			UserId: user.ID,
			Name:   user.Name,
			School: user.School,
			Class:  user.Class,
		}, nil
	}
	if req.Username != "" {
		user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			logger.Extract(ctx).Error("Error when get user by username", zap.Error(err))
			return GetUserResponse{}, err
		}
		return GetUserResponse{
			UserId: user.ID,
			Name:   user.Name,
			School: user.School,
			Class:  user.Class,
		}, nil

	}
	return GetUserResponse{}, nil
}

func (s *service) DeleteUser(ctx context.Context, req DeleteUserRequest) error {
	err := s.userRepo.DeleteUser(ctx, req.UserId)
	if err != nil {
		logger.Extract(ctx).Error("Error when delete user", zap.Error(err))
		return err
	}

	err = s.userPasswordRepo.DeleteUserPassword(ctx, req.UserId)
	if err != nil {
		logger.Extract(ctx).Error("Error when delete user password", zap.Error(err))
		return err
	}

	return nil
}

// Allow user to delete themselves or admin to delete other users
func (s *service) ValidateTokenToDeleteUser(ctx context.Context, token string, userId int) error {
	err := s.jwtService.ValidateToken(ctx, token)
	if err != nil {
		return err
	}

	currentUserId, _, err := s.jwtService.ExtractToken(ctx, token)
	if err != nil {
		return err
	}

	if currentUserId == userId {
		return nil
	}

	currentUser, err := s.userRepo.GetUserByUserId(ctx, currentUserId)
	if err != nil {
		return err
	}

	if currentUser.Role == model.UserRoleAdmin {
		return nil
	}

	return ErrForbiddenToDeleteUser

}

func (s *service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	userPassword, err := s.userPasswordRepo.GetUserPasswordByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return LoginResponse{}, ErrUserNotFound
		} else {
			return LoginResponse{}, err
		}
	}
	if err = s.encryptService.ComparePassword(ctx, req.Password, userPassword.Password); err != nil {
		return LoginResponse{}, err
	}
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return LoginResponse{}, err
	}
	accessToken, err := s.jwtService.GenerateToken(ctx, user.ID)
	if err != nil {
		logger.Extract(ctx).Error("Error when generate token", zap.Error(err))
		return LoginResponse{}, err
	}
	return LoginResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, token string) (RefreshTokenResponse, error) {
	username, exp, err := s.jwtService.ExtractToken(ctx, token)
	if err != nil {
		logger.Extract(ctx).Error("Error when validate token", zap.Error(err))
		return RefreshTokenResponse{}, err
	}
	if time.Unix(exp, 0).Before(time.Now().Add(time.Minute * 5)) {
		tokenString, err := s.jwtService.GenerateToken(ctx, username)
		if err != nil {
			logger.Extract(ctx).Error("Error when generate token", zap.Error(err))
		}
		return RefreshTokenResponse{
			AccessToken: tokenString,
		}, nil
	}
	return RefreshTokenResponse{AccessToken: token}, nil
}

func (s *service) GetListUser(ctx context.Context, req GetListUserRequest) (GetListUserResponse, error) {
	return GetListUserResponse{}, nil
}

type service struct {
	userRepo         repo.IUserRepo
	userPasswordRepo repo.IUserPasswordRepo
	encryptService   ICryptographyService
	jwtService       ITokenService
}

func NewService(repo repo.IUserRepo, userPasswordRepo repo.IUserPasswordRepo, encryptService ICryptographyService, jwtService ITokenService) IService {
	return &service{
		userRepo:         repo,
		userPasswordRepo: userPasswordRepo,
		encryptService:   encryptService,
		jwtService:       jwtService,
	}
}
