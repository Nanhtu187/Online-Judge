package handler

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/service"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
)

type OnlineJudgeHandler struct {
	pb.UnimplementedOnlineJudgeServiceServer
	problemSvc    service.ProblemService
	submissionSvc service.SubmissionService
}

func NewOnlineJudgeHandler(problemSvc service.ProblemService, submissionSvc service.SubmissionService) *OnlineJudgeHandler {
	return &OnlineJudgeHandler{
		problemSvc:    problemSvc,
		submissionSvc: submissionSvc,
	}
}

// Problem Methods
func (h *OnlineJudgeHandler) UpsertProblem(ctx context.Context, req *pb.UpsertProblemRequest) (*pb.UpsertProblemResponse, error) {
	return h.problemSvc.UpsertProblem(ctx, req)
}

func (h *OnlineJudgeHandler) UpsertTestCases(ctx context.Context, req *pb.UpsertTestCasesRequest) (*pb.UpsertTestCasesResponse, error) {
	return h.problemSvc.UpsertTestCases(ctx, req)
}

func (h *OnlineJudgeHandler) GetProblem(ctx context.Context, req *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	return h.problemSvc.GetProblem(ctx, req)
}

func (h *OnlineJudgeHandler) ListProblems(ctx context.Context, req *pb.ListProblemsRequest) (*pb.ListProblemsResponse, error) {
	return h.problemSvc.ListProblems(ctx, req)
}

func (h *OnlineJudgeHandler) ListTestCases(ctx context.Context, req *pb.ListTestCasesRequest) (*pb.ListTestCasesResponse, error) {
	return h.problemSvc.ListTestCases(ctx, req)
}

// Submission Methods
func (h *OnlineJudgeHandler) SubmitCode(ctx context.Context, req *pb.SubmitCodeRequest) (*pb.SubmitCodeResponse, error) {
	return h.submissionSvc.SubmitCode(ctx, req)
}

func (h *OnlineJudgeHandler) ListSubmissions(ctx context.Context, req *pb.ListSubmissionsRequest) (*pb.ListSubmissionsResponse, error) {
	return h.submissionSvc.ListSubmissions(ctx, req)
}

func (h *OnlineJudgeHandler) GetSubmissionResultDetail(ctx context.Context, req *pb.GetSubmissionResultRequest) (*pb.GetSubmissionResultDetailResponse, error) {
	return h.submissionSvc.GetSubmissionResultDetail(ctx, req)
}
