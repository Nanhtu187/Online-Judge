package repository

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/models"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"go.uber.org/zap"
)

type IamRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	AssignRole(ctx context.Context, userRole *models.UserRole) error
	
	CheckPermission(ctx context.Context, userID, permission string, scope models.Scope, scopeID string) (bool, error)
	GetGlobalPermissions(ctx context.Context, userID string) ([]string, error)
}

type iamRepository struct {
	provider *database.Provider
	logger   *zap.Logger
}

func NewIamRepository(provider *database.Provider, logger *zap.Logger) IamRepository {
	return &iamRepository{provider: provider, logger: logger}
}

func (r *iamRepository) CreateUser(ctx context.Context, user *models.User) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	return db.Create(user).Error
}

func (r *iamRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *iamRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *iamRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}
	if err := db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *iamRepository) AssignRole(ctx context.Context, userRole *models.UserRole) error {
	db, err := database.GetTx(ctx)
	if err != nil {
		return err
	}
	return db.Create(userRole).Error
}

func (r *iamRepository) CheckPermission(ctx context.Context, userID, permission string, scope models.Scope, scopeID string) (bool, error) {
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return false, err
	}

	// 1. Check via Roles
	var count int64
	query := db.Table("user_roles ur").
		Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("ur.user_id = ? AND p.name = ?", userID, permission)

	// Global roles apply everywhere
	// Scoped roles apply only to specific scope
	if scope == models.ScopeGlobal {
		query = query.Where("ur.scope = ?", models.ScopeGlobal)
	} else {
		query = query.Where("(ur.scope = ? OR (ur.scope = ? AND ur.scope_id = ?))", models.ScopeGlobal, models.ScopeGroup, scopeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	// 2. Check Direct Permissions
	query = db.Table("user_permissions up").
		Joins("JOIN permissions p ON up.permission_id = p.id").
		Where("up.user_id = ? AND p.name = ?", userID, permission)

	if scope == models.ScopeGlobal {
		query = query.Where("up.scope = ?", models.ScopeGlobal)
	} else {
		query = query.Where("(up.scope = ? OR (up.scope = ? AND up.scope_id = ?))", models.ScopeGlobal, models.ScopeGroup, scopeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *iamRepository) GetGlobalPermissions(ctx context.Context, userID string) ([]string, error) {
	db, err := database.GetReadonly(ctx)
	if err != nil {
		return nil, err
	}

	var perms []string

	// 1. From Roles
	err = db.Table("user_roles ur").
		Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("ur.user_id = ? AND ur.scope = ?", userID, models.ScopeGlobal).
		Pluck("p.name", &perms).Error
	if err != nil {
		return nil, err
	}

	// 2. From Direct Permissions
	var directPerms []string
	err = db.Table("user_permissions up").
		Joins("JOIN permissions p ON up.permission_id = p.id").
		Where("up.user_id = ? AND up.scope = ?", userID, models.ScopeGlobal).
		Pluck("p.name", &directPerms).Error
	if err != nil {
		return nil, err
	}

	perms = append(perms, directPerms...)
	return perms, nil
}
