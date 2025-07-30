package main

import (
	"log/slog"
	"os"
	"time"
	_ "time/tzdata"

	"tg-home-bot/internal/config"
	"tg-home-bot/internal/echo"
	"tg-home-bot/internal/middleware"
	"tg-home-bot/internal/sensor"
	ha "tg-home-bot/pkg/home-assistant"

	tele "gopkg.in/telebot.v3"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		slog.Error("init config", "error", err)
		return
	}

	initLogger(cfg.App.LogLevel)

	slog.Info("conf: permitted users", "users", cfg.Telegram.PermitUsers)

	_, err = initBot(cfg)
	if err != nil {
		slog.Error("init bot", "error", err)
		return
	}
}

func initBot(config *config.Config) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("TG_API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, c tele.Context) {
			if c != nil {
				slog.Error("bot error", "update_id", c.Update().ID, "error", err.Error())
				return
			}

			slog.Error("bot error", "error", err.Error())
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	haProvider := ha.NewService(config.HomeAssistant.URL, config.HomeAssistant.Token)

	b.Use(middleware.PermitUsers(config.Telegram.PermitUsers))

	echo.RegisterHandler(b)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	sensor.RegisterHandler(b, sensor.NewService(haProvider, loc, config.HomeAssistant.Timeout))

	b.Start()

	return b, nil
}

func initLogger(level string) {
	var lvl slog.Level

	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		slog.Default().Error("unmarshal log level", "level", level, "error", err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))

	slog.SetDefault(log)
}
