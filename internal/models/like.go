package models

// DB Model
type Like struct {
	UserID uint `gorm:"primaryKey"` // ID of the user who liked the post
	PostID uint `gorm:"primaryKey"` // ID of the post that was liked
}

func (Like) TableName() string {
	return "like"
}