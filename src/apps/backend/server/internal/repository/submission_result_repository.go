package repository

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"go.uber.org/zap"
)

type SubmissionResultRepository interface {
	ListBySubmission(ctx context.Context, submissionID string) ([]*models.TestCaseResult, error)
	Create(ctx context.Context, result *models.TestCaseResult) error
}

type submissionResultRepository struct {
	provider *database.Provider
	logger   *zap.Logger
}

func NewSubmissionResultRepository(provider *database.Provider, logger *zap.Logger) SubmissionResultRepository {
	return &submissionResultRepository{provider: provider, logger: logger}
}

func (r *submissionResultRepository) ListBySubmission(ctx context.Context, submissionID string) ([]*models.TestCaseResult, error) {
	var results []*models.TestCaseResult
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	err = db.Where("submission_id = ?", submissionID).Find(&results).Error
	if err != nil {
		r.logger.Error("failed to list submission results", zap.String("submission_id", submissionID), zap.Error(err))
	}
	return results, err
}

func (r *submissionResultRepository) Create(ctx context.Context, result *models.TestCaseResult) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	err = db.Create(result).Error
	if err != nil {
		r.logger.Error("failed to create submission result", zap.Error(err))
	}
	return err
}
