package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type SubscriptionRepository struct {
	db *database.Database
}

func NewSubscriptionRepository(db *database.Database) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			user_id, name, description, amount, currency,
			billing_period, billing_day, custom_period_days,
			start_date, next_billing_date, is_active, category, icon, color
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		sub.UserID, sub.Name, sub.Description, sub.Amount, sub.Currency,
		sub.BillingPeriod, sub.BillingDay, sub.CustomPeriodDays,
		sub.StartDate, sub.NextBillingDate, sub.IsActive, sub.Category, sub.Icon, sub.Color,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *SubscriptionRepository) GetByID(id, userID int) (*models.Subscription, error) {
	var sub models.Subscription
	query := `SELECT * FROM subscriptions WHERE id = $1 AND user_id = $2`
	err := r.db.Get(&sub, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetAll(userID int) ([]models.Subscription, error) {
	var subs []models.Subscription
	query := `
		SELECT * FROM subscriptions
		WHERE user_id = $1
		ORDER BY next_billing_date ASC
	`
	err := r.db.Select(&subs, query, userID)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) Update(sub *models.Subscription) error {
	query := `
		UPDATE subscriptions SET
			name = $1, description = $2, amount = $3, currency = $4,
			billing_period = $5, billing_day = $6, custom_period_days = $7,
			next_billing_date = $8, is_active = $9, category = $10, icon = $11, color = $12
		WHERE id = $13 AND user_id = $14
	`
	result, err := r.db.Exec(
		query,
		sub.Name, sub.Description, sub.Amount, sub.Currency,
		sub.BillingPeriod, sub.BillingDay, sub.CustomPeriodDays,
		sub.NextBillingDate, sub.IsActive, sub.Category, sub.Icon, sub.Color,
		sub.ID, sub.UserID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

func (r *SubscriptionRepository) Delete(id, userID int) error {
	query := `DELETE FROM subscriptions WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

func (r *SubscriptionRepository) GetUpcomingBillings(daysAhead int) ([]models.Subscription, error) {
	var subs []models.Subscription
	query := `
		SELECT * FROM subscriptions
		WHERE is_active = true
		AND next_billing_date <= $1
		ORDER BY next_billing_date ASC
	`
	endDate := time.Now().AddDate(0, 0, daysAhead)
	err := r.db.Select(&subs, query, endDate)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) GetStats(userID int) (*models.SubscriptionStats, error) {
	stats := &models.SubscriptionStats{}

	// Get total and active count
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_active = true) as active
		FROM subscriptions
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(&stats.TotalSubscriptions, &stats.ActiveSubscriptions)
	if err != nil {
		return nil, err
	}

	// Calculate monthly total (approximate)
	query = `
		SELECT COALESCE(SUM(
			CASE
				WHEN billing_period = 'monthly' THEN amount
				WHEN billing_period = 'yearly' THEN amount / 12
				WHEN billing_period = 'weekly' THEN amount * 4
				WHEN billing_period = 'custom' AND custom_period_days IS NOT NULL THEN amount * 30.0 / custom_period_days
				ELSE 0
			END
		), 0) as monthly_total
		FROM subscriptions
		WHERE user_id = $1 AND is_active = true
	`
	err = r.db.QueryRow(query, userID).Scan(&stats.MonthlyTotal)
	if err != nil {
		return nil, err
	}

	stats.YearlyTotal = stats.MonthlyTotal * 12

	// Get next payment
	var nextSub models.Subscription
	query = `
		SELECT * FROM subscriptions
		WHERE user_id = $1 AND is_active = true
		ORDER BY next_billing_date ASC
		LIMIT 1
	`
	err = r.db.Get(&nextSub, query, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		stats.NextPayment = &nextSub
	}

	return stats, nil
}
