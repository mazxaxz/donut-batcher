package rabbitmq

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
	"github.com/mazxaxz/donut-batcher/pkg/requestid"
)

type Callback func(ctx context.Context, delivery amqp.Delivery) (bool, error)

func (c *Client) Subscribe(ctx context.Context, cfg config.Subscriber, cb Callback) {
	ch, err := c.connection.Channel()
	if err != nil {
		c.logger.Error(err)
		return
	}
	defer func() { _ = ch.Close() }()

	durable := !cfg.Exclusive
	autoDelete := cfg.Exclusive
	if _, err := ch.QueueDeclare(cfg.Queue, durable, autoDelete, cfg.Exclusive, false, nil); err != nil {
		c.logger.Error(err)
		return
	}

	hostname, _ := os.Hostname()
	uid := uuid.NewString()
	msgs, err := ch.Consume(cfg.Queue, fmt.Sprintf("%s_%s", hostname, uid), false, cfg.Exclusive, false, false, nil)
	if err != nil {
		c.logger.Error(err)
		return
	}

	hold := make(chan bool)
	go func() {
		<-ctx.Done()
		hold <- true
	}()

	go func() {
		for d := range msgs {
			if d.CorrelationId == "" {
				d.CorrelationId = requestid.NewRequestID()
			}
			ctx = requestid.New(ctx, d.CorrelationId)

			ack, err := cb(ctx, d)
			if err != nil {
				c.logger.Error(err)
			}
			if ack {
				if err := d.Ack(false); err != nil {
					c.logger.Error(err)
				}
			} else {
				if err := d.Nack(false, true); err != nil {
					c.logger.Error(err)
				}
			}
		}
	}()
	<-hold
}
