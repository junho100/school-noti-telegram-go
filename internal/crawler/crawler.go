package crawler

import (
	"school-noti-telegram-go/internal/models"
)

type Crawler interface {
	FetchNotices() ([]models.Notice, error)
	FilterByKeywords(notices []models.Notice, keywords []string) []models.Notice
}
