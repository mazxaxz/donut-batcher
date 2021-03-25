package batch

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOrg "go.mongodb.org/mongo-driver/mongo"

	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
	"github.com/mazxaxz/donut-batcher/pkg/money"
)

type BatchResult struct {
	ID     primitive.ObjectID
	Status Status
}

func (c *serviceContext) Batch(ctx context.Context, t transaction.Transaction) (BatchResult, error) {
	callback := func(sessCtx mongoOrg.SessionContext) (interface{}, error) {
		currency, err := money.CurrencyFrom(t.Currency)
		if err != nil {
			return nil, err
		}
		filter := bson.D{{"userId", t.UserID}, {"status", StatusUndispatched}, {"currency", currency}}
		result := c.mongo.FindOne(ctx, _collectionName, filter)
		if err := result.Err(); err != nil && !errors.Is(err, mongoOrg.ErrNoDocuments) {
			return nil, err
		}

		var b Batch
		if err := result.Decode(&b); err != nil {
			if !errors.Is(err, mongoOrg.ErrNoDocuments) {
				return nil, err
			} else {
				b = NewBatch(t.UserID, currency)
				insertResult, err := c.mongo.InsertOne(ctx, _collectionName, b)
				if err != nil {
					return nil, err
				}
				b.ID = insertResult.InsertedID.(primitive.ObjectID)
			}
		}
		investment, err := money.CalculateInvestment(t.Amount)
		if err != nil {
			return nil, err
		}
		amount, err := money.Add(b.Amount.String(), investment)
		if err != nil {
			return nil, err
		}

		if threshold, exists := c.threshold[b.Currency]; exists {
			exceeded, err := money.GreaterThanOrEqual(amount, threshold)
			if err != nil {
				return nil, err
			}
			if exceeded {
				b.Status = StatusReadyToDispatch
			}
		}

		b.UpdatedDate = time.Now().UTC()
		b.TransactionIDs = append(b.TransactionIDs, t.ID)
		b.Amount, err = primitive.ParseDecimal128(amount)
		if err != nil {
			return nil, err
		}

		filter = bson.D{{"_id", b.ID}}
		update := bson.D{
			{"$set", bson.D{
				{"amount", b.Amount},
				{"status", b.Status},
				{"transactionIds", b.TransactionIDs},
				{"updatedDate", b.UpdatedDate},
			}},
		}
		if err := c.mongo.UpdateOne(ctx, _collectionName, filter, update); err != nil {
			return nil, err
		}
		return BatchResult{ID: b.ID, Status: b.Status}, nil
	}

	result, err := c.mongo.WithinTransaction(ctx, callback)
	if err != nil {
		return BatchResult{}, err
	}
	if s, ok := result.(BatchResult); ok {
		return s, nil
	}
	return BatchResult{}, nil
}
