package home_assistant

import "fmt"

type service struct{}

var _ Service = &service{}

func NewService() *service {
	return &service{}
}

func (s *service) GetSensorState(sensor string) (*SensorState, error) {
	// TODO: implement method
	value := fmt.Sprintf("No value for '%s' sensor 😢", sensor)

	return &SensorState{Value: value}, nil
}
