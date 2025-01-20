package model

// DB Model
type Comment struct {
	ID          uint
	ContentText string
	UserID      uint
	PostID      uint
}

// Request Models
type CreateCommentRequest struct {
	ContentText string
}

// Response Models