package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	userRepo *repository.UserRepository
}

func NewBot(token string, db *database.Database) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is empty")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Printf("Telegram Bot authorized on account %s", api.Self.UserName)

	return &Bot{
		api:      api,
		userRepo: repository.NewUserRepository(db),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	if strings.HasPrefix(text, "/start") {
		b.sendMessage(chatID, fmt.Sprintf(
			"👋 Добро пожаловать в FlowPay!\n\n"+
				"Ваш Chat ID: %d\n\n"+
				"Используйте этот ID в настройках веб-приложения для получения уведомлений о подписках.",
			chatID,
		))
	} else if strings.HasPrefix(text, "/help") {
		b.sendMessage(chatID,
			"📖 Помощь по боту FlowPay\n\n"+
			"/start - Получить Chat ID\n"+
			"/help - Показать эту справку\n\n"+
			"Чтобы начать получать уведомления:\n"+
			"1. Скопируйте ваш Chat ID\n"+
			"2. Введите его в настройках Telegram в веб-приложении\n"+
			"3. Включите уведомления",
		)
	} else {
		b.sendMessage(chatID, "Используйте /start для получения Chat ID или /help для справки")
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) SendNotification(chatID int64, subscriptionName string, amount float64, currency string, daysUntil int, billingDate string) error {
	var message string

	if daysUntil == 0 {
		message = fmt.Sprintf(
			"🔔 <b>Сегодня списание!</b>\n\n"+
				"Подписка: <b>%s</b>\n"+
				"Сумма: <b>%.2f %s</b>\n"+
				"Дата: <b>%s</b>",
			subscriptionName, amount, currency, billingDate,
		)
	} else if daysUntil == 1 {
		message = fmt.Sprintf(
			"⏰ <b>Напоминание о подписке</b>\n\n"+
				"Подписка: <b>%s</b>\n"+
				"Сумма: <b>%.2f %s</b>\n"+
				"Списание: <b>завтра</b> (%s)",
			subscriptionName, amount, currency, billingDate,
		)
	} else {
		message = fmt.Sprintf(
			"⏰ <b>Напоминание о подписке</b>\n\n"+
				"Подписка: <b>%s</b>\n"+
				"Сумма: <b>%.2f %s</b>\n"+
				"Списание через: <b>%d дней</b> (%s)",
			subscriptionName, amount, currency, daysUntil, billingDate,
		)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"

	if _, err := b.api.Send(msg); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
