package transactionmessagehandlers

import (
	"context"

	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
)

type handlerContext struct {
}

func New(dispatchPublisher rabbitmq.Publisher) *handlerContext {
	return nil
}

func (c *handlerContext) Handle(ctx context.Context, delivery amqp.Delivery) error {
	return nil
}
