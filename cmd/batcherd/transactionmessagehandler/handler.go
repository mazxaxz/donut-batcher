package transactionmessagehandler

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
)

type handlerContext struct {
	dispatchPublisher rabbitmq.Publisher
	logger            *logrus.Logger
}

func New(dispatchPublisher rabbitmq.Publisher, l *logrus.Logger) *handlerContext {
	c := handlerContext{
		dispatchPublisher: dispatchPublisher,
		logger:            l,
	}
	return &c
}

func (c *handlerContext) Handle(ctx context.Context, delivery amqp.Delivery) (bool, error) {
	switch delivery.Type {
	case transaction.MessageTypeTransaction:
		return true, nil
	default:
		c.logger.Warnf("unknown type: '%s'", delivery.Type)
		return true, rabbitmq.ErrUnknownMessageType
	}
}
