package models

// DB Model
type Comment struct {
	ID          uint
	UserID      uint
	PostID      uint
	ContentText string
}

func (Comment) TableName() string {
	return "comment"
}

// Request Models
type CreateCommentRequest struct {
	ContentText string `json:"content_text"`
}

// Response Models