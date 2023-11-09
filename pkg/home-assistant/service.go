package home_assistant

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var ErrNoSensor = errors.New("no sensor")

type service struct {
	pathStates  string
	restyClient *resty.Client
	token       string
}

func NewService(url, token string) Service {
	return &service{
		pathStates:  fmt.Sprintf("%s/api%s", url, "/states"),
		restyClient: resty.New(),
		token:       token,
	}
}

func (s *service) GetSensorState(sensor string) (*SensorRawState, error) {
	sensors, err := s.GetSensorsState()
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

func (s *service) GetSensorsState() ([]SensorRawState, error) {
	var (
		result []SensorRawState
	)

	resp, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.token).
		SetResult(&result).
		Get(s.pathStates)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(string(resp.Body()))
	}

	return result, nil
}
