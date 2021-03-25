package transactionmessagehandler

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
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
)

type handlerContext struct {
	batchSvc          batch.Service
	dispatchPublisher rabbitmq.Publisher
	logger            *logrus.Logger
}

func New(bSvc batch.Service, dispatchPublisher rabbitmq.Publisher, l *logrus.Logger) *handlerContext {
	c := handlerContext{
		batchSvc:          bSvc,
		dispatchPublisher: dispatchPublisher,
		logger:            l,
	}
	return &c
}

func (c *handlerContext) Handle(ctx context.Context, delivery amqp.Delivery) (bool, error) {
	switch delivery.Type {
	case transaction.MessageTypeTransaction:
		var msg transaction.Transaction
		if err := json.Unmarshal(delivery.Body, &msg); err != nil {
			return true, errors.Wrap(err, fmt.Sprintf("Body: %s", string(delivery.Body[:])))
		}

		result, err := c.batchSvc.Batch(ctx, msg)
		if err != nil {
			switch err {
			// TODO:
			}
		}
		if result.Status == batch.StatusReadyToDispatch {
			payload := dispatch.Dispatch{BatchID: result.ID.String()}
			err := c.dispatchPublisher.Publish(ctx, payload, dispatch.MessageTypeDispatch)
			if err != nil {
				return true, errors.Wrap(err, fmt.Sprintf("could not send dispatch event, BatchID: %s", result.ID.String()))
			}
		}
		return true, nil
	default:
		c.logger.Warnf("unknown type: '%s'", delivery.Type)
		return true, rabbitmq.ErrUnknownMessageType
	}
}
