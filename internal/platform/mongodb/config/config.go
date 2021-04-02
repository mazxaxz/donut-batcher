package config

import "encoding/json"

type Config struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

func (c *Config) UnmarshalEnvironmentValue(data string) error {
	return json.Unmarshal([]byte(data), &c)
}
