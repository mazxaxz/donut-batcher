package rabbitmq

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
)

type Client struct {
	connection *amqp.Connection
	logger     *logrus.Logger
}

func NewClient(ctx context.Context, cfg config.Config, l *logrus.Logger) (*Client, error) {
	c := Client{
		logger: l,
	}
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
