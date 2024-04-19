package repo

import (
	"context"
	"github.com/Nanhtu187/Online-Judge/app/iam/model"
	"gorm.io/gorm"
)

type IUserPasswordRepo interface {
	GetUserPassword(ctx context.Context, userId int) (model.UserPassword, error)
	UpsertUserPassword(ctx context.Context, user model.UserPassword) error
}

type userPasswordRepo struct {
	db *gorm.DB
}

func (r *userPasswordRepo) GetUserPassword(ctx context.Context, userId int) (model.UserPassword, error) {
	var userPassword model.UserPassword
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&userPassword).Error
	return userPassword, err
}

func (r *userPasswordRepo) UpsertUserPassword(ctx context.Context, user model.UserPassword) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

func NewUserPasswordRepo(db *gorm.DB) IUserPasswordRepo {
	return &userPasswordRepo{
		db: db,
	}
}
