package models

import "gorm.io/gorm"

// DB Model
type Comment struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`       // ID of the user who created the comment
	PostID      uint      `gorm:"not null"`       // ID of the post the comment belongs to
	ContentText string    `gorm:"size:256;not null"` // Content text of the comment
}

func (Comment) TableName() string {
	return "comment"
}

// Request Models
type CreateCommentRequest struct {
	ContentText string `json:"content_text"`
}

// Response Models