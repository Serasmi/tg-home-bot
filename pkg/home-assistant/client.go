package home_assistant

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	statesPath = "/states"
)

type client struct {
	restyClient *resty.Client
	token       string
	url         string
}

func NewClient(url, token string) *client {
	return &client{restyClient: resty.New(), token: token, url: url}
}

func (c client) getSensorsState() ([]SensorRawState, error) {
	resp, err := c.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.token).
		Get(c.withPath(statesPath))
	if err != nil {
		return nil, err
	}

	var result []SensorRawState

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c client) withPath(path string) string {
	return fmt.Sprintf("%s/api%s", c.url, path)
}
