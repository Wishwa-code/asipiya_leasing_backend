package models

// LoginRequest defines the structure for authentication attempts 👤
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}