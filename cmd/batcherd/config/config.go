package config

import (
	"github.com/Netflix/go-env"

	rabbitConfig "github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
)

type Config struct {
	MQClient                rabbitConfig.Client     `env:"MQ_CLIENT,required=true"`
	MQTransactionSubscriber rabbitConfig.Subscriber `env:"MQ_TRANSACTION_SUBSCRIBER,required=true"`
	MQDispatchSubscriber    rabbitConfig.Subscriber `env:"MQ_DISPATCH_SUBSCRIBER,required=true"`
	MQDispatchPublisher     rabbitConfig.Publisher  `env:"MQ_DISPATCH_PUBLISHER,required=true"`
}

func Init() (Config, error) {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
