package models

import (
	"time"
)

type Notification struct {
	ID               int       `db:"id" json:"id"`
	SubscriptionID   int       `db:"subscription_id" json:"subscription_id"`
	NotifyDaysBefore int       `db:"notify_days_before" json:"notify_days_before"`
	IsEnabled        bool      `db:"is_enabled" json:"is_enabled"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type NotificationHistory struct {
	ID               int       `db:"id" json:"id"`
	SubscriptionID   int       `db:"subscription_id" json:"subscription_id"`
	UserID           int       `db:"user_id" json:"user_id"`
	NotificationType string    `db:"notification_type" json:"notification_type"`
	SentAt           time.Time `db:"sent_at" json:"sent_at"`
	Status           string    `db:"status" json:"status"`
	ErrorMessage     *string   `db:"error_message" json:"error_message"`
}

type NotificationPayload struct {
	UserID           int
	SubscriptionName string
	Amount           float64
	Currency         string
	DaysUntilBilling int
	BillingDate      time.Time
}
