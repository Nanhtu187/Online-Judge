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

func (s *Server) UpsertContest(ctx context.Context, request *judger_server.UpsertContestRequest) (*judger_server.UpsertContestResponse, error) {
	req, err := validateUpsertContestRequest(request)
	if err != nil {
		return nil, err
	}
	id, err := s.service.UpsertContest(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.UpsertContestResponse{
		Code:    200,
		Message: "success",
		Contest: &judger_server.ContestResponseData{
			ContestId: int32(id),
		},
	}, nil
}

func (s *Server) GetContest(ctx context.Context, request *judger_server.GetContestRequest) (*judger_server.GetContestResponse, error) {
	req, err := validateGetContestRequest(request)
	if err != nil {
		return nil, err
	}
	contest, err := s.service.GetContest(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetContestResponse{
		Code:    200,
		Message: "success",
		Data:    contest,
	}, nil
}

func (s *Server) GetListContests(ctx context.Context, request *judger_server.GetListContestsRequest) (*judger_server.GetListContestsResponse, error) {
	req, err := validateGetListContestsRequest(request)
	if err != nil {
		return nil, err
	}
	contests, err := s.service.GetListContests(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetListContestsResponse{
		Code:    200,
		Message: "success",
		Data: &judger_server.GetListContestsResponseData{
			Contests: contests,
		},
	}, nil
}

func (s *Server) GetSubmission(ctx context.Context, request *judger_server.GetSubmissionRequest) (*judger_server.GetSubmissionResponse, error) {
	req, err := validateGetSubmissionRequest(request)
	if err != nil {
		return nil, err
	}
	submission, err := s.service.GetSubmission(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetSubmissionResponse{
		Code:    200,
		Message: "success",
		Data:    submission,
	}, nil
}

func (s *Server) GetListSubmissions(ctx context.Context, request *judger_server.GetListSubmissionsRequest) (*judger_server.GetListSubmissionsResponse, error) {
	req, err := validateGetListSubmissionsRequest(request)
	if err != nil {
		return nil, err
	}
	submissions, err := s.service.GetListSubmissions(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetListSubmissionsResponse{
		Code:    200,
		Message: "success",
		Data: &judger_server.GetListSubmissionsResponseData{
			Submissions: submissions,
		},
	}, nil
}

func (s *Server) GetResult(ctx context.Context, request *judger_server.GetResultRequest) (*judger_server.GetResultResponse, error) {
	req, err := validateGetResultRequest(request)
	if err != nil {
		return nil, err
	}
	result, err := s.service.GetResult(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetResultResponse{
		Code:    200,
		Message: "success",
		Data:    result,
	}, nil
}

func (s *Server) GetListResults(ctx context.Context, request *judger_server.GetListResultsRequest) (*judger_server.GetListResultsResponse, error) {
	req, err := validateGetListResultsRequest(request)
	if err != nil {
		return nil, err
	}
	results, err := s.service.GetListResults(ctx, req)
	if err != nil {
		return nil, err
	}
	return &judger_server.GetListResultsResponse{
		Code:    200,
		Message: "success",
		Data: &judger_server.GetListResultsResponseData{
			Results: results,
		},
	}, nil
}

func NewServer(db *gorm.DB, rdb *redis.Client) *Server {
	return &Server{
		service: NewService(db, rdb),
	}
}
