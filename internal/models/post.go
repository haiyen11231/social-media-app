package models

import "time"

// DB Model
type Post struct {
	ID               uint
	UserID           uint
	ContentText      string
	ContentImagePath string
	Visible          bool
	Comments         []*Comment
	LikedUsers       []*User
}

// Request Models
type CreatePostRequest struct {
	UserID uint `json:"user_id"`
	ContentText      string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible          bool `json:"visible"`
}

type EditPostRequest struct {
	ContentText      string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible          bool `json:"visible"`
}

// Response Models
type PostDetailResponse struct {
	PostID uint `json:"post_id"`
	UserID uint `json:"user_id"`
	ContentText string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible bool `json:"visible"`
	CreatedTime time.Time `json:"created_time"`
}