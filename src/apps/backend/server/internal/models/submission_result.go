package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestCaseResult struct {
	ID           string           `gorm:"type:char(36);primaryKey" json:"id"`
	SubmissionID string           `gorm:"type:char(36);not null;index" json:"submission_id"`
	TestCaseID   string           `gorm:"type:char(36);not null;index" json:"test_case_id"`
	Status       SubmissionStatus `gorm:"type:varchar(20);not null" json:"status"`
	ActualOutput string           `gorm:"type:text" json:"actual_output"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (t *TestCaseResult) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}
