package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/internal/telegram"
)

type Scheduler struct {
	cron      *cron.Cron
	subRepo   *repository.SubscriptionRepository
	notifRepo *repository.NotificationRepository
	userRepo  *repository.UserRepository
	bot       *telegram.Bot
}

func NewScheduler(
	subRepo *repository.SubscriptionRepository,
	notifRepo *repository.NotificationRepository,
	userRepo *repository.UserRepository,
	bot *telegram.Bot,
) *Scheduler {
	return &Scheduler{
		cron:      cron.New(),
		subRepo:   subRepo,
		notifRepo: notifRepo,
		userRepo:  userRepo,
		bot:       bot,
	}
}

func (s *Scheduler) Start() {
	// Run notification check every hour
	_, err := s.cron.AddFunc("@hourly", s.checkNotifications)
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
		return
	}

	s.cron.Start()
	log.Println("Scheduler started")

	// Run immediately on start
	go s.checkNotifications()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) checkNotifications() {
	log.Println("Checking for notifications...")

	// Get all active subscriptions for the next 30 days
	subscriptions, err := s.subRepo.GetUpcomingBillings(30)
	if err != nil {
		log.Printf("Error fetching upcoming billings: %v", err)
		return
	}

	for _, sub := range subscriptions {
		// Get notifications for this subscription
		notifications, err := s.notifRepo.GetBySubscription(sub.ID)
		if err != nil {
			log.Printf("Error fetching notifications for subscription %d: %v", sub.ID, err)
			continue
		}

		// Get user's telegram settings
		telegramSettings, err := s.userRepo.GetTelegramSettings(sub.UserID)
		if err != nil || telegramSettings == nil || !telegramSettings.IsEnabled || telegramSettings.ChatID == nil {
			continue
		}

		// Calculate days until billing
		now := time.Now()
		daysUntil := int(sub.NextBillingDate.Sub(now).Hours() / 24)

		// Check each notification setting
		for _, notif := range notifications {
			if !notif.IsEnabled {
				continue
			}

			// Check if we should send notification
			if daysUntil == notif.NotifyDaysBefore {
				// Check if we already sent this notification today
				history, _ := s.notifRepo.GetHistory(sub.UserID, 100)
				alreadySent := false

				for _, h := range history {
					if h.SubscriptionID == sub.ID &&
						h.SentAt.After(now.Add(-24*time.Hour)) &&
						h.Status == "sent" {
						alreadySent = true
						break
					}
				}

				if !alreadySent {
					s.sendNotification(telegramSettings, &sub, daysUntil)
				}
			}
		}
	}

	log.Println("Notification check completed")
}

func (s *Scheduler) sendNotification(settings *models.TelegramSettings, sub *models.Subscription, daysUntil int) {
	if settings.ChatID == nil {
		return
	}

	billingDate := sub.NextBillingDate.Format("02.01.2006")

	err := s.bot.SendNotification(
		*settings.ChatID,
		sub.Name,
		sub.Amount,
		sub.Currency,
		daysUntil,
		billingDate,
	)

	// Save to history
	history := &models.NotificationHistory{
		SubscriptionID:   sub.ID,
		UserID:           sub.UserID,
		NotificationType: "telegram",
		Status:           "sent",
	}

	if err != nil {
		log.Printf("Error sending notification: %v", err)
		history.Status = "failed"
		errMsg := err.Error()
		history.ErrorMessage = &errMsg
	} else {
		log.Printf("Notification sent for subscription %s (user %d)", sub.Name, sub.UserID)
	}

	if err := s.notifRepo.CreateHistory(history); err != nil {
		log.Printf("Error saving notification history: %v", err)
	}
}
