package model

import "gorm.io/gorm"

type Contest struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	StartTime   string `json:"start_time" gorm:"type:timestamp;not null"`
	EndTime     string `json:"end_time" gorm:"type:timestamp;not null"`
}
