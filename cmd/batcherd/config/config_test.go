package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should assign all values", func(t *testing.T) {
		// arrange
		os.Clearenv()
		os.Setenv("HTTP", "{\"port\":8085}")
		os.Setenv("THRESHOLD_USD", "100")
		os.Setenv("MONGO_CLIENT", "{\"uri\":\"mongodb://root:password@mongo:27017/admin\",\"database\":\"donut\"}")
		os.Setenv("MQ_CLIENT", "{\"uri\":\"amqp://user:secret@rabbit:5672/\"}")
		os.Setenv("MQ_TRANSACTION_SUBSCRIBER", "{\"queue\":\"Donut.Q.Transaction\",\"prefetch_count\":10}")
		os.Setenv("MQ_TRANSACTION_PUBLISHER", "{\"exchange\":\"Donut.T.Topic\",\"queue\":\"Donut.Q.Transaction\",\"routing_key\":\"Donut.K.Transaction\",\"kind\":\"topic\"}")
		os.Setenv("MQ_DISPATCH_SUBSCRIBER", "{\"queue\":\"Donut.Q.Dispatch\",\"prefetch_count\":10}")
		os.Setenv("MQ_DISPATCH_PUBLISHER", "{\"exchange\":\"Donut.T.Topic\",\"queue\":\"Donut.Q.Dispatch\",\"routing_key\":\"Donut.K.Dispatch\",\"kind\":\"topic\"}")
		os.Setenv("LOGGER", "{\"log_level\":\"info\",\"output_type\":\"json\"}")

		// act
		result, err := Load()

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 8085, result.HTTP.Port)
		assert.Equal(t, "100", result.ThresholdUSD)
		assert.Equal(t, "mongodb://root:password@mongo:27017/admin", result.MongoClient.URI)
		assert.Equal(t, "donut", result.MongoClient.Database)
		assert.Equal(t, "amqp://user:secret@rabbit:5672/", result.MQClient.URI)

		assert.Equal(t, 10, result.MQTransactionSubscriber.PrefetchCount)
		assert.Equal(t, "Donut.Q.Transaction", result.MQTransactionSubscriber.Queue)
		assert.Equal(t, false, result.MQTransactionSubscriber.Exclusive)
		assert.Equal(t, "Donut.Q.Transaction", result.MQTransactionPublisher.Queue)
		assert.Equal(t, "Donut.K.Transaction", result.MQTransactionPublisher.RoutingKey)
		assert.Equal(t, "Donut.T.Topic", result.MQTransactionPublisher.Exchange)
		assert.Equal(t, "topic", result.MQTransactionPublisher.Kind)

		assert.Equal(t, 10, result.MQDispatchSubscriber.PrefetchCount)
		assert.Equal(t, "Donut.Q.Dispatch", result.MQDispatchSubscriber.Queue)
		assert.Equal(t, false, result.MQDispatchSubscriber.Exclusive)
		assert.Equal(t, "Donut.Q.Dispatch", result.MQDispatchPublisher.Queue)
		assert.Equal(t, "Donut.K.Dispatch", result.MQDispatchPublisher.RoutingKey)
		assert.Equal(t, "Donut.T.Topic", result.MQDispatchPublisher.Exchange)
		assert.Equal(t, "topic", result.MQDispatchPublisher.Kind)

		assert.Equal(t, "info", result.Logger.LogLevel)
		assert.Equal(t, "json", result.Logger.OutputType)
	})

	t.Run("should assign default not required values", func(t *testing.T) {
		// arrange
		os.Clearenv()
		os.Setenv("HTTP", "{\"port\":8085}")
		os.Setenv("MONGO_CLIENT", "{\"uri\":\"mongodb://root:password@mongo:27017/admin\",\"database\":\"donut\"}")
		os.Setenv("MQ_CLIENT", "{\"uri\":\"amqp://user:secret@rabbit:5672/\"}")
		os.Setenv("MQ_TRANSACTION_SUBSCRIBER", "{\"queue\":\"Donut.Q.Transaction\",\"prefetch_count\":10}")
		os.Setenv("MQ_TRANSACTION_PUBLISHER", "{\"exchange\":\"Donut.T.Topic\",\"queue\":\"Donut.Q.Transaction\",\"routing_key\":\"Donut.K.Transaction\",\"kind\":\"topic\"}")
		os.Setenv("MQ_DISPATCH_SUBSCRIBER", "{\"queue\":\"Donut.Q.Dispatch\",\"prefetch_count\":10}")
		os.Setenv("MQ_DISPATCH_PUBLISHER", "{\"exchange\":\"Donut.T.Topic\",\"queue\":\"Donut.Q.Dispatch\",\"routing_key\":\"Donut.K.Dispatch\",\"kind\":\"topic\"}")

		// act
		result, err := Load()

		// assert
		assert.NoError(t, err)
		assert.Equal(t, "100", result.ThresholdUSD)
		assert.Equal(t, "", result.Logger.LogLevel)
		assert.Equal(t, "", result.Logger.OutputType)
	})

	t.Run("should return error, no required fields specified", func(t *testing.T) {
		// arrange
		os.Clearenv()

		// act
		_, err := Load()

		// assert
		assert.Error(t, err)
	})
}
