package iam

import (
	"context"
	"errors"
	"github.com/Nanhtu187/Online-Judge/app/iam/model"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/logger"
	"github.com/Nanhtu187/Online-Judge/app/iam/repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IService interface {
	UpsertUser(ctx context.Context, req UpsertUserRequest) (UpsertUserResponse, error)
	GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error)
	DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error)
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (RefreshTokenResponse, error)
	GetListUser(ctx context.Context, req GetListUserRequest) (GetListUserResponse, error)
}

// UpsertUser Upsert user info and password
func (s *service) UpsertUser(ctx context.Context, req UpsertUserRequest) (UpsertUserResponse, error) {
	userId, err := s.upsertUserInfo(ctx, req)
	if err != nil {
		return UpsertUserResponse{}, err
	}
	if req.Password != "" {
		err = s.upsertUserPassword(ctx, userId, req.Password)
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
		if validateCreateUserErr := validateCreateUserRequest(req); validateCreateUserErr != nil {
			return 0, validateCreateUserErr
		}

		user = model.User{
			Username: req.Username,
			Name:     req.Name,
			School:   req.School,
			Class:    req.Class,
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
func (s *service) upsertUserPassword(ctx context.Context, userId int, password string) error {
	password, err := s.encryptService.EncryptPassword(ctx, password)
	if err != nil {
		logger.Extract(ctx).Error("Error when encrypt password", zap.Error(err))
		return err
	}
	userPassword, err := s.userPasswordRepo.GetUserPassword(ctx, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userPassword = model.UserPassword{
			UserId:   userId,
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

func (s *service) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	return DeleteUserResponse{}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	return LoginResponse{}, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (RefreshTokenResponse, error) {
	return RefreshTokenResponse{}, nil
}

func (s *service) GetListUser(ctx context.Context, req GetListUserRequest) (GetListUserResponse, error) {
	return GetListUserResponse{}, nil
}

type service struct {
	userRepo         repo.IUserRepo
	userPasswordRepo repo.IUserPasswordRepo
	encryptService   ICryptographyService
}

func NewService(repo repo.IUserRepo, userPasswordRepo repo.IUserPasswordRepo, encryptService ICryptographyService) IService {
	return &service{
		userRepo:         repo,
		userPasswordRepo: userPasswordRepo,
		encryptService:   encryptService,
	}
}
