package main

import (
	"os"
	"tg-home-bot/internal/config"
	"tg-home-bot/internal/echo"
	"tg-home-bot/internal/middleware"
	"tg-home-bot/internal/sensor"
	ha "tg-home-bot/pkg/home-assistant"
	"tg-home-bot/pkg/logging"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	logger := logging.NewLogger()

	cfg := config.InitConfig(logger)

	logger.Info(cfg.Telegram.PermitUsers)

	logger.SetLevel(cfg.App.LogLevel)

	_, err := initBot(cfg, logger)
	if err != nil {
		logger.Fatal("Failed bot init", err)
	}
}

func initBot(config *config.Config, logger *logging.Logger) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("TG_API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	haProvider := ha.NewService(ha.NewClient(config.HomeAssistant.URL, config.HomeAssistant.Token))

	b.Use(middleware.PermitUsers(config.Telegram.PermitUsers))

	echo.RegisterHandler(b, logger)
	sensor.RegisterHandler(b, sensor.NewUseCase(haProvider), logger)

	b.Start()

	return b, nil
}
