package main

import (
	"log"
	"time"

	"school-noti-telegram-go/internal/config"
	"school-noti-telegram-go/internal/crawler"
	"school-noti-telegram-go/internal/notifier"
	"school-noti-telegram-go/internal/storage"

	"github.com/robfig/cron/v3"
)

func main() {
	// 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("설정을 불러오는데 실패했습니다: %v", err)
	}

	// 컴포넌트 초기화
	crawlerSvc := crawler.NewCrawler(cfg)
	notifierSvc := notifier.NewTelegramNotifier(cfg)
	storageSvc, err := storage.NewRedisStorage(cfg)
	if err != nil {
		log.Fatalf("Redis 초기화 실패: %v", err)
	}
	defer storageSvc.Close()

	// 한국 시간대 설정
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatalf("시간대 설정 실패: %v", err)
	}

	// cron 작업 설정
	c := cron.New(cron.WithLocation(loc))

	// 매일 오전 10시와 오후 10시(22시)에 실행
	if _, err := c.AddFunc("0 10,22 * * *", func() {
		if err := runCrawlingJob(crawlerSvc, notifierSvc, storageSvc, cfg.Keywords); err != nil {
			log.Printf("작업 실행 중 오류 발생: %v", err)
		}
	}); err != nil {
		log.Fatalf("크론 작업 설정 실패: %v", err)
	}

	// 크론 작업 시작
	c.Start()

	// 프로그램 시작 시 즉시 한 번 실행
	if err := runCrawlingJob(crawlerSvc, notifierSvc, storageSvc, cfg.Keywords); err != nil {
		log.Printf("초기 작업 실행 중 오류 발생: %v", err)
	}

	// 프로그램이 종료되지 않도록 대기
	select {}
}

func runCrawlingJob(c crawler.Crawler, n notifier.TelegramNotifier, s *storage.RedisStorage, keywords []string) error {
	notices, err := c.FetchNotices()
	if err != nil {
		return err
	}

	filteredNotices := c.FilterByKeywords(notices, keywords)

	for _, notice := range filteredNotices {
		if !s.IsNoticeProcessed(notice.ID) {
			if err := n.SendNotice(notice); err != nil {
				log.Printf("공지사항 전송 실패: %v", err)
				continue
			}

			if err := s.MarkNoticeAsProcessed(notice.ID); err != nil {
				log.Printf("공지사항 처리 상태 저장 실패: %v", err)
			}
		}
	}

	return nil
}
