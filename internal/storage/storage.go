package storage

import (
	"school-noti-telegram-go/internal/models"
)

type Storage interface {
	SaveNotice(notice models.Notice) error
	GetLastCheckedNotices() ([]models.Notice, error)
	IsNoticeSent(noticeID string) bool
}
