package models

// DB Model
type Like struct {
	UserID uint `gorm:"primaryKey"`
	PostID uint `gorm:"primaryKey"`
}

func (Like) TableName() string {
	return "like"
}