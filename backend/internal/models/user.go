package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// User represents a user in the database.
type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	FirstName       *string   `json:"firstname,omitempty"`
	LastName        *string   `json:"lastname,omitempty"`
	DateOfBirth     *string   `json:"dateofbirth,omitempty"`
	Avatar          *string   `json:"avatar,omitempty"`
	Nickname        *string   `json:"nickname,omitempty"`
	AboutMe         *string   `json:"aboutme,omitempty"`
	IsProfilePublic bool      `json:"is_profile_public"`
	CreatedAt       time.Time `json:"created_at"`
}

// LoginRequest represents the request body for a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// StepOneCredintial represents the user's credentials on step One of registration.
type StepOneCredintial struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// String returns a pretty-printed multiline JSON-style representation of the User.
func (user User) String() string {
	data, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		return fmt.Sprintf("User{error: %v}", err)
	}
	return string(data)
}
