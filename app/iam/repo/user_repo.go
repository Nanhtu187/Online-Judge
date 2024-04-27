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
	DeleteUser(ctx context.Context, userId int) error
	ExistedUser(ctx context.Context, username string) (bool, error)
}

type userRepo struct {
	db *gorm.DB
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Table("users").Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepo) GetUserByUserId(ctx context.Context, userId int) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", userId).First(&user).Error
	return user, err
}

func (r *userRepo) UpsertUser(ctx context.Context, user model.User) (int, error) {
	err := r.db.WithContext(ctx).Table("users").Save(&user).Error
	return user.ID, err
}

func (r *userRepo) DeleteUser(ctx context.Context, userId int) error {
	return r.db.WithContext(ctx).Where("id = ?", userId).Delete(&model.User{}).Error
}

func (r *userRepo) GetDeletedUser(ctx context.Context, userId int) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Unscoped().Table("users").Where("id = ?", userId).First(&user).Error
	return user, err
}

func (r *userRepo) ExistedUser(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

var _ IUserRepo = &userRepo{}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &userRepo{
		db: db,
	}
}
