package iam

import (
	"context"
	"github.com/Nanhtu187/Online-Judge/app/iam/config"
	"github.com/Nanhtu187/Online-Judge/app/iam/repo"
	"github.com/Nanhtu187/Online-Judge/proto/rpc/iam/v1"
	"gorm.io/gorm"
)

type Server struct {
	iam.UnimplementedIamServiceServer
	service IService
}

func (s *Server) UpsertUser(ctx context.Context, request *iam.UpsertUserRequest) (*iam.UpsertUserResponse, error) {
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

func InitServer(db *gorm.DB, conf *config.Config) *Server {
	userRepo := repo.NewUserRepo(db)
	userPasswordRepo := repo.NewUserPasswordRepo(db)
	encryptService := NewEncryptoService(conf.PasswordEncryptKey)
	svc := NewService(userRepo, userPasswordRepo, encryptService)
	return NewServer(svc)
}

func NewServer(svc IService) *Server {
	return &Server{service: svc}
}
