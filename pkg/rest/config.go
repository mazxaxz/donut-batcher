package rest

import "encoding/json"

type Config struct {
	Port int `json:"port"`
}

func (c *Config) UnmarshalEnvironmentValue(data string) error {
	return json.Unmarshal([]byte(data), &c)
}
