package model

import "gorm.io/gorm"

type Result struct {
	gorm.Model
	SubmissionID uint   `json:"submission_id" gorm:"not null"`
	Status       string `json:"status" gorm:"type:varchar(50);not null"`
	TestCaseID   uint   `json:"test_case_id" gorm:"not null"`
	Score        int    `json:"score" gorm:"not null"`
	Time         int    `json:"time" gorm:"not null"`
	Memory       int    `json:"memory" gorm:"not null"`
}
