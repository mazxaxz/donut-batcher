package dispatchmessagehandler

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/message/dispatch"
)

type handlerContext struct {
	logger *logrus.Logger
}

func New(l *logrus.Logger) *handlerContext {
	c := handlerContext{
		logger: l,
	}
	return &c
}

func (c *handlerContext) Handle(ctx context.Context, delivery amqp.Delivery) (bool, error) {
	switch delivery.Type {
	case dispatch.MessageTypeDispatch:
		return true, nil
	default:
		c.logger.Warnf("unknown type: '%s'", delivery.Type)
		return true, rabbitmq.ErrUnknownMessageType
	}
}
