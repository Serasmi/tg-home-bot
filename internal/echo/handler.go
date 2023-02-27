package echo

import (
	"tg-home-bot/pkg/logging"

	"gopkg.in/telebot.v3"
)

const (
	echoPath = "/echo"
)

func RegisterHandler(bot *telebot.Bot, logger *logging.Logger) {
	bot.Handle(echoPath, echoHandler)
}

func echoHandler(c telebot.Context) error {
	command := c.Message().Payload
	if command == "" {
		return nil
	}

	return c.Send(command)
}
