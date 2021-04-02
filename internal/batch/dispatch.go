package batch

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOrg "go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoBatchID = errors.New("no batch id was provided")
)

func (c *serviceContext) Dispatch(ctx context.Context, batchID string) error {
	if batchID == "" {
		return ErrNoBatchID
	}
	ID, err := primitive.ObjectIDFromHex(batchID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", ID}, {"status", StatusReadyToDispatch}}
	result := c.mongo.FindOne(ctx, _collectionName, filter)
	if err := result.Err(); err != nil && !errors.Is(err, mongoOrg.ErrNoDocuments) {
		return err
	}

	var b Batch
	if err := result.Decode(&b); err != nil {
		switch err {
		case mongoOrg.ErrNoDocuments:
			return nil
		default:
			return err
		}
	}
	if err := c.bankSDK.Send(ctx, b.UserID, b.Amount.String(), b.Currency.String()); err != nil {
		return err
	}

	b.Status = StatusDispatched
	b.UpdatedDate = time.Now().UTC()
	b.DispatchedDate = time.Now().UTC()

	filter = bson.D{{"_id", b.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"status", b.Status},
			{"updatedDate", b.UpdatedDate},
			{"dispatchedDate", b.DispatchedDate},
		}},
	}
	return c.mongo.UpdateOne(ctx, _collectionName, filter, update)
}
