package iam

import (
	"context"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ICryptographyService interface {
	EncryptPassword(ctx context.Context, password string) (string, error)
	ComparePassword(ctx context.Context, password, hashedPassword string) error
}

type cryptographyService struct {
	secretKey int
}

func (s *cryptographyService) EncryptPassword(ctx context.Context, password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), s.secretKey)
	if err != nil {
		logger.Extract(ctx).Error("Error when encrypt password", zap.Error(err))
		return "", err
	}

	return string(hashed), nil
}

func (s *cryptographyService) ComparePassword(ctx context.Context, password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.Extract(ctx).Error("Error when compare password", zap.Error(err))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrPasswordNotMatch
		}
		return err
	}

	return nil

}

func NewEncryptoService(secretKey int) ICryptographyService {
	return &cryptographyService{
		secretKey: secretKey,
	}
}
