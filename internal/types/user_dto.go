package types

import (
	"fmt"
	"strings"
	"time"
)

// RegisterRequestDTO handles the incoming payload for user registration
type RegisterRequestDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate ensures incoming data meets basic validation rules before processing
func (r *RegisterRequestDTO) Validate() error {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" || !strings.Contains(r.Email, "@") {
		return fmt.Errorf("a valid email address is required")
	}
	if len(r.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if strings.TrimSpace(r.FirstName) == "" || strings.TrimSpace(r.LastName) == "" {
		return fmt.Errorf("first name and last name cannot be blank")
	}
	return nil
}

// LoginRequestDTO handles authentication payloads
type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponseDTO is the clean, safe object returned to the client.
// It hides password hashes, tokens, and internal database flags.
type UserResponseDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsVerified bool     `json:"is_verified"`
	CreatedAt time.Time `json:"created_at"`
}

