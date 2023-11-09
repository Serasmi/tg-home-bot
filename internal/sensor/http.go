package sensor

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

const (
	startPath = "/start"

	Humidity    = "Humidity"
	Temperature = "Temperature"
)

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}

	btnHumidity    = menu.Text(Humidity)
	btnTemperature = menu.Text(Temperature)
)

type sensorService interface {
	SensorValue(sensor string) (string, error)
}

func RegisterHandler(bot *tele.Bot, uc sensorService) {
	menu.Reply(
		menu.Row(btnHumidity),
		menu.Row(btnTemperature),
	)

	bot.Handle(startPath, startHandler)
	bot.Handle(&btnHumidity, humidityHandler(uc))
	bot.Handle(&btnTemperature, temperatureHandler(uc))
}

func startHandler(c tele.Context) error {
	return c.Send("What do you want?", menu)
}

func humidityHandler(uc sensorService) tele.HandlerFunc {
	return func(c tele.Context) error {
		value, err := uc.SensorValue(strings.ToLower(Humidity))
		if err != nil {
			return err
		}

		return c.Send(value)
	}
}

func temperatureHandler(uc sensorService) tele.HandlerFunc {
	return func(c tele.Context) error {
		value, err := uc.SensorValue(strings.ToLower(Temperature))
		if err != nil {
			return err
		}

		return c.Send(value)
	}
}
