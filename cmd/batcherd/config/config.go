package config

import (
	"github.com/Netflix/go-env"

	mongoConfig "github.com/mazxaxz/donut-batcher/internal/platform/mongodb/config"
	rabbitConfig "github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
	"github.com/mazxaxz/donut-batcher/pkg/logger"
	"github.com/mazxaxz/donut-batcher/pkg/rest"
)

type Config struct {
	HTTP                    rest.Config             `env:"HTTP"`
	ThresholdUSD            string                  `env:"THRESHOLD_USD,default=100"`
	MongoClient             mongoConfig.Config      `env:"MONGO_CLIENT,required=true"`
	MQClient                rabbitConfig.Config     `env:"MQ_CLIENT,required=true"`
	MQTransactionSubscriber rabbitConfig.Subscriber `env:"MQ_TRANSACTION_SUBSCRIBER,required=true"`
	MQTransactionPublisher  rabbitConfig.Publisher  `env:"MQ_TRANSACTION_PUBLISHER,required=true"`
	MQDispatchSubscriber    rabbitConfig.Subscriber `env:"MQ_DISPATCH_SUBSCRIBER,required=true"`
	MQDispatchPublisher     rabbitConfig.Publisher  `env:"MQ_DISPATCH_PUBLISHER,required=true"`
	Logger                  logger.Config           `env:"LOGGER"`
}

func Load() (Config, error) {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
