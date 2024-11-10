package main

import (
	"fmt"
	"log"
	"time"

	"school-noti-telegram-go/internal/config"
	"school-noti-telegram-go/internal/crawler"
	"school-noti-telegram-go/internal/models"
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
	notifierSvc, err := notifier.NewTelegramNotifier(cfg)
	if err != nil {
		log.Fatalf("Telegram 초기화 실패: %v", err)
	}
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

	c := cron.New(cron.WithLocation(loc))

	// 매일 오전 11시와 오후 11시에 실행
	if _, err := c.AddFunc("0 11,23 * * *", func() {
		now := time.Now().In(loc)
		checkTime := now.Format("2006-01-02 15:04")

		notices, err := crawlerSvc.FetchAllNotices()
		if err != nil {
			log.Printf("공지사항 크롤링 실패: %v", err)
			return
		}

		// 새로운 공지사항 필터링
		var newNotices []models.Notice
		for _, notice := range notices {
			if !storageSvc.IsNoticeProcessed(notice.ID) {
				newNotices = append(newNotices, notice)
				if err := storageSvc.MarkNoticeAsProcessed(notice.ID); err != nil {
					log.Printf("공지사항 처리 상태 저장 실패: %v", err)
				}
			}
		}

		// 새로운 공지사항이 없는 경우
		if len(newNotices) == 0 {
			message := fmt.Sprintf("발견된 공지사항이 없습니다.\n확인 시각: %s", checkTime)
			if err := notifierSvc.SendMessage(message); err != nil {
				log.Printf("메시지 전송 실패: %v", err)
			}
			return
		}

		// 새로운 공지사항 전송
		for _, notice := range newNotices {
			if err := notifierSvc.SendNotice(notice); err != nil {
				log.Printf("공지사항 전송 실패: %v", err)
			}
		}
	}); err != nil {
		log.Fatalf("크론 작업 설정 실패: %v", err)
	}

	c.Start()

	// 프로그램 종료 방지
	select {}
}
