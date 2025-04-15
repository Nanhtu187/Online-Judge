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
}

type service struct {
	db  *gorm.DB
	rdb *redis.Client
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
	err := s.db.Create(&submission).Error
	return int(submission.ID), err
}

func NewService(db *gorm.DB, rdb *redis.Client) IService {
	return &service{
		db:  db,
		rdb: rdb,
	}
}
