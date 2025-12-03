package models

import (
	"time"
)

type Subscription struct {
	ID                int       `db:"id" json:"id"`
	UserID            int       `db:"user_id" json:"user_id"`
	TeamID            *int      `db:"team_id" json:"team_id"`
	Name              string    `db:"name" json:"name"`
	Description       *string   `db:"description" json:"description"`
	Amount            float64   `db:"amount" json:"amount"`
	Currency          string    `db:"currency" json:"currency"`
	BillingPeriod     string    `db:"billing_period" json:"billing_period"` // monthly, yearly, weekly, custom
	BillingDay        *int      `db:"billing_day" json:"billing_day"`
	CustomPeriodDays  *int      `db:"custom_period_days" json:"custom_period_days"`
	StartDate         time.Time `db:"start_date" json:"start_date"`
	NextBillingDate   time.Time `db:"next_billing_date" json:"next_billing_date"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	Category          *string   `db:"category" json:"category"`
	Icon              *string   `db:"icon" json:"icon"`
	Color             *string   `db:"color" json:"color"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	Notifications     []Notification `json:"notifications,omitempty"`
}

type CreateSubscriptionRequest struct {
	Name             string   `json:"name"`
	Description      *string  `json:"description"`
	Amount           float64  `json:"amount"`
	Currency         string   `json:"currency"`
	BillingPeriod    string   `json:"billing_period"`
	BillingDay       *int     `json:"billing_day"`
	CustomPeriodDays *int     `json:"custom_period_days"`
	StartDate        string   `json:"start_date"` // YYYY-MM-DD
	Category         *string  `json:"category"`
	Icon             *string  `json:"icon"`
	Color            *string  `json:"color"`
	NotifyDaysBefore []int    `json:"notify_days_before"`
}

type UpdateSubscriptionRequest struct {
	Name             *string  `json:"name"`
	Description      *string  `json:"description"`
	Amount           *float64 `json:"amount"`
	Currency         *string  `json:"currency"`
	BillingPeriod    *string  `json:"billing_period"`
	BillingDay       *int     `json:"billing_day"`
	CustomPeriodDays *int     `json:"custom_period_days"`
	IsActive         *bool    `json:"is_active"`
	Category         *string  `json:"category"`
	Icon             *string  `json:"icon"`
	Color            *string  `json:"color"`
}

type SubscriptionStats struct {
	TotalSubscriptions int     `json:"total_subscriptions"`
	ActiveSubscriptions int    `json:"active_subscriptions"`
	MonthlyTotal       float64 `json:"monthly_total"`
	YearlyTotal        float64 `json:"yearly_total"`
	NextPayment        *Subscription `json:"next_payment"`
}
