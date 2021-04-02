package logger

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel   string `json:"log_level"`
	OutputType string `json:"output_type"`
}

func (c *Config) UnmarshalEnvironmentValue(data string) error {
	return json.Unmarshal([]byte(data), &c)
}

func Configure(l *logrus.Logger, config Config) error {
	if config.LogLevel != "" {
		level, err := logrus.ParseLevel(config.LogLevel)
		if err != nil {
			return err
		}
		l.SetLevel(level)
	}
	switch config.OutputType {
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		fallthrough
	default:
		l.SetFormatter(&logrus.TextFormatter{})
	}
	return nil
}
