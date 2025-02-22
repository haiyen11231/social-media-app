package models

// DB Model
type Following struct {
	UserID     uint `gorm:"primaryKey"`
	FollowerID uint `gorm:"primaryKey"`
}

func (Following) TableName() string {
	return "following"
}