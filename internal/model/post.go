package model

import "time"

// DB Model
type Post struct {
	ID               uint
	ContentText      string
	ContentImagePath string
	UserID           uint
	Visible          bool
	Comments         []*Comment
	LikedUsers       []*User
}

// Request Models
type CreatePostRequest struct {
	UserID uint
	ContentText      string
	ContentImagePath string
	Visible          bool
}

type EditPostRequest struct {
	ContentText      string
	ContentImagePath string
	Visible          bool
}

// Response Models
type PostDetailResponse struct {
	PostID uint
	UserID uint
	ContentText string
	ContentImagePath string
	Visible bool
	CreatedTime time.Time
}