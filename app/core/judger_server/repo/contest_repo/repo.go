package contest_repo

import (
	"context"

	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/model"
	"gorm.io/gorm"
)

type IContestRepo interface {
	UpsertContest(ctx context.Context, contestId int, title, description string, startTime, endTime string) (int, error)
	GetContest(ctx context.Context, contestId int) (model.Contest, error)
	GetListContestPreview(ctx context.Context, request GetListContestRequest) ([]ContestPreview, error)
}

type contestRepo struct {
	db *gorm.DB
}

func (c *contestRepo) UpsertContest(ctx context.Context, contestId int, title, description string, startTime, endTime string) (int, error) {
	var contest model.Contest
	if contestId != 0 {
		if err := c.db.First(&contest, contestId).Error; err != nil {
			return 0, err
		}
		if title != "" {
			contest.Title = title
		}

		if description != "" {
			contest.Description = description
		}

		if startTime != "" {
			contest.StartTime = startTime
		}

		if endTime != "" {
			contest.EndTime = endTime
		}

		if err := c.db.Save(&contest).Error; err != nil {
			return 0, err
		}

		return int(contest.ID), nil

	} else {
		contest.Title = title
		contest.Description = description
		contest.StartTime = startTime
		contest.EndTime = endTime
		if err := c.db.Create(&contest).Error; err != nil {
			return 0, err
		}
		return int(contest.ID), nil
	}
}

func (c *contestRepo) GetContest(ctx context.Context, contestId int) (model.Contest, error) {
	var contest model.Contest
	if err := c.db.First(&contest, contestId).Error; err != nil {
		return model.Contest{}, err
	}

	return contest, nil
}

func (c *contestRepo) GetListContestPreview(ctx context.Context, request GetListContestRequest) ([]ContestPreview, error) {
	var contests []ContestPreview
	var query = c.db.Model(&model.Contest{})
	if request.CreatedBy != 0 {
		query = query.Where("created_by = ?", request.CreatedBy)
	}

	if request.Keywords != "" {
		query = query.Where("title LIKE ?", "%"+request.Keywords+"%")
	}

	query = query.Offset(request.Offset).Limit(request.Limit).Find(&contests)
	if err := query.Error; err != nil {
		return nil, err
	}

	return contests, nil
}

func NewContestRepo(db *gorm.DB) IContestRepo {
	return &contestRepo{
		db: db,
	}
}
