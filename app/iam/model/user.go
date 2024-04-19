package model

import "gorm.io/gorm"

type UserRole string

const (
	UserRoleAdmin         UserRole = "admin"
	UserRoleContestant    UserRole = "contestant"
	UserRoleProblemSetter UserRole = "problem_setter"
)

type User struct {
	gorm.Model
	ID       int      `json:"id" gorm:"primaryKey"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	School   string   `json:"school"`
	Class    string   `json:"class"`
	Role     UserRole `json:"role"`
}
