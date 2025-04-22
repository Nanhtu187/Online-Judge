package repo

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"gorm.io/gorm"
)

type IRepository interface {
	CreateSubmission(ctx context.Context, submission *model.Submission) error
	GetSubmissionByID(ctx context.Context, id int32) (*model.Submission, error)
	CreateContest(ctx context.Context, contest *model.Contest) error
	GetContestByID(ctx context.Context, id int32) (*model.Contest, error)
	GetAllContests(ctx context.Context) ([]*model.Contest, error)
	GetAllSubmissions(ctx context.Context) ([]*model.Submission, error)
	GetResultByID(ctx context.Context, id int32) (*model.Result, error)
	GetAllResults(ctx context.Context) ([]*model.Result, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) CreateSubmission(ctx context.Context, submission *model.Submission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *repository) GetSubmissionByID(ctx context.Context, id int32) (*model.Submission, error) {
	var submission model.Submission
	err := r.db.WithContext(ctx).First(&submission, id).Error
	return &submission, err
}

func (r *repository) CreateContest(ctx context.Context, contest *model.Contest) error {
	return r.db.WithContext(ctx).Create(contest).Error
}

func (r *repository) GetContestByID(ctx context.Context, id int32) (*model.Contest, error) {
	var contest model.Contest
	err := r.db.WithContext(ctx).First(&contest, id).Error
	return &contest, err
}

func (r *repository) GetAllContests(ctx context.Context) ([]*model.Contest, error) {
	var contests []*model.Contest
	err := r.db.WithContext(ctx).Find(&contests).Error
	return contests, err
}

func (r *repository) GetAllSubmissions(ctx context.Context) ([]*model.Submission, error) {
	var submissions []*model.Submission
	err := r.db.WithContext(ctx).Find(&submissions).Error
	return submissions, err
}

func (r *repository) GetResultByID(ctx context.Context, id int32) (*model.Result, error) {
	var result model.Result
	err := r.db.WithContext(ctx).First(&result, id).Error
	return &result, err
}

func (r *repository) GetAllResults(ctx context.Context) ([]*model.Result, error) {
	var results []*model.Result
	err := r.db.WithContext(ctx).Find(&results).Error
	return results, err
}

func NewRepository(db *gorm.DB) IRepository {
	return &repository{
		db: db,
	}
}
