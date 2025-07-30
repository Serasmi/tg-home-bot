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
	ShortName    string
	ID           string
	Class        string
	Icon         string
}

func (s Sensor) matched(state SensorRawState) bool {
	if s.Class != "" {
		return s.Class == state.Attributes.DeviceClass
	}

	return s.ID == state.EntityID
}

var Sensors = []Sensor{
	{
		Name:         "humidity",
		FriendlyName: "Humidity",
		ShortName:    "Humidity",
		Class:        "humidity",
		Icon:         "💧",
	},
	{
		Name:         "temperature",
		FriendlyName: "Temperature",
		ShortName:    "Temp",
		Class:        "temperature",
		Icon:         "🌡️",
	},
	{
		Name:         "rpi_cpu_temp",
		FriendlyName: "RPI Temperature",
		ShortName:    "RPI Temp",
		ID:           "sensor.sensor_rpi_cpu_temp",
		Icon:         "🌡️",
	},
	{
		Name:         "nas_server_state",
		FriendlyName: "NAS server state",
		ShortName:    "Server",
		ID:           "binary_sensor.192_168_2_7",
		Icon:         "🖥️",
	},
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
