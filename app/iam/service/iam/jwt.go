package iam

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type ITokenService interface {
	GenerateToken(ctx context.Context, userName string) (string, error)
	ValidateToken(ctx context.Context, token string) error
	ExtractToken(ctx context.Context, token string) (string, int64, error)
}

type tokenService struct {
	secretKey []byte
}

func (s *tokenService) GenerateToken(ctx context.Context, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": userName,
			"exp":      time.Now().Add(time.Minute * 4).Unix(),
		})
	return token.SignedString(s.secretKey)
}

func (s *tokenService) ValidateToken(ctx context.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrInvalidToken
	}
	return nil
}

func (s *tokenService) ExtractToken(ctx context.Context, tokenString string) (string, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return "", 0, err

	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, ErrInvalidToken
	}
	return claims["username"].(string), int64(claims["exp"].(float64)), nil
}

func NewTokenService(secretKey []byte) ITokenService {
	return &tokenService{
		secretKey: secretKey,
	}
}
