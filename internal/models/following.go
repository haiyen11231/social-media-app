package models

// DB Model
type Following struct {
	UserID     uint `gorm:"primaryKey"`     // The user who is being followed
	FollowerID uint `gorm:"primaryKey"`     // The user who is following
}

func (Following) TableName() string {
	return "following"
}