package server

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/proto/rpc/judger_server"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	judger_server.UnimplementedJudgerServiceServer
	service IService
}

func (s *Server) UpsertProblem(ctx context.Context, request *judger_server.UpsertProblemRequest) (*judger_server.UpsertProblemResponse, error) {
	req, err := validateUpsertProblemRequest(request)
	if err != nil {
		return nil, err
	}
	id, err := s.service.UpsertProblem(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.UpsertProblemResponse{
		Code:    200,
		Message: "success",
		Data: &judger_server.ProblemResponseData{
			ProblemId: int32(id),
		},
	}, nil
}

func NewServer(db *gorm.DB, rdb *redis.Client) *Server {
	return &Server{
		service: NewService(db, rdb),
	}
}
