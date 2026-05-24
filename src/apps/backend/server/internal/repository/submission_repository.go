package repository

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"go.uber.org/zap"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) error
	List(ctx context.Context, problemID string, page, pageSize int) ([]*models.Submission, error)
	GetByID(ctx context.Context, id string) (*models.Submission, error)
	UpdateStatus(ctx context.Context, id string, status models.SubmissionStatus) error
}

type submissionRepository struct {
	provider *database.Provider
	logger   *zap.Logger
}

func NewSubmissionRepository(provider *database.Provider, logger *zap.Logger) SubmissionRepository {
	return &submissionRepository{provider: provider, logger: logger}
}

func (r *submissionRepository) Create(ctx context.Context, s *models.Submission) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	err = db.Create(s).Error
	if err != nil {
		r.logger.Error("failed to create submission", zap.Error(err))
	}
	return err
}

func (r *submissionRepository) List(ctx context.Context, problemID string, page, pageSize int) ([]*models.Submission, error) {
	var submissions []*models.Submission
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	query := db.Model(&models.Submission{})
	if problemID != "" {
		query = query.Where("problem_id = ?", problemID)
	}
	
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&submissions).Error
	if err != nil {
		r.logger.Error("failed to list submissions", zap.String("problem_id", problemID), zap.Error(err))
	}
	return submissions, err
}

func (r *submissionRepository) GetByID(ctx context.Context, id string) (*models.Submission, error) {
	var submission models.Submission
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	err = db.First(&submission, "id = ?", id).Error
	if err != nil {
		r.logger.Error("failed to get submission by id", zap.String("id", id), zap.Error(err))
	}
	return &submission, err
}

func (r *submissionRepository) UpdateStatus(ctx context.Context, id string, status models.SubmissionStatus) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	err = db.Model(&models.Submission{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		r.logger.Error("failed to update submission status", zap.String("id", id), zap.Error(err))
	}
	return err
}
