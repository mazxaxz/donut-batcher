package batch

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	mongoOrg "go.mongodb.org/mongo-driver/mongo"

	"github.com/mazxaxz/donut-batcher/internal/platform/mongodb"
	"github.com/mazxaxz/donut-batcher/pkg/banksdk"
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
	"github.com/mazxaxz/donut-batcher/pkg/money"
)

const (
	_collectionName = "batches"
)

type Service interface {
	mongodb.Indexer

	Paginate(ctx context.Context, limit, offset int, asc bool, status *Status) ([]Batch, error)
	Batch(ctx context.Context, t transaction.Transaction) (BatchResult, error)
	Dispatch(ctx context.Context, batchID string) error
}

type serviceContext struct {
	mongo     mongodb.Clienter
	bankSDK   banksdk.Clienter
	logger    *logrus.Logger
	threshold map[money.Currency]string
}

func New(mc mongodb.Clienter, bc banksdk.Clienter, l *logrus.Logger, threshold map[string]string) (Service, error) {
	c := serviceContext{
		mongo:     mc,
		bankSDK:   bc,
		logger:    l,
		threshold: make(map[money.Currency]string),
	}
	for k, v := range threshold {
		if v == "" {
			return nil, money.ErrZeroAmount
		}
		currency, err := money.CurrencyFrom(k)
		if err != nil {
			return nil, err
		}
		c.threshold[currency] = v
	}
	return &c, nil
}

func (c *serviceContext) Index(ctx context.Context) {
	timeout, _ := context.WithTimeout(ctx, 30*time.Second)

	indexes := []mongoOrg.IndexModel{
		{Keys: bson.D{{"status", 1}}},
		{Keys: bson.D{{"createdDate", -1}}},
		{Keys: bson.D{{"_id", 1}, {"status", 1}}},
		{Keys: bson.D{{"userId", 1}, {"status", 1}, {"currency", 1}}},
	}

	var wg sync.WaitGroup
	for _, idx := range indexes {
		wg.Add(1)
		go func(idx mongoOrg.IndexModel) {
			defer wg.Done()
			if err := c.mongo.CreateIndex(timeout, _collectionName, idx); err != nil {
				c.logger.Error(err)
			}
		}(idx)
	}
	wg.Wait()
}
