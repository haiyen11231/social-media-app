package models

import (
	"time"

	"gorm.io/gorm"
)

// DB Model
type Post struct {
	gorm.Model
	ContentText     string    `gorm:"size:500"`       // Content text of the post
	ContentImagePath string   `gorm:"size:256"`       // Path to the image in MinIO
	UserID          uint      `gorm:"not null"`       // ID of the user who created the post
	Visible         bool      `gorm:"not null"`       // Visibility of the post
	Comments        []*Comment `gorm:"foreignKey:PostID"` // Comments on the post
	LikedUsers      []*User   `gorm:"many2many:like;foreignKey:id;joinForeignKey:post_id;References:id;joinReferences:user_id"` // Users who liked the post
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