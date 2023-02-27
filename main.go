package main

import (
	"log"
	"os"
	"tg-home-bot/internal/config"
	"tg-home-bot/pkg/logging"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	logger := logging.NewLogger()

	cfg := config.InitConfig(logger)

	logger.SetLevel(cfg.App.LogLevel)

	logger.Info("Start app")

	pref := tele.Settings{
		Token:  os.Getenv("TG_API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Start()
}
