package dispatchmessagehandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/batch"
	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/message/dispatch"
)

type handlerContext struct {
	batchSvc batch.Service
	logger   *logrus.Logger
}

func New(bSvc batch.Service, l *logrus.Logger) *handlerContext {
	c := handlerContext{
		batchSvc: bSvc,
		logger:   l,
	}
	return &c
}

func (c *handlerContext) Handle(ctx context.Context, delivery amqp.Delivery) (bool, error) {
	switch delivery.Type {
	case dispatch.MessageTypeDispatch:
		var msg dispatch.Dispatch
		if err := json.Unmarshal(delivery.Body, &msg); err != nil {
			return true, errors.Wrap(err, fmt.Sprintf("Body: %s", string(delivery.Body[:])))
		}

		if err := c.batchSvc.Dispatch(ctx, msg.BatchID); err != nil {
			switch err {
			case batch.ErrNoBatchID:
				return true, err
			default:
				return false, err
			}
		}
		return true, nil
	default:
		c.logger.Warnf("unknown type: '%s'", delivery.Type)
		return true, rabbitmq.ErrUnknownMessageType
	}
}
