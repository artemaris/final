package models

import "time"

// User represents a registered user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never serialize password
	CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a request to authenticate a user
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents a successful authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
