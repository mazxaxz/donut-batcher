package config

import "encoding/json"

type Config struct {
	URI string `json:"uri"`
}

func (c *Config) UnmarshalEnvironmentValue(data string) error {
	return json.Unmarshal([]byte(data), &c)
}

type Subscriber struct {
	Queue         string `json:"queue"`
	Exclusive     bool   `json:"exclusive"`
	PrefetchCount int    `json:"prefetch_count"`
}

func (s *Subscriber) UnmarshalEnvironmentValue(data string) error {
	return json.Unmarshal([]byte(data), &s)
}

type Publisher struct {
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
	Queue      string `json:"queue"`
	Kind       string `json:"kind"`
}

func (p *Publisher) UnmarshalEnvironmentValue(data string) error {
	return json.Unmarshal([]byte(data), &p)
}
