package notifier

import (
	"school-noti-telegram-go/internal/models"
)

type TelegramNotifier interface {
	SendMessage(message string) error
	SendNotice(notice models.Notice) error
}
