package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionStatus string

type SubmissionType string

const (
	TypeTest     SubmissionType = "TEST"
	TypeOfficial SubmissionType = "OFFICIAL"
)

const (
	StatusPending       SubmissionStatus = "PENDING"
	StatusRunning       SubmissionStatus = "RUNNING"
	StatusAccepted      SubmissionStatus = "ACCEPTED"
	StatusFailed        SubmissionStatus = "FAILED"
	StatusWrongAnswer   SubmissionStatus = "WRONG_ANSWER"
	StatusRuntimeError  SubmissionStatus = "RUNTIME_ERROR"
	StatusLimitExceeded SubmissionStatus = "LIMIT_EXCEEDED"
	StatusCompileError  SubmissionStatus = "COMPILE_ERROR"
)

type Submission struct {
	ID             string           `gorm:"type:char(36);primaryKey" json:"id"`
	ProblemID      string           `gorm:"type:char(36);not null;index" json:"problem_id"`
	CodeContent    string           `gorm:"type:text;not null" json:"code_content"`
	Status         SubmissionStatus `gorm:"type:varchar(20);not null" json:"status"`
	Language       string           `gorm:"type:varchar(20);not null" json:"language"`
	SubmissionType SubmissionType   `gorm:"type:varchar(20);not null;default:'OFFICIAL'" json:"submission_type"`
	UserID         string           `gorm:"type:char(36);not null;index" json:"user_id"`
	ProblemTitle   string           `gorm:"->"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (s *Submission) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
