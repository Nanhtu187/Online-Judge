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
		Data: &judger_server.UpsertProblemResponseData{
			ProblemId: int32(id),
		},
	}, nil
}

func (s *Server) GetProblem(ctx context.Context, request *judger_server.GetProblemRequest) (*judger_server.GetProblemResponse, error) {
	req, err := validateGetProblemRequest(request)
	if err != nil {
		return nil, err
	}
	problem, err := s.service.GetProblem(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetProblemResponse{
		Code:    200,
		Message: "success",
		Data:    problem,
	}, nil
}

func (s *Server) GetProblems(ctx context.Context, request *judger_server.GetProblemsRequest) (*judger_server.GetProblemsResponse, error) {
	req, err := validateGetProblemsRequest(request)
	if err != nil {
		return nil, err
	}
	problems, err := s.service.GetProblems(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetProblemsResponse{
		Code:    200,
		Message: "success",
		Data: &judger_server.GetProblemsResponseData{
			Problems: problems,
		},
	}, nil
}

func (s *Server) UpsertSubmissions(ctx context.Context, request *judger_server.UpsertSubmissionRequest) (*judger_server.UpsertSubmissionResponse, error) {
	req, err := validateUpsertSubmissionRequest(request)
	if err != nil {
		return nil, err
	}
	id, err := s.service.UpsertSubmissions(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.UpsertSubmissionResponse{
		Code:    200,
		Message: "success",
		Data: &judger_server.UpsertSubmissionResponseData{
			SubmissionId: int32(id),
		},
	}, nil
}

func NewServer(db *gorm.DB, rdb *redis.Client) *Server {
	return &Server{
		service: NewService(db, rdb),
	}
}
