package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mazxaxz/donut-batcher/internal/platform/mongodb/config"
	"github.com/mazxaxz/donut-batcher/pkg/logger"
)

type Clienter interface {
	Find(ctx context.Context, coll string, filter interface{}, opt *options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, coll string, filter interface{}) *mongo.SingleResult
	UpdateOne(ctx context.Context, coll string, filter, update interface{}) error
	InsertOne(ctx context.Context, coll string, doc interface{}) (*mongo.InsertOneResult, error)
	WithinTransaction(ctx context.Context, cb TransactionCallback) (result interface{}, err error)
	CreateIndex(ctx context.Context, collectionName string, spec mongo.IndexModel) error
}

type clientContext struct {
	client *mongo.Client
	db     string
	logger *logrus.Logger
}

func New(ctx context.Context, cfg config.Config, l *logrus.Logger) (Clienter, error) {
	c := clientContext{
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

func (c *clientContext) dispose(ctx context.Context) {
	<-ctx.Done()
	if err := c.client.Disconnect(ctx); err != nil {
		hostname, _ := os.Hostname()
		entry := logger.Log{
			Hostname:  hostname,
			Severity:  logrus.ErrorLevel.String(),
			Message:   err.Error(),
			Timestamp: time.Now().UTC(),
		}
		c.logger.Error(entry)
	}
}
