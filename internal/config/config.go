package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken     string   `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID       int64    `mapstructure:"TELEGRAM_CHAT_ID"`
	SchoolNoticeURL      string   `mapstructure:"SCHOOL_NOTICE_URL"`
	SchoolNoticeKeywords []string `mapstructure:"SCHOOL_NOTICE_KEYWORDS"`
	DeptGeneralURL       string   `mapstructure:"DEPT_GENERAL_URL"`
	DeptScholarshipURL   string   `mapstructure:"DEPT_SCHOLARSHIP_URL"`
	DeptNoticeKeywords   []string `mapstructure:"DEPT_NOTICE_KEYWORDS"`
	RedisAddr            string   `mapstructure:"REDIS_ADDR"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// 필수 설정값 검증
	if config.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN이 설정되지 않았습니다")
	}
	if config.TelegramChatID == 0 {
		return nil, fmt.Errorf("TELEGRAM_CHAT_ID가 설정되지 않았습니다")
	}
	if config.SchoolNoticeURL == "" {
		return nil, fmt.Errorf("SCHOOL_NOTICE_URL이 설정되지 않았습니다")
	}
	if config.DeptGeneralURL == "" {
		return nil, fmt.Errorf("DEPT_GENERAL_URL이 설정되지 않았습니다")
	}
	if config.DeptScholarshipURL == "" {
		return nil, fmt.Errorf("DEPT_SCHOLARSHIP_URL이 설정되지 않았습니다")
	}
	if len(config.DeptNoticeKeywords) == 0 {
		return nil, fmt.Errorf("DEPT_NOTICE_KEYWORDS가 설정되지 않았습니다")
	}
	if config.RedisAddr == "" {
		config.RedisAddr = "localhost:6379"
	}

	return config, nil
}
