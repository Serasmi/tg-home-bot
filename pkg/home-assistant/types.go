package home_assistant

import (
	"strings"
	"time"
)

type Service interface {
	GetSensorState(string) (*SensorRawState, error)
}

type Attributes struct {
	DeviceClass       string `json:"device_class"`
	FriendlyName      string `json:"friendly_name"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
}

type SensorRawState struct {
	Attributes Attributes `json:"attributes"`
	State      string     `json:"state"`
	EntityID   string     `json:"entity_id"`
	ChangedAt  Time       `json:"last_changed"`
	UpdatedAt  Time       `json:"last_updated"`
}

type Time struct {
	*time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")

	if str == "null" || str == "" {
		return nil
	}

	pTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = &pTime

	return nil
}
