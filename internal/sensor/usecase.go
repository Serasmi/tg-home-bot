package sensor

import (
	"fmt"
	ha "tg-home-bot/pkg/home-assistant"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type useCase struct {
	titleCaser cases.Caser
	iotService ha.Service
}

var _ sensorUseCase = &useCase{}

func NewUseCase(s ha.Service) *useCase {
	return &useCase{
		titleCaser: cases.Title(language.English),
		iotService: s,
	}
}

func (u *useCase) SensorValue(sensor string) (string, error) {
	state, err := u.iotService.GetSensorState(sensor)
	if err != nil {
		return "", err
	}

	return u.formatSensorValue(sensor, state)
}

func (u *useCase) formatSensorValue(sensor string, state *ha.SensorRawState) (string, error) {
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
