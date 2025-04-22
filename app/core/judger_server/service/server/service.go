package server

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IService interface {
	UpsertProblem(ctx context.Context, request UpsertProblemRequest) (int, error)
	GetProblem(ctx context.Context, request GetProblemRequest) (ProblemDetail, error)
	GetListProblem(ctx context.Context, request GetListProblemRequest) ([]ProblemPreview, error)
	UpsertSubmissions(ctx context.Context, request UpsertSubmissionRequest) (int, error)
	UpsertContest(ctx context.Context, request UpsertContestRequest) (int, error)
	GetContest(ctx context.Context, request GetContestRequest) (Contest, error)
	GetListContests(ctx context.Context, request GetListContestsRequest) ([]Contest, error)
	GetSubmission(ctx context.Context, request GetSubmissionRequest) (Submission, error)
	GetListSubmissions(ctx context.Context, request GetListSubmissionsRequest) ([]Submission, error)
	GetResult(ctx context.Context, request GetResultRequest) (Result, error)
	GetListResults(ctx context.Context, request GetListResultsRequest) ([]Result, error)
}

type service struct {
	db  *gorm.DB
	rdb *redis.Client
	repo IRepository
}

func (s *service) UpsertProblem(ctx context.Context, request UpsertProblemRequest) (int, error) {
	if request.ProblemId == 0 {
		s.db.Create(&model.Problem{
			Title:       request.Title,
			Description: request.Description,
		})
	} else {
		var problem model.Problem
		if err := s.db.First(&problem, request.ProblemId).Error; err != nil {
			return 0, err
		}
		if request.Title != "" {
			problem.Title = request.Title
		}
		if request.Description != "" {
			problem.Description = request.Description
		}
		s.db.Save(&problem)
	}
	return 0, nil
}

func (s *service) GetProblem(ctx context.Context, request GetProblemRequest) (ProblemDetail, error) {
	return ProblemDetail{}, nil
}

func (s *service) GetListProblem(ctx context.Context, request GetListProblemRequest) ([]ProblemPreview, error) {
	return nil, nil
}

func (s *service) UpsertSubmissions(ctx context.Context, request UpsertSubmissionRequest) (int, error) {
	submission := model.Submission{
		ProblemId: request.ProblemId,
		ContestId: request.ContestId,
		UserId:    request.UserId,
		Language:  request.Language,
		Code:      request.Code,
	}
	err := s.repo.CreateSubmission(ctx, &submission)
	return int(submission.ID), err
}

func (s *service) UpsertContest(ctx context.Context, request UpsertContestRequest) (int, error) {
	contest := model.Contest{
		ContestId:   request.ContestId,
		Name:        request.Name,
		Description: request.Description,
		StartTime:   request.StartTime,
		EndTime:     request.EndTime,
		ProblemIds:  request.ProblemIds,
	}
	err := s.repo.CreateContest(ctx, &contest)
	return int(contest.ID), err
}

func (s *service) GetContest(ctx context.Context, request GetContestRequest) (Contest, error) {
	contest, err := s.repo.GetContestByID(ctx, request.ContestId)
	if err != nil {
		return Contest{}, err
	}
	return Contest{
		ContestId:   contest.ContestId,
		Name:        contest.Name,
		Description: contest.Description,
		StartTime:   contest.StartTime,
		EndTime:     contest.EndTime,
		ProblemIds:  contest.ProblemIds,
	}, nil
}

func (s *service) GetListContests(ctx context.Context, request GetListContestsRequest) ([]Contest, error) {
	contests, err := s.repo.GetAllContests(ctx)
	if err != nil {
		return nil, err
	}
	var result []Contest
	for _, contest := range contests {
		result = append(result, Contest{
			ContestId:   contest.ContestId,
			Name:        contest.Name,
			Description: contest.Description,
			StartTime:   contest.StartTime,
			EndTime:     contest.EndTime,
			ProblemIds:  contest.ProblemIds,
		})
	}
	return result, nil
}

func (s *service) GetSubmission(ctx context.Context, request GetSubmissionRequest) (Submission, error) {
	submission, err := s.repo.GetSubmissionByID(ctx, request.SubmissionId)
	if err != nil {
		return Submission{}, err
	}
	return Submission{
		SubmissionId: submission.SubmissionId,
		ProblemId:    submission.ProblemId,
		Language:     submission.Language,
		Code:         submission.Code,
		Input:        submission.Input,
		Output:       submission.Output,
		TestCaseId:   submission.TestCaseId,
	}, nil
}

func (s *service) GetListSubmissions(ctx context.Context, request GetListSubmissionsRequest) ([]Submission, error) {
	submissions, err := s.repo.GetAllSubmissions(ctx)
	if err != nil {
		return nil, err
	}
	var result []Submission
	for _, submission := range submissions {
		result = append(result, Submission{
			SubmissionId: submission.SubmissionId,
			ProblemId:    submission.ProblemId,
			Language:     submission.Language,
			Code:         submission.Code,
			Input:        submission.Input,
			Output:       submission.Output,
			TestCaseId:   submission.TestCaseId,
		})
	}
	return result, nil
}

func (s *service) GetResult(ctx context.Context, request GetResultRequest) (Result, error) {
	result, err := s.repo.GetResultByID(ctx, request.ResultId)
	if err != nil {
		return Result{}, err
	}
	return Result{
		ResultId:    result.ResultId,
		ProblemId:   result.ProblemId,
		Language:    result.Language,
		Code:        result.Code,
		Input:       result.Input,
		Output:      result.Output,
		TestCaseId:  result.TestCaseId,
	}, nil
}

func (s *service) GetListResults(ctx context.Context, request GetListResultsRequest) ([]Result, error) {
	results, err := s.repo.GetAllResults(ctx)
	if err != nil {
		return nil, err
	}
	var resultList []Result
	for _, result := range results {
		resultList = append(resultList, Result{
			ResultId:    result.ResultId,
			ProblemId:   result.ProblemId,
			Language:    result.Language,
			Code:        result.Code,
			Input:       result.Input,
			Output:      result.Output,
			TestCaseId:  result.TestCaseId,
		})
	}
	return resultList, nil
}

func NewService(db *gorm.DB, rdb *redis.Client, repo IRepository) IService {
	return &service{
		db:  db,
		rdb: rdb,
		repo: repo,
	}
}
