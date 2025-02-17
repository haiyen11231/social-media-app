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