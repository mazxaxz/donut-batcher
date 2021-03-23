package rabbitmq

import (
	"context"
	"encoding/json"
	"os"

	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
	"github.com/mazxaxz/donut-batcher/pkg/requestid"
)

type Publisher interface {
	Publish(ctx context.Context, data interface{}, msgType string) error
}

type publisherContext struct {
	channel *amqp.Channel
	cfg     config.Publisher
}

func NewPublisher(ctx context.Context, c *Client, cfg config.Publisher) (Publisher, error) {
	p := publisherContext{
		cfg: cfg,
	}

	ch, err := c.connection.Channel()
	if err != nil {
		return nil, err
	}
	p.channel = ch
	go p.close(ctx)

	if err := p.channel.ExchangeDeclare(
		cfg.Exchange,
		cfg.Kind,
		true,
		false,
		false,
		false,
		nil); err != nil {

		return nil, err
	}

	if err := p.channel.QueueBind(cfg.Queue, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *publisherContext) Publish(ctx context.Context, data interface{}, msgType string) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rid, exists := requestid.From(ctx)
	if !exists {
		ctx = requestid.Context(ctx)
		rid, _ = requestid.From(ctx)
	}
	hostname, _ := os.Hostname()
	msg := amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: rid,
		Type:          msgType,
		AppId:         hostname,
		Body:          payload,
	}
	if err := c.channel.Publish(c.cfg.Exchange, c.cfg.RoutingKey, false, false, msg); err != nil {
		return err
	}
	return nil
}

func (c *publisherContext) close(ctx context.Context) {
	<-ctx.Done()
	func() { _ = c.channel.Close() }()
}
