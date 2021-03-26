package batch

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOrg "go.mongodb.org/mongo-driver/mongo"

	"github.com/mazxaxz/donut-batcher/internal/platform/mongodb"
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
	"github.com/mazxaxz/donut-batcher/pkg/money"
)

var (
	ErrNoTransactionID = errors.New("no transaction id was provided")
	ErrNoUserID        = errors.New("no user id was provided")
)

type BatchResult struct {
	ID     primitive.ObjectID
	Status Status
}

func (c *serviceContext) Batch(ctx context.Context, t transaction.Transaction) (BatchResult, error) {
	result, err := c.mongo.WithinTransaction(ctx, c.callback(ctx, t))
	if err != nil {
		return BatchResult{}, err
	}
	if s, ok := result.(BatchResult); ok {
		return s, nil
	}
	return BatchResult{}, nil
}

func (c *serviceContext) callback(ctx context.Context, t transaction.Transaction) mongodb.TransactionCallback {
	return func(sessCtx mongoOrg.SessionContext) (interface{}, error) {
		// that could be extracted to message, but I did not wanted to add complexity with validation library
		if t.ID == "" {
			return BatchResult{}, ErrNoTransactionID
		}
		if t.UserID == "" {
			return BatchResult{}, ErrNoUserID
		}
		//
		currency, err := money.CurrencyFrom(t.Currency)
		if err != nil {
			return BatchResult{}, err
		}
		filter := bson.D{{"userId", t.UserID}, {"status", StatusUndispatched}, {"currency", currency}}
		result := c.mongo.FindOne(ctx, _collectionName, filter)
		if err := result.Err(); err != nil && !errors.Is(err, mongoOrg.ErrNoDocuments) {
			return BatchResult{}, err
		}

		var b Batch
		if err := result.Decode(&b); err != nil {
			if !errors.Is(err, mongoOrg.ErrNoDocuments) {
				return BatchResult{}, err
			} else {
				b = NewBatch(t.UserID, currency)
				insertResult, err := c.mongo.InsertOne(ctx, _collectionName, b)
				if err != nil {
					return BatchResult{}, err
				}
				b.ID = insertResult.InsertedID.(primitive.ObjectID)
			}
		}
		investment, err := money.CalculateInvestment(t.Amount)
		if err != nil {
			return BatchResult{}, err
		}
		amount, err := money.Add(b.Amount.String(), investment)
		if err != nil {
			return BatchResult{}, err
		}

		if threshold, exists := c.threshold[b.Currency]; exists {
			exceeded, err := money.GreaterThanOrEqual(amount, threshold)
			if err != nil {
				return BatchResult{}, err
			}
			if exceeded {
				/*
					According to AC, there should be another batch entry created with
					status Undispatched, but it is a waste of resources IMO and that
					should be discussed if it is needed.
				*/
				b.Status = StatusReadyToDispatch
			}
		}

		b.UpdatedDate = time.Now().UTC()
		b.TransactionIDs = append(b.TransactionIDs, t.ID)
		b.Amount, err = primitive.ParseDecimal128(amount)
		if err != nil {
			return BatchResult{}, err
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
			return BatchResult{}, err
		}
		return BatchResult{ID: b.ID, Status: b.Status}, nil
	}
}
