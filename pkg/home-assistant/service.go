package home_assistant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Service interface {
	GetSensorState(context.Context, Sensor) (*SensorRawState, error)
	GetSensorsState(context.Context, ...Sensor) ([]SensorRawState, error)
}

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

func (s *service) GetSensorState(ctx context.Context, sensor Sensor) (*SensorRawState, error) {
	states, err := s.GetSensorsState(ctx, sensor)
	if err != nil {
		return nil, err
	}

	if len(states) == 0 {
		return nil, ErrNoSensor
	}

	return &states[0], nil
}

func (s *service) GetSensorsState(ctx context.Context, sensors ...Sensor) ([]SensorRawState, error) {
	states, err := s.fetchSensorsState(ctx)
	if err != nil {
		return nil, err
	}

	if len(sensors) == 0 {
		return states, nil
	}

	res := make([]SensorRawState, len(sensors))

loop:
	for i := range sensors {
		for j := range states {
			if sensors[i].matched(states[j]) {
				res[i] = states[j]
				continue loop
			}
		}

		res[i] = SensorRawState{
			Attributes: Attributes{},
			State:      StateUnknown,
			ChangedAt:  newTime(time.Now()),
			UpdatedAt:  newTime(time.Now()),
		}
	}

	return res, nil
}

func (s *service) fetchSensorsState(ctx context.Context) ([]SensorRawState, error) {
	var (
		result []SensorRawState
	)

	resp, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.token).
		SetContext(ctx).
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
