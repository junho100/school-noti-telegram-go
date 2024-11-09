package main

import (
	"log"
	"time"

	"school-noti-telegram-go/internal/config"
	"school-noti-telegram-go/internal/crawler"
	"school-noti-telegram-go/internal/notifier"
	"school-noti-telegram-go/internal/storage"
)

func main() {
	// 1. 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("설정을 불러오는데 실패했습니다: %v", err)
	}

	// 2. 컴포넌트 초기화
	crawler := crawler.NewCrawler(cfg)
	notifier := notifier.NewTelegramNotifier(cfg)
	storage := storage.NewStorage(cfg)

	// 3. 주기적으로 실행할 메인 로직
	for {
		if err := runCrawlingJob(crawler, notifier, storage, cfg.Keywords); err != nil {
			log.Printf("작업 실행 중 오류 발생: %v", err)
		}

		// 24시간 대기
		time.Sleep(24 * time.Hour)
	}
}

func runCrawlingJob(c crawler.Crawler, n notifier.TelegramNotifier, s storage.Storage, keywords []string) error {
	// 1. 새로운 공지사항 수집
	notices, err := c.FetchNotices()
	if err != nil {
		return err
	}

	// 2. 키워드 필터링
	filteredNotices := c.FilterByKeywords(notices, keywords)

	// 3. 새로운 공지사항 확인 및 알림 전송
	for _, notice := range filteredNotices {
		if !s.IsNoticeSent(notice.ID) {
			if err := n.SendNotice(notice); err != nil {
				log.Printf("공지사항 전송 실패: %v", err)
				continue
			}

			if err := s.SaveNotice(notice); err != nil {
				log.Printf("공지사항 저장 실패: %v", err)
			}
		}
	}

	return nil
}
