package echo

import (
	"gopkg.in/telebot.v3"
)

const (
	echoPath = "/echo"
	sayPath  = "/say"
)

func RegisterHandler(bot *telebot.Bot) {
	bot.Handle(echoPath, echoHandler)
	bot.Handle(sayPath, echoHandler)
}

func echoHandler(c telebot.Context) error {
	command := c.Message().Payload
	if command == "" {
		return nil
	}

	return c.Send(command)
}
