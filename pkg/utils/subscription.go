package utils

import (
	"time"
)

// CalculateNextBillingDate calculates the next billing date based on the billing period
func CalculateNextBillingDate(startDate time.Time, billingPeriod string, billingDay *int, customDays *int) time.Time {
	now := time.Now()
	current := startDate

	switch billingPeriod {
	case "monthly":
		// Move to the next month if start date is in the past
		for current.Before(now) {
			current = current.AddDate(0, 1, 0)
		}
		// Adjust to billing day if specified
		if billingDay != nil && *billingDay > 0 && *billingDay <= 31 {
			year, month, _ := current.Date()
			current = time.Date(year, month, *billingDay, 0, 0, 0, 0, current.Location())
			// If the adjusted date is before now, move to next month
			if current.Before(now) {
				current = current.AddDate(0, 1, 0)
			}
		}

	case "yearly":
		// Move to the next year if start date is in the past
		for current.Before(now) {
			current = current.AddDate(1, 0, 0)
		}

	case "weekly":
		// Move forward by weeks
		for current.Before(now) {
			current = current.AddDate(0, 0, 7)
		}

	case "custom":
		if customDays != nil && *customDays > 0 {
			// Move forward by custom days
			for current.Before(now) {
				current = current.AddDate(0, 0, *customDays)
			}
		}
	}

	return current
}

// UpdateNextBillingDate updates the next billing date after a payment
func UpdateNextBillingDate(currentDate time.Time, billingPeriod string, customDays *int) time.Time {
	switch billingPeriod {
	case "monthly":
		return currentDate.AddDate(0, 1, 0)
	case "yearly":
		return currentDate.AddDate(1, 0, 0)
	case "weekly":
		return currentDate.AddDate(0, 0, 7)
	case "custom":
		if customDays != nil && *customDays > 0 {
			return currentDate.AddDate(0, 0, *customDays)
		}
	}
	return currentDate
}
