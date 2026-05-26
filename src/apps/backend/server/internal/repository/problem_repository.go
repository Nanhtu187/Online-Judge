package repository

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"go.uber.org/zap"
)

type ProblemRepository interface {
	CreateProblem(ctx context.Context, problem *models.Problem) error
	GetByID(ctx context.Context, id string) (*models.Problem, error)
	ListProblems(ctx context.Context, page, pageSize int) ([]*models.Problem, error)
	ListTestCases(ctx context.Context, problemID string, isSample *bool) ([]*models.TestCase, error)
	UpdateProblem(ctx context.Context, problem *models.Problem) error
	DeleteTestCases(ctx context.Context, problemID string) error
	CreateTestCase(ctx context.Context, testCase *models.TestCase) error

	GetTagByName(ctx context.Context, name string) (*models.Tag, error)
	CreateTag(ctx context.Context, tag *models.Tag) error
	ReplaceTags(ctx context.Context, problemID string, tags []*models.Tag) error
}

type problemRepository struct {
	provider *database.Provider
	logger   *zap.Logger
}

func NewProblemRepository(provider *database.Provider, logger *zap.Logger) ProblemRepository {
	return &problemRepository{provider: provider, logger: logger}
}

func (r *problemRepository) CreateProblem(ctx context.Context, problem *models.Problem) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(problem).Error; err != nil {
		r.logger.Error("failed to create problem", zap.Error(err))
		return err
	}
	return nil
}

func (r *problemRepository) GetByID(ctx context.Context, id string) (*models.Problem, error) {
	var problem models.Problem
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	err = db.Preload("Tags").First(&problem, "id = ?", id).Error
	if err != nil {
		r.logger.Error("failed to get problem by id", zap.String("id", id), zap.Error(err))
	}
	return &problem, err
}

func (r *problemRepository) ListProblems(ctx context.Context, page, pageSize int) ([]*models.Problem, error) {
	var problems []*models.Problem
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	err = db.Preload("Tags").Offset((page - 1) * pageSize).Limit(pageSize).Find(&problems).Error
	if err != nil {
		r.logger.Error("failed to list problems", zap.Error(err))
	}
	return problems, err
}

func (r *problemRepository) ListTestCases(ctx context.Context, problemID string, isSample *bool) ([]*models.TestCase, error) {
	var testCases []*models.TestCase
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	
	query := db.Where("problem_id = ?", problemID)
	if isSample != nil {
		query = query.Where("is_sample = ?", *isSample)
	}

	err = query.Find(&testCases).Error
	if err != nil {
		r.logger.Error("failed to list test cases", zap.String("problem_id", problemID), zap.Error(err))
	}
	return testCases, err
}

func (r *problemRepository) CreateTestCase(ctx context.Context, testCase *models.TestCase) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(testCase).Error; err != nil {
		r.logger.Error("failed to create test case", zap.Error(err))
		return err
	}
	return nil
}

func (r *problemRepository) UpdateProblem(ctx context.Context, problem *models.Problem) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	if err := db.Model(&models.Problem{}).Where("id = ?", problem.ID).Updates(problem).Error; err != nil {
		r.logger.Error("failed to update problem", zap.Error(err))
		return err
	}
	return nil
}

func (r *problemRepository) DeleteTestCases(ctx context.Context, problemID string) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	return db.Where("problem_id = ?", problemID).Delete(&models.TestCase{}).Error
}

func (r *problemRepository) GetTagByName(ctx context.Context, name string) (*models.Tag, error) {
	var tag models.Tag
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	if err := db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *problemRepository) CreateTag(ctx context.Context, tag *models.Tag) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	return db.Create(tag).Error
}

func (r *problemRepository) ReplaceTags(ctx context.Context, problemID string, tags []*models.Tag) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	return db.Model(&models.Problem{ID: problemID}).Association("Tags").Replace(tags)
}
