package models

import (
	"time"
)

// User represents a user of the uptime monitor system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is never returned in JSON
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Role      string    `json:"role"` // "admin", "user"
	Active    bool      `json:"active"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the data needed for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// AuthResponse represents the data returned after successful authentication
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
