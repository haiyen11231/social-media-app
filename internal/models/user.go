package models

import "time"

// DB Model
type User struct {
	ID          uint `gorm:"primaryKey"`
	FirstName   string `gorm:"size:50;not null"`
	LastName    string `gorm:"size:50;not null"`
	DateOfBirth time.Time `gorm:"not null"`
	Email string `gorm:"size:50;not null"`
	Username string `gorm:"size:50;not null;index:idx_username"`
	HashedPassword string `gorm:"size:50;not null"`
	Salt string `gorm:"size:20;not null"`
	Following []*User `gorm:"many2many:following;foreignKey:id;joinForeignKey:user_id;References:id;JoinReferences:follower_id"`
	Followers []*User `gorm:"many2many:following;foreignKey:id;joinForeignKey:follower_id;References:id;JoinReferences:user_id"`
	Posts []*Post `gorm:"foreignKey:UserID"`
}

// Request Models
type LogInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	FirstName string `json:"first_name"`
	LastName    string `json:"last_name"`
	DoB string `json:"dob"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EditUserRequest struct {
	FirstName string `json:"first_name"`
	LastName    string `json:"last_name"`
	DoB string `json:"dob"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response Models
type MessageResponse struct {
	Message string `json:"message"`
}