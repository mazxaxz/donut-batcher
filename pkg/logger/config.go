package logger

import "github.com/sirupsen/logrus"

type Config struct {
	LogLevel string `json:"log_level"`
}

func Configure(l *logrus.Logger, config Config) error {
	if config.LogLevel != "" {
		level, err := logrus.ParseLevel(config.LogLevel)
		if err != nil {
			return err
		}
		l.SetLevel(level)
	}
	return nil
}
