package batch

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	Batch(ctx context.Context, t transaction.Transaction) (BatchResult, error)
	Dispatch(ctx context.Context, ID string) error
}

type serviceContext struct {
	mongo     *mongodb.Client
	bankSDK   banksdk.Client
	threshold map[money.Currency]string
}

func New(mc *mongodb.Client, bc banksdk.Client, threshold map[string]string) (Service, error) {
	c := serviceContext{
		mongo:     mc,
		bankSDK:   bc,
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

func (c *serviceContext) Index(ctx context.Context) error {
	timeout, _ := context.WithTimeout(ctx, 30*time.Second)

	idx := mongoOrg.IndexModel{Keys: bson.D{{"userId", 1}, {"status", 1}, {"currency", 1}}}
	if err := c.mongo.CreateIndex(timeout, _collectionName, idx); err != nil {
		return err
	}
	return nil
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

		batchResult := BatchResult{ID: b.ID, Status: b.Status}
		return batchResult, nil
	}

	result, err := c.mongo.Transaction(ctx, callback)
	if err != nil {
		return BatchResult{}, err
	}
	if s, ok := result.(BatchResult); ok {
		return s, nil
	}
	return BatchResult{}, nil
}

func (c *serviceContext) Dispatch(ctx context.Context, ID string) error {
	return nil
}
