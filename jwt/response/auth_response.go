package response

import (
	"time"
)

// SignIn struct to describe login response.
type SignIn struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

// User struct to describe User response.
// The instance should be created from models.User.ToResponse()
type User struct {
	ID           int         `json:"id"`
	Email        string      `json:"email"`
	Fullname     string      `json:"fullname"`
	Phone        string      `json:"phone"`
	Token        interface{} `json:"token"`
	Status       string      `json:"status"`
	Avatar       interface{} `json:"avatar"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	VerifiedAt   interface{} `json:"verified_at"`
	BlockedAt    interface{} `json:"blocked_at"`
	DeletedAt    interface{} `json:"deleted_at"`
	LastAccessAt interface{} `json:"last_access_at"`
}
