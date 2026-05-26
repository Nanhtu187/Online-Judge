package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

type Role struct {
	ID   string `gorm:"type:char(36);primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
}

type Permission struct {
	ID   string `gorm:"type:char(36);primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
}

type Scope string

const (
	ScopeGlobal Scope = "GLOBAL"
	ScopeGroup  Scope = "GROUP"
)

type UserRole struct {
	ID      string `gorm:"type:char(36);primaryKey" json:"id"`
	UserID  string `gorm:"type:char(36);not null;index" json:"user_id"`
	RoleID  string `gorm:"type:char(36);not null" json:"role_id"`
	Scope   Scope  `gorm:"type:varchar(20);not null;default:'GLOBAL'" json:"scope"`
	ScopeID string `gorm:"type:char(36)" json:"scope_id"`
}

func (ur *UserRole) BeforeCreate(tx *gorm.DB) (err error) {
	if ur.ID == "" {
		ur.ID = uuid.New().String()
	}
	return
}

type UserPermission struct {
	ID           string `gorm:"type:char(36);primaryKey" json:"id"`
	UserID       string `gorm:"type:char(36);not null;index" json:"user_id"`
	PermissionID string `gorm:"type:char(36);not null" json:"permission_id"`
	Scope        Scope  `gorm:"type:varchar(20);not null;default:'GLOBAL'" json:"scope"`
	ScopeID      string `gorm:"type:char(36)" json:"scope_id"`
}

func (up *UserPermission) BeforeCreate(tx *gorm.DB) (err error) {
	if up.ID == "" {
		up.ID = uuid.New().String()
	}
	return
}
