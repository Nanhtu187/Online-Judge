package submission_repo

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"gorm.io/gorm"
)

type ISubmissionRepo interface {
	UpsertSubmission(ctx context.Context, submissionId int, problemId, contestId, userId int, language, code, status string, score, time, memory int) (int, error)
	GetSubmission(ctx context.Context, submissionId int) (model.Submission, error)
	GetListSubmissionPreview(ctx context.Context, request GetListSubmissionRequest) ([]SubmissionPreview, error)
}

type submissionRepo struct {
	db *gorm.DB
}

func (s *submissionRepo) UpsertSubmission(ctx context.Context, submissionId int, problemId, contestId, userId int, language, code, status string, score, time, memory int) (int, error) {
	var submission model.Submission
	if submissionId != 0 {
		if err := s.db.First(&submission, submissionId).Error; err != nil {
			return 0, err
		}
		if problemId != 0 {
			submission.ProblemID = uint(problemId)
		}
		if contestId != 0 {
			submission.ContestID = uint(contestId)
		}
		if userId != 0 {
			submission.UserID = uint(userId)
		}
		if language != "" {
			submission.Language = language
		}
		if code != "" {
			submission.Code = code
		}
		if status != "" {
			submission.Status = status
		}
		if score != 0 {
			submission.Score = score
		}
		if time != 0 {
			submission.Time = time
		}
		if memory != 0 {
			submission.Memory = memory
		}

		if err := s.db.Save(&submission).Error; err != nil {
			return 0, err
		}

		return int(submission.ID), nil

	} else {
		submission.ProblemID = uint(problemId)
		submission.ContestID = uint(contestId)
		submission.UserID = uint(userId)
		submission.Language = language
		submission.Code = code
		submission.Status = status
		submission.Score = score
		submission.Time = time
		submission.Memory = memory
		if err := s.db.Create(&submission).Error; err != nil {
			return 0, err
		}
		return int(submission.ID), nil
	}
}

func (s *submissionRepo) GetSubmission(ctx context.Context, submissionId int) (model.Submission, error) {
	var submission model.Submission
	if err := s.db.First(&submission, submissionId).Error; err != nil {
		return model.Submission{}, err
	}

	return submission, nil
}

func (s *submissionRepo) GetListSubmissionPreview(ctx context.Context, request GetListSubmissionRequest) ([]SubmissionPreview, error) {
	var submissions []SubmissionPreview
	var query = s.db.Model(&model.Submission{})
	if request.ContestId != 0 {
		query = query.Where("contest_id = ?", request.ContestId)
	}

	if request.UserId != 0 {
		query = query.Where("user_id = ?", request.UserId)
	}

	if request.Keywords != "" {
		query = query.Where("code LIKE ?", "%"+request.Keywords+"%")
	}

	query = query.Offset(request.Offset).Limit(request.Limit).Take(&submissions)
	if err := query.Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func NewSubmissionRepo(db *gorm.DB) ISubmissionRepo {
	return &submissionRepo{
		db: db,
	}
}
