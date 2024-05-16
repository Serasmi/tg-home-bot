package config

import (
	"time"

	"tg-home-bot/pkg/logger"

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

var config Config

func GetConfig() *Config {
	return &config
}

func InitConfig() *Config {
	logger.GetLogger().Debug("init config")

	err := godotenv.Load()
	if err != nil {
		logger.GetLogger().Fatalf("Error loading .env file: %s", err.Error())
	}

	err = env.Parse(&config)
	if err != nil {
		logger.GetLogger().Fatal("init config error ", err)
	}

	return &config
}
