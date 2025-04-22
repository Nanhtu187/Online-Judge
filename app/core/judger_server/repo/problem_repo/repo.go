package problem_repo

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"gorm.io/gorm"
)

type IProblemRepo interface {
	UpsertProblem(ctx context.Context, problemId int, title, desciption string) (int, error)
	GetProblem(ctx context.Context, problemId int) (model.Problem, error)
	GetListProblemPreview(ctx context.Context, request GetListProblemRequest) ([]ProblemPreview, error)
}

type problemRepo struct {
	db *gorm.DB
}

func (p *problemRepo) UpsertProblem(ctx context.Context, problemId int, title, description string) (int, error) {
	var problem model.Problem
	if problemId != 0 {
		if err := p.db.First(&problem, problemId).Error; err != nil {
			return 0, err
		}
		if title != "" {
			problem.Title = title
		}

		if description != "" {
			problem.Description = description
		}

		if err := p.db.Save(&problem).Error; err != nil {
			return 0, err
		}

		return int(problem.ID), nil

	} else {
		problem.Title = title
		problem.Description = description
		if err := p.db.Create(&problem).Error; err != nil {
			return 0, err
		}
		return int(problem.ID), nil
	}
}

func (p *problemRepo) GetProblem(ctx context.Context, problemId int) (model.Problem, error) {
	var problem model.Problem
	if err := p.db.First(&problem, problemId).Error; err != nil {
		return model.Problem{}, err
	}

	return problem, nil
}
func (p *problemRepo) GetListProblemPreview(ctx context.Context, request GetListProblemRequest) ([]ProblemPreview, error) {
	var problems []ProblemPreview
	var query = p.db.Model(&model.Problem{})
	if request.ContestId != 0 {
		query = query.Where("contest_id = ?", request.ContestId)
	}

	if request.CreatedBy != 0 {
		query = query.Where("created_by = ?", request.CreatedBy)
	}

	if request.Keywords != "" {
		query = query.Where("title LIKE ?", "%"+request.Keywords+"%")
	}

	query = query.Offset(request.Offset).Limit(request.Limit).Take(&problems)
	if err := query.Error; err != nil {
		return nil, err
	}

	return problems, nil
}

func NewProblemRepo(db *gorm.DB) IProblemRepo {
	return &problemRepo{
		db: db,
	}
}
