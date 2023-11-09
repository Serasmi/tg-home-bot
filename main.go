package main

import (
	"os"
	"time"

	"tg-home-bot/internal/config"
	"tg-home-bot/internal/echo"
	"tg-home-bot/internal/middleware"
	"tg-home-bot/internal/sensor"
	ha "tg-home-bot/pkg/home-assistant"
	"tg-home-bot/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

func main() {
	cfg := config.InitConfig()

	logger.GetLogger().SetLevel(cfg.App.LogLevel)
	logger.GetLogger().Infof("[init] permit users: %v", cfg.Telegram.PermitUsers)

	_, err := initBot(cfg)
	if err != nil {
		logger.GetLogger().Fatal("[init] bot init", err)
	}
}

func initBot(config *config.Config) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("TG_API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, c tele.Context) {
			if c != nil {
				logger.GetLogger().WithField("update_id", c.Update().ID).Error(err)
				return
			}

			logger.GetLogger().Error(err)
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	haProvider := ha.NewService(config.HomeAssistant.URL, config.HomeAssistant.Token)

	b.Use(middleware.PermitUsers(config.Telegram.PermitUsers))

	echo.RegisterHandler(b)
	sensor.RegisterHandler(b, sensor.NewService(haProvider))

	b.Start()

	return b, nil
}
