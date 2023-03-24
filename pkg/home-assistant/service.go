package home_assistant

import (
	"errors"
)

var ErrNoSensor = errors.New("no sensor")

type service struct {
	client *client
}

var _ Service = &service{}

func NewService(c *client) *service {
	return &service{client: c}
}

func (s *service) GetSensorState(sensor string) (*SensorRawState, error) {
	sensors, err := s.client.getSensorsState()
	if err != nil {
		return nil, err
	}

	for _, state := range sensors {
		if state.Attributes.DeviceClass == sensor {
			return &state, nil
		}
	}

	return nil, ErrNoSensor
}
