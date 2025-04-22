package server

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/common/logger"
	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/repo/problem_repo"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IService interface {
	UpsertProblem(ctx context.Context, request UpsertProblemRequest) (int, error)
	GetProblem(ctx context.Context, request GetProblemRequest) (ProblemDetail, error)
	GetListProblem(ctx context.Context, request GetListProblemRequest) ([]ProblemPreview, error)
	UpsertContest(ctx context.Context, request UpsertContestRequest) (int, error)
	GetContest(ctx context.Context, request GetContestRequest) (ContestDetail, error)
	GetListContest(ctx context.Context, request GetListContestRequest) ([]ContestPreview, error)
	UpsertSubmissions(ctx context.Context, request UpsertSubmissionRequest) (int, error)
	GetSubmission(ctx context.Context, request GetSubmissionRequest) (SubmissionDetail, error)
	GetListSubmission(ctx context.Context, request GetListSubmissionRequest) ([]SubmissionPreview, error)
	GetListResult(ctx context.Context, request GetListResultRequest) ([]ResultPreview, error)
	GetResult(ctx context.Context, request GetResultRequest) (ResultDetail, error)
}

type service struct {
	problemRepo    problem_repo.IProblemRepo
	contestRepo    IContestRepo
	submissionRepo ISubmissionRepo
	resultRepo     IResultRepo
}

func (s *service) UpsertProblem(ctx context.Context, request UpsertProblemRequest) (int, error) {
	problemId, err := s.problemRepo.UpsertProblem(ctx, int(request.ProblemId), request.Title, request.Description)
	if err != nil {
		logger.Extract(ctx).Error("Get problem error: %v", zap.Error(err))
		return 0, err
	}

	return problemId, nil
}

func (s *service) GetProblem(ctx context.Context, request GetProblemRequest) (ProblemDetail, error) {
	problem, err := s.problemRepo.GetProblem(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get problem error: %v", zap.Error(err))
		return ProblemDetail{}, err
	}

	return ProblemDetail{
		ProblemId:   problem.ID,
		Title:       problem.Title,
		Description: problem.Description,
	}, nil
}

func (s *service) GetListProblem(ctx context.Context, request GetListProblemRequest) ([]ProblemPreview, error) {
	problems, err := s.problemRepo.GetListProblemPreview(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get list problem preview error: %v", zap.Error(err))
		return nil, err
	}

	var result []ProblemPreview
	for _, problem := range problems {
		result = append(result, ProblemPreview{
			ProblemId: problem.ID,
			Title:     problem.Title,
		})
	}
	return result, nil

}

func (s *service) UpsertContest(ctx context.Context, request UpsertContestRequest) (int, error) {
	problemId, err := s.problemRepo.UpsertProblem(ctx, request.Title)
	if err != nil {
		logger.Extract(ctx).Error("Get problem error: %v", zap.Error(err))
		return 0, err
	}

	return problemId, nil
}

func (s *service) GetContest(ctx context.Context, request GetContestRequest) (ContestDetail, error) {
	contest, err := s.contestRepo.GetContest(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get contest error: %v", zap.Error(err))
		return ContestDetail{}, err
	}

	return ContestDetail{
		ContestId:   contest.ID,
		Title:       contest.Title,
		Description: contest.Description,
	}, nil
}

func (s *service) GetListContest(ctx context.Context, request GetListContestRequest) ([]ContestPreview, error) {
	contests, err := s.contestRepo.GetListContestPreview(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get list contest preview error: %v", err)
		return nil, err
	}

	var result []ContestPreview
	for _, contest := range contests {
		result = append(result, ContestPreview{
			ContestId: contest.ID,
			Title:     contest.Title,
		})
	}
	return result, nil
}

func (s *service) GetSubmission(ctx context.Context, request GetSubmissionRequest) (SubmissionDetail, error) {
	submission, err := s.submissionRepo.GetSubmission(ctx)
	if err != nil {
		logger.Extract(ctx).Error("Get submission error: %v", err)
		return SubmissionDetail{}, err
	}

	return SubmissionDetail{
		SubmissionId: submission.ID,
		ProblemId:    submission.ProblemId,
		UserId:       submission.UserId,
	}, nil
}

func (s *service) GetListSubmission(ctx context.Context, request GetListSubmissionRequest) ([]SubmissionPreview, error) {
	submissions, err := s.submissionRepo.GetListSubmissionPreview(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get list submission preview error: %v", err)
		return nil, err
	}

	var result []SubmissionPreview
	for _, submission := range submissions {
		result = append(result, SubmissionPreview{
			SubmissionId: submission.ID,
			ProblemId:    submission.ProblemId,
			UserId:       submission.UserId,
		})
	}
	return result, nil
}

func (s *service) GetListResult(ctx context.Context, request GetListResultRequest) ([]ResultPreview, error) {
	results, err := s.resultRepo.GetListResultPreview(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get list result preview error: %v", err)
		return nil, err
	}

	var result []ResultPreview
	for _, resultItem := range results {
		result = append(result, ResultPreview{
			ResultId:     resultItem.ID,
			SubmissionId: resultItem.SubmissionId,
			Status:       resultItem.Status,
		})
	}
	return result, nil
}

func (s *service) GetResult(ctx context.Context, request GetResultRequest) (ResultDetail, error) {
	result, err := s.resultRepo.GetResult(ctx, request)
	if err != nil {
		logger.Extract(ctx).Error("Get result error: %v", err)
		return ResultDetail{}, err
	}

	return ResultDetail{
		ResultId:     result.ID,
		SubmissionId: result.SubmissionId,
		Status:       result.Status,
	}, nil
}

func (s *service) UpsertSubmissions(ctx context.Context, request UpsertSubmissionRequest) (int, error) {
	submission := model.Submission{
		ProblemId: request.ProblemId,
		ContestId: request.ContestId,
		UserId:    request.UserId,
		Language:  request.Language,
		Code:      request.Code,
	}
	err := s.db.Create(&submission).Error
	return int(submission.ID), err
}

func NewService(db *gorm.DB, rdb *redis.Client) IService {
	return &service{
		db:  db,
		rdb: rdb,
	}
}
