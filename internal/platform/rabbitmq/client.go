package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
)

type Client struct {
	connection *amqp.Connection
}

func NewClient(ctx context.Context, cfg config.Client) (*Client, error) {
	c := Client{}
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, err
	}
	c.connection = conn
	go c.close(ctx)

	return &c, err
}

func (c *Client) close(ctx context.Context) {
	<-ctx.Done()
	func() { _ = c.connection.Close() }()
}
