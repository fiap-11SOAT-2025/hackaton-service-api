package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoStatus string

const (
	StatusPending    VideoStatus = "PENDING"
	StatusProcessing VideoStatus = "PROCESSING"
	StatusDone       VideoStatus = "DONE"
	StatusError      VideoStatus = "ERROR"
)

type Video struct {
	ID           string         `gorm:"type:uuid;primary_key;" json:"id"`
	UserID       string         `gorm:"type:uuid;index;not null" json:"user_id"`
	FileName     string         `json:"file_name"`
	InputBucket  string         `json:"input_bucket"`
	InputKey     string         `json:"input_key"`
	OutputBucket string         `json:"output_bucket"`
	OutputKey    string         `json:"output_key"`
	Status       VideoStatus    `gorm:"index;default:'PENDING'" json:"status"`
	ErrorMessage string         `json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewVideo(userID, fileName, inputKey string) *Video {
	return &Video{
		ID:        uuid.New().String(),
		UserID:    userID,
		FileName:  fileName,
		InputKey:  inputKey,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}