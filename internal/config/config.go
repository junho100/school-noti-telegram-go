package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken string   `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID   int64    `mapstructure:"TELEGRAM_CHAT_ID"`
	SchoolNoticeURL  string   `mapstructure:"SCHOOL_NOTICE_URL"`
	Keywords         []string `mapstructure:"KEYWORDS"`
	CheckInterval    string   `mapstructure:"CHECK_INTERVAL"`
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
	if len(config.Keywords) == 0 {
		return nil, fmt.Errorf("KEYWORDS가 설정되지 않았습니다")
	}
	if config.CheckInterval == "" {
		config.CheckInterval = "24h" // 기본값 설정
	}

	return config, nil
}
