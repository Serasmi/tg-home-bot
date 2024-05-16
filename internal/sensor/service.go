package sensor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	ha "tg-home-bot/pkg/home-assistant"
	"tg-home-bot/pkg/logger"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service interface {
	GetSensorValue(ctx context.Context, sensor ha.Sensor) (string, error)
	GetSensorsValue(ctx context.Context, sensors ...ha.Sensor) (string, error)
}

type service struct {
	location   *time.Location
	timeout    time.Duration
	titleCaser cases.Caser
	iotService ha.Service
}

func NewService(s ha.Service, location *time.Location, timeout time.Duration) Service {
	return &service{
		location:   location,
		timeout:    timeout,
		titleCaser: cases.Title(language.English),
		iotService: s,
	}
}

func (s *service) GetSensorsValue(ctx context.Context, sensors ...ha.Sensor) (string, error) {
	if len(sensors) == 0 {
		return "", nil
	}

	states, err := s.iotService.GetSensorsState(ctx, sensors...)
	if err != nil {
		return "", err
	}

	if len(states) == 0 {
		return "", ha.ErrNoSensor
	}

	if len(states) != len(sensors) {
		return "", errors.New("invalid HA response")
	}

	res := make([]string, len(states))
	for i := range states {
		res[i], err = s.formatSensorValue(sensors[i], &states[i])
		if err != nil {
			logger.GetLogger().WithFields(logger.Fields{
				"sensor": sensors[i],
				"state":  states[i],
			}).WithError(err).Error("format sensor value")

			return "", err
		}
	}

	return strings.Join(res, "\n\n"), nil
}

func (s *service) GetSensorValue(ctx context.Context, sensor ha.Sensor) (string, error) {
	state, err := s.iotService.GetSensorState(ctx, sensor)
	if err != nil {
		return "", err
	}

	return s.formatSensorValue(sensor, state)
}

func (s *service) formatSensorValue(sensor ha.Sensor, state *ha.SensorRawState) (string, error) {
	return fmt.Sprintf(
		"%s%s is %s.\nLast updated at %s",
		s.getIcon(sensor),
		sensor.FriendlyName,
		state.State+state.Attributes.UnitOfMeasurement,
		state.UpdatedAt.In(s.location).Format(time.RFC822),
	), nil
}

func (s *service) getIcon(sensor ha.Sensor) string {
	switch sensor {
	case ha.SensorHumidity:
		return "💧"
	case ha.SensorTemperature, ha.SensorRPITemperature:
		return "🌡️"
	default:
		return ""
	}
}
