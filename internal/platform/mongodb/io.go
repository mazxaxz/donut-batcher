package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *clientContext) Find(ctx context.Context, coll string, filter interface{}, opt *options.FindOptions) (*mongo.Cursor, error) {
	return c.client.Database(c.db).Collection(coll).Find(ctx, filter, opt)
}

func (c *clientContext) FindOne(ctx context.Context, coll string, filter interface{}) *mongo.SingleResult {
	return c.client.Database(c.db).Collection(coll).FindOne(ctx, filter)
}

func (c *clientContext) UpdateOne(ctx context.Context, coll string, filter, update interface{}) error {
	_, err := c.client.Database(c.db).Collection(coll).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientContext) InsertOne(ctx context.Context, coll string, doc interface{}) (*mongo.InsertOneResult, error) {
	return c.client.Database(c.db).Collection(coll).InsertOne(ctx, doc)
}
