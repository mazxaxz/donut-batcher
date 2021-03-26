package batch

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOrg "go.mongodb.org/mongo-driver/mongo"

	mockMongodb "github.com/mazxaxz/donut-batcher/internal/platform/mongodb/mock"
	"github.com/mazxaxz/donut-batcher/pkg/banksdk"
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
	"github.com/mazxaxz/donut-batcher/pkg/money"
)

func TestBatch(t *testing.T) {
	t.Run("should return mongo error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		mockMongoClient.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).Return(nil, mongoOrg.ErrClientDisconnected)

		// act
		result, err := svcCtx.Batch(context.Background(), give)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err, mongoOrg.ErrClientDisconnected)
	})

	t.Run("should return mongo error", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		give := transaction.Transaction{ID: batchID.Hex(), UserID: "11", Currency: "USD"}
		want := BatchResult{ID: batchID, Status: StatusUndispatched}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		mockMongoClient.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).Return(want, nil)

		// act
		result, err := svcCtx.Batch(context.Background(), give)

		// assert
		assert.Equal(t, want, result)
		assert.NoError(t, err)
	})
}

func TestCallback(t *testing.T) {
	t.Run("should return no transaction id error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{UserID: "11", Currency: "USD"}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err, ErrNoTransactionID)
	})

	t.Run("should return no user id error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{ID: "1", Currency: "USD"}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err, ErrNoUserID)
	})

	t.Run("should return invalid currency error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{ID: "1", UserID: "11", Currency: "invalid"}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err, money.ErrInvalidCurrencyCode)
	})

	t.Run("should return find one error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{ID: "1", UserID: "11", Currency: "USD"}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(mongoOrg.ErrClientDisconnected)
		filter := bson.D{{"userId", give.UserID}, {"status", StatusUndispatched}, {"currency", money.Currency(give.Currency)}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err, mongoOrg.ErrClientDisconnected)
	})

	t.Run("should return find one error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{ID: "1", UserID: "11", Currency: "USD"}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(mongoOrg.ErrClientDisconnected)
		filter := bson.D{{"userId", give.UserID}, {"status", StatusUndispatched}, {"currency", money.Currency(give.Currency)}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err, mongoOrg.ErrClientDisconnected)
	})

	t.Run("should return decode error", func(t *testing.T) {
		// arrange
		give := transaction.Transaction{ID: "1", UserID: "11", Currency: "USD"}
		want := BatchResult{}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(nil)
		singleResult.EXPECT().Decode(gomock.Any()).Return(errors.New("random error"))
		filter := bson.D{{"userId", give.UserID}, {"status", StatusUndispatched}, {"currency", money.Currency(give.Currency)}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.Error(t, err)
	})

	t.Run("should add transaction to the existing batch", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		give := transaction.Transaction{ID: "1", UserID: "11", Amount: "11.11", Currency: "USD"}
		want := BatchResult{ID: batchID, Status: StatusUndispatched}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(nil)
		singleResult.EXPECT().Decode(gomock.Any()).Do(func(b *Batch) {
			b.ID = batchID
			b.UserID = give.UserID
			b.Amount, _ = primitive.ParseDecimal128("0")
			b.Currency = money.Currency(give.Currency)
			b.Status = StatusUndispatched
		}).Return(nil)
		filter := bson.D{{"userId", give.UserID}, {"status", StatusUndispatched}, {"currency", money.Currency(give.Currency)}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		filter = bson.D{{"_id", batchID}}
		mockMongoClient.EXPECT().UpdateOne(gomock.Any(), _collectionName, filter, gomock.Any()).Return(nil)

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.NoError(t, err)
	})

	t.Run("should add transaction to the newly created batch", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		give := transaction.Transaction{ID: "1", UserID: "11", Amount: "11.11", Currency: "USD"}
		want := BatchResult{ID: batchID, Status: StatusUndispatched}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "100"},
		}

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(mongoOrg.ErrNoDocuments)
		singleResult.EXPECT().Decode(gomock.Any()).Return(mongoOrg.ErrNoDocuments)
		filter := bson.D{{"userId", give.UserID}, {"status", StatusUndispatched}, {"currency", money.Currency(give.Currency)}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		mockMongoClient.EXPECT().InsertOne(gomock.Any(), _collectionName, gomock.Any()).Return(&mongoOrg.InsertOneResult{InsertedID: batchID}, nil)

		filter = bson.D{{"_id", batchID}}
		mockMongoClient.EXPECT().UpdateOne(gomock.Any(), _collectionName, filter, gomock.Any()).Return(nil)

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.NoError(t, err)
	})

	t.Run("should return ready to dispatch", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		give := transaction.Transaction{ID: "1", UserID: "11", Amount: "11.11", Currency: "USD"}
		want := BatchResult{ID: batchID, Status: StatusReadyToDispatch}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svcCtx := serviceContext{
			mongo:     mockMongoClient,
			bankSDK:   banksdk.New(),
			logger:    logrus.New(),
			threshold: map[money.Currency]string{"USD": "0.5"},
		}

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(nil)
		singleResult.EXPECT().Decode(gomock.Any()).Do(func(b *Batch) {
			b.ID = batchID
			b.UserID = give.UserID
			b.Amount, _ = primitive.ParseDecimal128("0")
			b.Currency = money.Currency(give.Currency)
			b.Status = StatusUndispatched
		}).Return(nil)
		filter := bson.D{{"userId", give.UserID}, {"status", StatusUndispatched}, {"currency", money.Currency(give.Currency)}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		filter = bson.D{{"_id", batchID}}
		mockMongoClient.EXPECT().UpdateOne(gomock.Any(), _collectionName, filter, gomock.Any()).Return(nil)

		// act
		ctx := context.Background()
		result, err := svcCtx.callback(ctx, give)(nil)

		// assert
		assert.Equal(t, want, result)
		assert.NoError(t, err)
	})
}
