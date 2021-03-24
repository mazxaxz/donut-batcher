package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mazxaxz/donut-batcher/internal/platform/mongodb/config"
)

type Client struct {
	client *mongo.Client
	db     string
	logger *logrus.Logger
}

func New(ctx context.Context, cfg config.Config, l *logrus.Logger) (*Client, error) {
	c := Client{
		db:     cfg.Database,
		logger: l,
	}
	timeout, _ := context.WithTimeout(ctx, 15*time.Second)
	client, err := mongo.Connect(timeout, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to mongodb")
	}
	c.client = client

	return &c, nil
}

func (c *Client) dispose(ctx context.Context) {
	<-ctx.Done()
	if err := c.client.Disconnect(ctx); err != nil {
		c.logger.Error(err)
	}
}
