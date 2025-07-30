package sensor

import (
	"context"

	ha "tg-home-bot/pkg/home-assistant"

	tele "gopkg.in/telebot.v3"
)

const (
	startPath      = "/start"
	allSensorsText = "All"
)

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
)

func RegisterHandler(bot *tele.Bot, uc Service) {
	var (
		btnAll            = menu.Text(allSensorsText)
		btnTemperature    = menu.Text(ha.SensorTemperature.FriendlyName)
		btnHumidity       = menu.Text(ha.SensorHumidity.FriendlyName)
		btnRPITemperature = menu.Text(ha.SensorRPITemperature.FriendlyName)
		btnServerState    = menu.Text(ha.SensorServerState.FriendlyName)
	)

	menu.Reply(
		menu.Row(btnAll),
		menu.Row(btnTemperature, btnHumidity, btnRPITemperature, btnServerState),
	)

	bot.Handle(startPath, startHandler)
	bot.Handle(&btnAll, handleSensors(uc))
	bot.Handle(&btnTemperature, handleSensor(uc, ha.SensorTemperature))
	bot.Handle(&btnHumidity, handleSensor(uc, ha.SensorHumidity))
	bot.Handle(&btnRPITemperature, handleSensor(uc, ha.SensorRPITemperature))
	bot.Handle(&btnServerState, handleSensor(uc, ha.SensorServerState))
}

func startHandler(c tele.Context) error {
	return c.Send("What do you want?", menu)
}

func handleSensor(uc Service, sensor ha.Sensor) tele.HandlerFunc {
	return func(c tele.Context) error {
		value, err := uc.GetSensorValue(context.Background(), sensor)
		if err != nil {
			return err
		}

		return c.Send(value)
	}
}

func handleSensors(uc Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		value, err := uc.GetSensorsValue(
			context.Background(),
			ha.SensorTemperature,
			ha.SensorHumidity,
			ha.SensorRPITemperature,
			ha.SensorServerState,
		)
		if err != nil {
			return err
		}

		return c.Send(value)
	}
}
