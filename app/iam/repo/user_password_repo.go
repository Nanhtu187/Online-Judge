package repo

import (
	"context"
	"github.com/Nanhtu187/Online-Judge/app/iam/model"
	"gorm.io/gorm"
)

type IUserPasswordRepo interface {
	GetUserPasswordByUserId(ctx context.Context, userId int) (model.UserPassword, error)
	GetUserPasswordByUsername(ctx context.Context, username string) (model.UserPassword, error)
	UpsertUserPassword(ctx context.Context, user model.UserPassword) error
}

type userPasswordRepo struct {
	db *gorm.DB
}

func (r *userPasswordRepo) GetUserPasswordByUserId(ctx context.Context, userId int) (model.UserPassword, error) {
	var userPassword model.UserPassword
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&userPassword).Error
	return userPassword, err
}

func (r *userPasswordRepo) GetUserPasswordByUsername(ctx context.Context, username string) (model.UserPassword, error) {
	var userPassword model.UserPassword
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&userPassword).Error
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
