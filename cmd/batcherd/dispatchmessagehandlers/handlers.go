package dispatchmessagehandlers

import (
	"context"

	"github.com/streadway/amqp"
)

type handlerContext struct {
}

func New() *handlerContext {
	return nil
}

func (c *handlerContext) Handle(ctx context.Context, delivery amqp.Delivery) error {
	return nil
}
