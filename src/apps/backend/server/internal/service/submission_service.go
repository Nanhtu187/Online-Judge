package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/repository"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"github.com/Nanhtu187/online-judge/src/packages/iam"
	"github.com/Nanhtu187/online-judge/src/packages/kafka"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
	kfk "github.com/segmentio/kafka-go"
)

type SubmissionService interface {
	SubmitCode(ctx context.Context, req *pb.SubmitCodeRequest) (*pb.SubmitCodeResponse, error)
	ListSubmissions(ctx context.Context, req *pb.ListSubmissionsRequest) (*pb.ListSubmissionsResponse, error)
	GetSubmissionResultDetail(ctx context.Context, req *pb.GetSubmissionResultRequest) (*pb.GetSubmissionResultDetailResponse, error)
}

type submissionService struct {
	repo       repository.SubmissionRepository
	resultRepo repository.SubmissionResultRepository
	provider   database.IProvider
	writer     *kfk.Writer
}

func NewSubmissionService(repo repository.SubmissionRepository, resultRepo repository.SubmissionResultRepository, provider database.IProvider, brokers []string) SubmissionService {
	writer := &kfk.Writer{
		Addr:     kfk.TCP(brokers...),
		Topic:    kafka.SubmissionRequestTopic,
		Balancer: &kfk.LeastBytes{},
	}
	return &submissionService{repo: repo, resultRepo: resultRepo, provider: provider, writer: writer}
}

func (s *submissionService) SubmitCode(ctx context.Context, req *pb.SubmitCodeRequest) (*pb.SubmitCodeResponse, error) {
	userID := iam.GetUserID(ctx)
	submission := &models.Submission{
		ProblemID:      req.ProblemId,
		CodeContent:    req.CodeContent,
		Status:         models.StatusPending,
		Language:       req.Language,
		SubmissionType: models.SubmissionType(pb.SubmissionType_name[int32(req.SubmissionType)][len("SUBMISSION_TYPE_"):]),
		UserID:         userID,
	}

	err := s.provider.Transact(ctx, func(ctx context.Context) error {
		return s.repo.Create(ctx, submission)
	})
	if err != nil {
		return nil, err
	}

	event := kafka.SubmissionEvent{
		SubmissionID:   submission.ID,
		ProblemID:      submission.ProblemID,
		CodeContent:    submission.CodeContent,
		Language:       submission.Language,
		SubmissionType: string(submission.SubmissionType),
	}
	eventJSON, _ := json.Marshal(event)

	err = s.writer.WriteMessages(ctx, kfk.Message{
		Key:   []byte(submission.ID),
		Value: eventJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to push submission event: %w", err)
	}

	return &pb.SubmitCodeResponse{SubmissionId: submission.ID}, nil
}

func (s *submissionService) ListSubmissions(ctx context.Context, req *pb.ListSubmissionsRequest) (*pb.ListSubmissionsResponse, error) {
	ctx = s.provider.Readonly(ctx)
	page := int(req.Page)
	if page == 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize == 0 {
		pageSize = 10
	}

	submissions, err := s.repo.List(ctx, req.ProblemId, req.UserId, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbSubmissions []*pb.Submission
	for _, sub := range submissions {
		pbSubmissions = append(pbSubmissions, &pb.Submission{
			Id:             sub.ID,
			ProblemId:      sub.ProblemID,
			CodeContent:    sub.CodeContent,
			Status:         pb.SubmissionStatus(pb.SubmissionStatus_value["SUBMISSION_STATUS_"+string(sub.Status)]),
			Language:       sub.Language,
			SubmissionType: pb.SubmissionType(pb.SubmissionType_value["SUBMISSION_TYPE_"+string(sub.SubmissionType)]),
			UserId:         sub.UserID,
			ProblemTitle:   sub.ProblemTitle,
			CreatedAt:      sub.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListSubmissionsResponse{Submissions: pbSubmissions}, nil
}

func (s *submissionService) GetSubmissionResultDetail(ctx context.Context, req *pb.GetSubmissionResultRequest) (*pb.GetSubmissionResultDetailResponse, error) {
	ctx = s.provider.Readonly(ctx)
	
	sub, err := s.repo.GetByID(ctx, req.SubmissionId)
	if err != nil {
		return nil, err
	}

	results, err := s.resultRepo.ListBySubmission(ctx, req.SubmissionId)
	if err != nil {
		return nil, err
	}

	var pbResults []*pb.TestCaseResult
	for _, r := range results {
		pbResults = append(pbResults, &pb.TestCaseResult{
			TestCaseId:   r.TestCaseID,
			Status:       pb.SubmissionStatus(pb.SubmissionStatus_value["SUBMISSION_STATUS_"+string(r.Status)]),
			ActualOutput: r.ActualOutput,
		})
	}

	return &pb.GetSubmissionResultDetailResponse{
		TestCaseResults: pbResults,
		Status:          pb.SubmissionStatus(pb.SubmissionStatus_value["SUBMISSION_STATUS_"+string(sub.Status)]),
	}, nil
}
