package home_assistant

import (
	"errors"
	"strings"
	"time"
)

const StateUnknown = "unknown"

var ErrNoSensor = errors.New("no sensor")

type Sensor struct {
	Name         string
	FriendlyName string
	ID           string
	Class        string
}

func (s Sensor) is(state SensorRawState) bool {
	if s.Class == SensorHumidity.Class || s.Class == SensorTemperature.Class {
		return s.Class == state.Attributes.DeviceClass
	}

	return s.ID == state.EntityID
}

var (
	SensorHumidity = Sensor{
		Name:         "humidity",
		FriendlyName: "Humidity",
		Class:        "humidity",
	}
	SensorTemperature = Sensor{
		Name:         "temperature",
		FriendlyName: "Temperature",
		Class:        "temperature",
	}
	SensorRPITemperature = Sensor{
		Name:         "rpi_cpu_temp",
		FriendlyName: "RPI Temperature",
		ID:           "sensor.sensor_rpi_cpu_temp",
	}
)

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

func newTime(t time.Time) Time {
	return Time{&t}
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
