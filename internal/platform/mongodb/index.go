package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Indexer interface {
	Index(ctx context.Context)
}

func (c *clientContext) CreateIndex(ctx context.Context, collectionName string, spec mongo.IndexModel) error {
	collection := c.client.Database(c.db).Collection(collectionName)
	if _, err := collection.Indexes().CreateOne(ctx, spec); err != nil {
		return errors.Wrap(err, "could not create mongodb index")
	}
	return nil
}
