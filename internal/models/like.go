package models

// DB Model
type Like struct {
	UserID uint
	PostID uint
}

func (Like) TableName() string {
	return "like"
}