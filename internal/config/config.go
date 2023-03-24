package config

import (
	"tg-home-bot/pkg/logging"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		LogLevel string `env:"APP_LOG_LEVEL" envDefault:"info"`
	}
	HomeAssistant struct {
		URL   string `env:"HA_API_BASEURL,required"`
		Token string `env:"HA_API_TOKEN,required"`
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

func InitConfig(logger *logging.Logger) *Config {
	logger.Debug("init config")

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err.Error())
	}

	err = env.Parse(&config)
	if err != nil {
		logger.Fatal("init config error ", err)
	}

	return &config
}
