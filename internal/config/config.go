package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		LogLevel string `env:"APP_LOG_LEVEL" envDefault:"info"`
	}
	HomeAssistant struct {
		Timeout time.Duration `env:"HA_TIMEOUT" envDefault:"3s"`
		URL     string        `env:"HA_API_BASEURL,required"`
		Token   string        `env:"HA_API_TOKEN,required"`
	}
	Telegram struct {
		PermitUsers []int64 `env:"TG_PERMIT_USERS"`
		Token       string  `env:"TG_API_TOKEN,required"`
	}
}

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("load .env file: %w", err)
	}

	var config Config

	err = env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("parse .env file: %w", err)
	}

	return &config, nil
}
