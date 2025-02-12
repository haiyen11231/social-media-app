package models

import "time"

// DB Model
type User struct {
	ID          uint
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Email string
	Username string
	HashedPassword string
	Salt string
	Following []*User
	Followers []*User
	Posts []*Post
}

// Request Models
type LogInRequest struct {
	Username string
	Password string
}

type SignUpRequest struct {
	FirstName string
	LastName    string
	DoB string
	Email string
	Username string
	Password string
}

type UpdateUserRequest struct {
	UserID uint
	FirstName string
	LastName    string
	DoB string
	Email string
	Username string
	Password string
}

// Response Models
type MessageResponse struct {
	Message string
}