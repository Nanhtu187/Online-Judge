package model

import "gorm.io/gorm"

type UserPassword struct {
	gorm.Model
	UserId   int    `json:"user_id"`
	Password string `json:"password"`
}
