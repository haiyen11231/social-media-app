package models

// Following represents the many-to-many relationship between users
type Following struct {
	UserID     uint `gorm:"primaryKey"`
	FollowerID uint `gorm:"primaryKey"`
}

// TableName explicitly sets the table name for the Following model
func (Following) TableName() string {
	return "following"
}