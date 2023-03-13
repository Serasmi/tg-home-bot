package sensor

import ha "tg-home-bot/pkg/home-assistant"

type useCase struct {
	iotService ha.Service
}

var _ sensorUseCase = &useCase{}

func NewUseCase(s ha.Service) *useCase {
	return &useCase{
		iotService: s,
	}
}

func (u *useCase) SensorValue(sensor string) (string, error) {
	value, err := u.iotService.GetSensorState(sensor)
	if err != nil {
		return "", err
	}

	return value.Value, nil
}
