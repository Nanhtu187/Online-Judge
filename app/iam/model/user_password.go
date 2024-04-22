package model

import "gorm.io/gorm"

type UserPassword struct {
	gorm.Model
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
