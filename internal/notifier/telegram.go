package notifier

import (
	"fmt"
	"school-noti-telegram-go/internal/config"
	"school-noti-telegram-go/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotifier(cfg *config.Config) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return nil, fmt.Errorf("텔레그램 봇 초기화 실패: %v", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: cfg.TelegramChatID,
	}, nil
}

// 일반 메시지 전송
func (t *TelegramNotifier) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("텔레그램 메시지 전송 실패: %v", err)
	}
	return nil
}

// 공지사항 전송
func (t *TelegramNotifier) SendNotice(notice models.Notice) error {
	var noticeType string
	switch notice.Type {
	case models.SchoolNotice:
		noticeType = "[학교 공지사항]"
	case models.DeptGeneralNotice:
		noticeType = "[학과 공지사항]"
	case models.DeptScholarshipNotice:
		noticeType = "[학과 장학금 공지사항]"
	}

	message := fmt.Sprintf(
		"%s\n제목: %s\n링크: %s\n작성일: %s",
		noticeType,
		notice.Title,
		notice.URL,
		notice.PostDate.Format("2006-01-02"),
	)

	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("텔레그램 공지사항 전송 실패: %v", err)
	}

	return nil
}
