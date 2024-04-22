package repo

import (
	"context"
	"github.com/Nanhtu187/Online-Judge/app/iam/model"
	"gorm.io/gorm"
)

type IUserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUserByUserId(ctx context.Context, userId int) (model.User, error)
	UpsertUser(ctx context.Context, user model.User) (int, error)
}

type userRepo struct {
	db *gorm.DB
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepo) GetUserByUserId(ctx context.Context, userId int) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	return user, err
}

func (r *userRepo) UpsertUser(ctx context.Context, user model.User) (int, error) {
	err := r.db.WithContext(ctx).Save(&user).Error
	return user.ID, err
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &userRepo{
		db: db,
	}
}
