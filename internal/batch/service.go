package batch

import (
	"context"

	"github.com/mazxaxz/donut-batcher/internal/platform/mongodb"
)

type Service interface {
	mongodb.Indexer

	Batch(ctx context.Context) error
	Dispatch(ctx context.Context, ID string) error
}

type serviceContext struct {
	mongo *mongodb.Client
}

func New(mongoClient *mongodb.Client) Service {
	c := serviceContext{
		mongo: mongoClient,
	}
	return &c
}

func (c *serviceContext) Index(ctx context.Context) error {
	return nil
}

func (c *serviceContext) Batch(ctx context.Context) error {
	return nil
}

func (c *serviceContext) Dispatch(ctx context.Context, ID string) error {
	return nil
}
