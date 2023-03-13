package home_assistant

type Service interface {
	GetSensorState(string) (*SensorState, error)
}

type SensorState struct {
	Value string
}
