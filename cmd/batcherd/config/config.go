package config

import (
	"github.com/Netflix/go-env"

	mongoConfig "github.com/mazxaxz/donut-batcher/internal/platform/mongodb/config"
	rabbitConfig "github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
	"github.com/mazxaxz/donut-batcher/pkg/logger"
)

type Config struct {
	Logger                  logger.Config           `env:"LOGGER"`
	MongoClient             mongoConfig.Config      `env:"MONGO_CLIENT,required=true"`
	MQClient                rabbitConfig.Client     `env:"MQ_CLIENT,required=true"`
	MQTransactionSubscriber rabbitConfig.Subscriber `env:"MQ_TRANSACTION_SUBSCRIBER,required=true"`
	MQDispatchSubscriber    rabbitConfig.Subscriber `env:"MQ_DISPATCH_SUBSCRIBER,required=true"`
	MQDispatchPublisher     rabbitConfig.Publisher  `env:"MQ_DISPATCH_PUBLISHER,required=true"`
}

func Load() (Config, error) {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
