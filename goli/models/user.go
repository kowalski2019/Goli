package models

import "time"

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email,omitempty"`
	Role              string    `json:"role"` // "admin" or "user"
	Phone             string    `json:"phone,omitempty"`
	Password          string    `json:"-"`                    // Never serialize password
	PasswordHash      string    `json:"-"`                    // Hashed password stored in DB
	TwoFAEmailEnabled int       `json:"two_fa_email_enabled"` // 0/1
	TwoFASmsEnabled   int       `json:"two_fa_sms_enabled"`   // 0/1
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserCreateRequest represents a request to create a user
type UserCreateRequest struct {
	Username          string `json:"username"`
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	TwoFAEmailEnabled int    `json:"two_fa_email_enabled"`
	TwoFASmsEnabled   int    `json:"two_fa_sms_enabled"`
}

// UserUpdateRequest represents a request to update a user
type UserUpdateRequest struct {
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	Password          string `json:"password,omitempty"`
	Role              string `json:"role,omitempty"`
	TwoFAEmailEnabled *int   `json:"two_fa_email_enabled,omitempty"`
	TwoFASmsEnabled   *int   `json:"two_fa_sms_enabled,omitempty"`
}
