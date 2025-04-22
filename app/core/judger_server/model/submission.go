package model

import (
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	ProblemID uint   `json:"problem_id" gorm:"not null"`
	ContestID uint   `json:"contest_id" gorm:"not null"`
	UserID    uint   `json:"user_id" gorm:"not null"`
	Language  string `json:"language" gorm:"type:varchar(50);not null"`
	Code      string `json:"code" gorm:"type:text;not null"`
	Status    string `json:"status" gorm:"type:varchar(50);not null"`
	Score     int    `json:"score" gorm:"not null"`
	Time      int    `json:"time" gorm:"not null"`
	Memory    int    `json:"memory" gorm:"not null"`
}
