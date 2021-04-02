package batch

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Disclaimer pointer to status parameter, is because it is faster that way
// USE POINTERS ONLY WHEN YOU NEED TO
func (c *serviceContext) Paginate(ctx context.Context, limit, offset int, asc bool, status *Status) ([]Batch, error) {
	filter := bson.M{}
	if status != nil {
		filter["status"] = *status
	}

	sort := -1
	if asc {
		sort = 1
	}
	opt := options.Find().SetSort(bson.M{"createdDate": sort}).SetSkip(int64(offset)).SetLimit(int64(limit))

	cursor, err := c.mongo.Find(ctx, _collectionName, filter, opt)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var batches []Batch
	if err := cursor.All(ctx, &batches); err != nil {
		return nil, err
	}
	if batches == nil {
		batches = make([]Batch, 0)
	}
	return batches, nil
}
