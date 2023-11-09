package sensor

import (
	"fmt"
	"time"
	_ "time/tzdata"

	ha "tg-home-bot/pkg/home-assistant"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service interface {
	SensorValue(sensor string) (string, error)
}

type service struct {
	titleCaser cases.Caser
	iotService ha.Service
}

func NewService(s ha.Service) Service {
	return &service{
		titleCaser: cases.Title(language.English),
		iotService: s,
	}
}

func (u *service) SensorValue(sensor string) (string, error) {
	state, err := u.iotService.GetSensorState(sensor)
	if err != nil {
		return "", err
	}

	return u.formatSensorValue(sensor, state)
}

func (u *service) formatSensorValue(sensor string, state *ha.SensorRawState) (string, error) {
	value := fmt.Sprintf("%s%s", state.State, state.Attributes.UnitOfMeasurement)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s is %s.\nLast updated at %s",
		u.titleCaser.String(sensor),
		value,
		state.UpdatedAt.In(loc).Format(time.RFC822),
	), nil
}
