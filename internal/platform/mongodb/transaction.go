package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionCallback func(sessCtx mongo.SessionContext) (interface{}, error)

func (c *Client) Transaction(ctx context.Context, cb TransactionCallback) (result interface{}, err error) {
	session, err := c.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, cb)
}
