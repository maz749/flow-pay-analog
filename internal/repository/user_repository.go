package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, username, auth_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	authType := user.AuthType
	if authType == "" {
		authType = "email"
	}
	return r.db.QueryRow(query, user.Email, user.PasswordHash, user.Username, authType).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) CreateTelegramUser(user *models.User) error {
	query := `
		INSERT INTO users (username, telegram_user_id, telegram_username, telegram_first_name, telegram_last_name, telegram_photo_url, auth_type)
		VALUES ($1, $2, $3, $4, $5, $6, 'telegram')
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, user.Username, user.TelegramUserID, user.TelegramUsername,
		user.TelegramFirstName, user.TelegramLastName, user.TelegramPhotoURL).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE telegram_user_id = $1`
	err := r.db.Get(&user, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetTelegramSettings(userID int) (*models.TelegramSettings, error) {
	var settings models.TelegramSettings
	query := `SELECT * FROM telegram_settings WHERE user_id = $1`
	err := r.db.Get(&settings, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *UserRepository) UpsertTelegramSettings(settings *models.TelegramSettings) error {
	query := `
		INSERT INTO telegram_settings (user_id, chat_id, is_enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			chat_id = EXCLUDED.chat_id,
			is_enabled = EXCLUDED.is_enabled,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, settings.UserID, settings.ChatID, settings.IsEnabled).
		Scan(&settings.ID, &settings.CreatedAt, &settings.UpdatedAt)
}
