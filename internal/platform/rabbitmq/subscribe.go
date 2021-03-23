package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
)

type Callback func(ctx context.Context, delivery amqp.Delivery) error

func (c *Client) Subscribe(ctx context.Context, cfg config.Subscriber, cb Callback) {
	ch, err := c.connection.Channel()
	if err != nil {
		// TODO log
		return
	}
	defer func() { _ = ch.Close() }()

	durable := !cfg.Exclusive
	autoDelete := cfg.Exclusive
	if _, err := ch.QueueDeclare(cfg.Queue, durable, autoDelete, cfg.Exclusive, false, nil); err != nil {
		// TODO log
		return
	}

	msgs, err := ch.Consume(cfg.Queue, "TODO", false, cfg.Exclusive, false, false, nil)
	if err != nil {
		// TODO log
		return
	}

	hold := make(chan bool)
	go func() {
		<-ctx.Done()
		hold <- true
	}()

	go func() {
		for d := range msgs {
			// TODO request id
			if err := cb(ctx, d); err != nil {
				// TODO log
			}
		}
	}()
	<-hold
}
