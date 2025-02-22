package models

import (
	"time"

	"gorm.io/gorm"
)

// DB Model
type Post struct {
    gorm.Model
    ContentText     string    `gorm:"size:500"`
    ContentImagePath string   `gorm:"size:256"`
	UserID          uint      `gorm:"not null"`
    Visible         bool      `gorm:"not null"`
    Comments        []*Comment `gorm:"foreignKey:PostID"`
    LikedUsers      []*User   `gorm:"many2many:like;foreignKey:id;joinForeignKey:post_id;References:id;joinReferences:user_id"`
}

func (Post) TableName() string {
    return "post"
}

// Request Models
type CreatePostRequest struct {
	UserID uint `json:"user_id"`
	ContentText      string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible          bool `json:"visible"`
}

type EditPostRequest struct {
	ContentText      *string `json:"content_text"`
	ContentImagePath *string `json:"content_image_path"`
	Visible          *bool `json:"visible"`
}

// Response Models
type PostDetailResponse struct {
	PostID uint `json:"post_id"`
	UserID uint `json:"user_id"`
	ContentText string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible bool `json:"visible"`
	CreatedAt time.Time `json:"created_at"`
}