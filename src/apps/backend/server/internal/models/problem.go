package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Problem struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	InputFormat  string         `gorm:"type:text" json:"input_format"`
	OutputFormat string         `gorm:"type:text" json:"output_format"`
	TimeLimit    int32          `gorm:"default:1000" json:"time_limit"` // in ms
	MemoryLimit  int32          `gorm:"default:256" json:"memory_limit"` // in MB
	TestCases    []TestCase     `gorm:"foreignKey:ProblemID" json:"test_cases"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Problem) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

type TestCase struct {
	ID             string         `gorm:"type:char(36);primaryKey" json:"id"`
	ProblemID      string         `gorm:"type:char(36);not null;index" json:"problem_id"`
	Input          string         `gorm:"type:text" json:"input"`
	ExpectedOutput string         `gorm:"type:text" json:"expected_output"`
	IsSample       bool           `gorm:"default:false" json:"is_sample"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (tc *TestCase) BeforeCreate(tx *gorm.DB) (err error) {
	if tc.ID == "" {
		tc.ID = uuid.New().String()
	}
	return
}
