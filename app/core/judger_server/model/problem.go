package model

import "gorm.io/gorm"

type Problem struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null;unique"`
	Description string `json:"description" gorm:"type:text;not null"`
}
