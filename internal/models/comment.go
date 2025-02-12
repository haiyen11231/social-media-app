package models

// DB Model
type Comment struct {
	ID          uint
	UserID      uint
	PostID      uint
	ContentText string
}

// Request Models
type CreateCommentRequest struct {
	ContentText string
}

// Response Models