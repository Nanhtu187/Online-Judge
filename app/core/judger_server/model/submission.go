package model

import (
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	ProblemId int32  `json:"problem_id" gorm:"index"`
	ContestId int32  `json:"contest_id" gorm:"index"`
	UserId    int32  `json:"user_id" gorm:"index"`
	Language  string `json:"language_id"`
	Code      string `json:"code"`
}
