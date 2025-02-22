package models

import "gorm.io/gorm"

// DB Model
type Comment struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	PostID      uint      `gorm:"not null"`
	ContentText string    `gorm:"size:256;not null"`
}

func (Comment) TableName() string {
	return "comment"
}

// Request Models
type CreateCommentRequest struct {
	ContentText string `json:"content_text"`
}

// Response Models