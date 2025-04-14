package server

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IService interface {
	UpsertProblem(ctx context.Context, request UpsertProblemRequest) (int, error)
	UpsertSubmissions(ctx context.Context, request UpsertSubmissionRequest) (int, error)
}

type service struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (s *service) UpsertProblem(ctx context.Context, request UpsertProblemRequest) (int, error) {
	if request.id == 0 {
		s.db.Create(&model.Problem{
			Title:       request.title,
			Description: request.description,
		})
	} else {
		var problem model.Problem
		if err := s.db.First(&problem, request.id).Error; err != nil {
			return 0, err
		}
		if request.title != "" {
			problem.Title = request.title
		}
		if request.description != "" {
			problem.Description = request.description
		}
		s.db.Save(&problem)
	}
	return 0, nil
}

func (s *service) UpsertSubmissions(ctx context.Context, request UpsertSubmissionRequest) (int, error) {
	submission := model.Submission{
		ProblemId: request.problemId,
		ContestId: request.contestId,
		UserId:    request.userId,
		Language:  request.language,
		Code:      request.code,
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
