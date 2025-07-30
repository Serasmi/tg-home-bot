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

func RegisterHandler(ctx context.Context, bot *tele.Bot, uc Service) {
	btnAll := menu.Text(allSensorsText)

	sensorButtons := make([]tele.Btn, len(ha.Sensors))
	for i := range ha.Sensors {
		sensorButtons[i] = menu.Text(ha.Sensors[i].ShortName)
	}

	menu.Reply(
		menu.Row(btnAll),
		menu.Row(sensorButtons...),
	)

	bot.Handle(startPath, startHandler)
	bot.Handle(&btnAll, handleSensors(ctx, uc, ha.Sensors))

	for i := range ha.Sensors {
		bot.Handle(&sensorButtons[i], handleSensor(uc, ha.Sensors[i]))
	}
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

func handleSensors(ctx context.Context, uc Service, sensors []ha.Sensor) tele.HandlerFunc {
	return func(c tele.Context) error {
		value, err := uc.GetSensorsValue(ctx, sensors...)
		if err != nil {
			return err
		}

		return c.Send(value)
	}
}
