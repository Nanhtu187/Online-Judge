package result_repo

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"gorm.io/gorm"
)

type IResultRepo interface {
	UpsertResult(ctx context.Context, resultId int, submissionId uint, status string, testCaseId uint, score, time, memory int) (int, error)
	GetResult(ctx context.Context, resultId int) (model.Result, error)
	GetListResults(ctx context.Context, request GetListResultsRequest) ([]ResultPreview, error)
}

type resultRepo struct {
	db *gorm.DB
}

func (r *resultRepo) UpsertResult(ctx context.Context, resultId int, submissionId uint, status string, testCaseId uint, score, time, memory int) (int, error) {
	var result model.Result
	if resultId != 0 {
		if err := r.db.First(&result, resultId).Error; err != nil {
			return 0, err
		}
		if submissionId != 0 {
			result.SubmissionID = submissionId
		}
		if status != "" {
			result.Status = status
		}
		if testCaseId != 0 {
			result.TestCaseID = testCaseId
		}
		if score != 0 {
			result.Score = score
		}
		if time != 0 {
			result.Time = time
		}
		if memory != 0 {
			result.Memory = memory
		}

		if err := r.db.Save(&result).Error; err != nil {
			return 0, err
		}

		return int(result.ID), nil

	} else {
		result.SubmissionID = submissionId
		result.Status = status
		result.TestCaseID = testCaseId
		result.Score = score
		result.Time = time
		result.Memory = memory
		if err := r.db.Create(&result).Error; err != nil {
			return 0, err
		}
		return int(result.ID), nil
	}
}

func (r *resultRepo) GetResult(ctx context.Context, resultId int) (model.Result, error) {
	var result model.Result
	if err := r.db.First(&result, resultId).Error; err != nil {
		return model.Result{}, err
	}

	return result, nil
}

func (r *resultRepo) GetListResults(ctx context.Context, request GetListResultsRequest) ([]ResultPreview, error) {
	var results []ResultPreview
	var query = r.db.Model(&model.Result{})
	if request.SubmissionId != 0 {
		query = query.Where("submission_id = ?", request.SubmissionId)
	}

	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}

	if request.TestCaseId != 0 {
		query = query.Where("test_case_id = ?", request.TestCaseId)
	}

	query = query.Offset(request.Offset).Limit(request.Limit).Take(&results)
	if err := query.Error; err != nil {
		return nil, err
	}

	return results, nil
}

func NewResultRepo(db *gorm.DB) IResultRepo {
	return &resultRepo{
		db: db,
	}
}
