package storage

import (
	"context"
	"fmt"
	"time"

	"school-noti-telegram-go/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(cfg *config.Config) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx := context.Background()

	// Redis 연결 테스트
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 연결 실패: %v", err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

// 공지사항 처리 여부 저장 (당일)
func (s *RedisStorage) MarkNoticeAsProcessed(noticeID string) error {
	// 키 형식: notice:{noticeID}
	key := fmt.Sprintf("notice:%s", noticeID)

	// 현재 시간 기준으로 다음 날 자정까지의 만료 시간 계산
	tomorrow := time.Now().Add(24 * time.Hour)
	expiryTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	ttl := time.Until(expiryTime)

	return s.client.Set(s.ctx, key, "1", ttl).Err()
}

// 공지사항이 이미 처리되었는지 확인
func (s *RedisStorage) IsNoticeProcessed(noticeID string) bool {
	key := fmt.Sprintf("notice:%s", noticeID)
	exists, err := s.client.Exists(s.ctx, key).Result()
	if err != nil {
		return false
	}
	return exists == 1
}

// Redis 연결 종료
func (s *RedisStorage) Close() error {
	return s.client.Close()
}
