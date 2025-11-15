package repository

import (
	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type NotificationRepository struct {
	db *database.Database
}

func NewNotificationRepository(db *database.Database) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (subscription_id, notify_days_before, is_enabled)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		notification.SubscriptionID,
		notification.NotifyDaysBefore,
		notification.IsEnabled,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *NotificationRepository) GetBySubscription(subscriptionID int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notifications WHERE subscription_id = $1 ORDER BY notify_days_before DESC`
	err := r.db.Select(&notifications, query, subscriptionID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) DeleteBySubscription(subscriptionID int) error {
	query := `DELETE FROM notifications WHERE subscription_id = $1`
	_, err := r.db.Exec(query, subscriptionID)
	return err
}

func (r *NotificationRepository) CreateHistory(history *models.NotificationHistory) error {
	query := `
		INSERT INTO notification_history
		(subscription_id, user_id, notification_type, status, error_message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sent_at
	`
	return r.db.QueryRow(
		query,
		history.SubscriptionID,
		history.UserID,
		history.NotificationType,
		history.Status,
		history.ErrorMessage,
	).Scan(&history.ID, &history.SentAt)
}

func (r *NotificationRepository) GetHistory(userID int, limit int) ([]models.NotificationHistory, error) {
	var history []models.NotificationHistory
	query := `
		SELECT * FROM notification_history
		WHERE user_id = $1
		ORDER BY sent_at DESC
		LIMIT $2
	`
	err := r.db.Select(&history, query, userID, limit)
	if err != nil {
		return nil, err
	}
	return history, nil
}
