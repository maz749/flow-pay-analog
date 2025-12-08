package models

import (
	"time"
)

type User struct {
	ID                int       `db:"id" json:"id"`
	Email             *string   `db:"email" json:"email,omitempty"`
	PasswordHash      *string   `db:"password_hash" json:"-"`
	Username          string    `db:"username" json:"username"`
	TelegramUserID    *int64    `db:"telegram_user_id" json:"telegram_user_id,omitempty"`
	TelegramUsername  *string   `db:"telegram_username" json:"telegram_username,omitempty"`
	TelegramFirstName *string   `db:"telegram_first_name" json:"telegram_first_name,omitempty"`
	TelegramLastName  *string   `db:"telegram_last_name" json:"telegram_last_name,omitempty"`
	TelegramPhotoURL  *string   `db:"telegram_photo_url" json:"telegram_photo_url,omitempty"`
	AuthType          string    `db:"auth_type" json:"auth_type"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type TelegramSettings struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	ChatID    *int64    `db:"chat_id" json:"chat_id"`
	IsEnabled bool      `db:"is_enabled" json:"is_enabled"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type TelegramAuthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}
