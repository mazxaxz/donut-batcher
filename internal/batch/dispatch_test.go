package batch

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOrg "go.mongodb.org/mongo-driver/mongo"

	mockMongodb "github.com/mazxaxz/donut-batcher/internal/platform/mongodb/mock"
	"github.com/mazxaxz/donut-batcher/pkg/banksdk"
	"github.com/mazxaxz/donut-batcher/pkg/money"
)

func TestDispatch(t *testing.T) {
	t.Run("should return no batch id error", func(t *testing.T) {
		// arrange
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svc, err := New(mockMongoClient, banksdk.New(), logrus.New(), map[string]string{})
		assert.NoError(t, err)

		// expected calls

		// act
		err = svc.Dispatch(context.Background(), "")

		// assert
		assert.Error(t, err, ErrNoBatchID)
	})

	t.Run("should return mongo object id parse error", func(t *testing.T) {
		// arrange
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svc, err := New(mockMongoClient, banksdk.New(), logrus.New(), map[string]string{})
		assert.NoError(t, err)

		// expected calls

		// act
		err = svc.Dispatch(context.Background(), "invalid")

		// assert
		assert.Error(t, err)
	})

	t.Run("should return find one error", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svc, err := New(mockMongoClient, banksdk.New(), logrus.New(), map[string]string{})
		assert.NoError(t, err)

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(mongoOrg.ErrClientDisconnected)
		filter := bson.D{{"_id", batchID}, {"status", StatusReadyToDispatch}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		// act
		err = svc.Dispatch(context.Background(), batchID.Hex())

		// assert
		assert.Error(t, err, mongoOrg.ErrClientDisconnected)
	})

	t.Run("should return nil, no batch with given id was found", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svc, err := New(mockMongoClient, banksdk.New(), logrus.New(), map[string]string{})
		assert.NoError(t, err)

		// expected calls
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(mongoOrg.ErrNoDocuments)
		singleResult.EXPECT().Decode(gomock.Any()).Return(mongoOrg.ErrNoDocuments)
		filter := bson.D{{"_id", batchID}, {"status", StatusReadyToDispatch}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		// act
		err = svc.Dispatch(context.Background(), batchID.Hex())

		// assert
		assert.NoError(t, err)
	})

	t.Run("should send money to bank and update dispatched batch", func(t *testing.T) {
		// arrange
		batchID := primitive.NewObjectID()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svc, err := New(mockMongoClient, banksdk.New(), logrus.New(), map[string]string{})
		assert.NoError(t, err)

		// expected calls
		/* there should be check if money was actually sent, but it's only a mock sdk */
		singleResult := mockMongodb.NewMockSingleResulter(mockCtrl)
		singleResult.EXPECT().Err().Return(nil)
		singleResult.EXPECT().Decode(gomock.Any()).Do(func(b *Batch) {
			b.ID = batchID
			b.UserID = "11"
			b.Amount, _ = primitive.ParseDecimal128("11.11")
			b.Currency = money.Currency("USD")
			b.Status = StatusReadyToDispatch
		}).Return(nil)
		filter := bson.D{{"_id", batchID}, {"status", StatusReadyToDispatch}}
		mockMongoClient.EXPECT().FindOne(gomock.Any(), _collectionName, filter).Return(singleResult)

		filter = bson.D{{"_id", batchID}}
		mockMongoClient.EXPECT().UpdateOne(gomock.Any(), _collectionName, filter, gomock.Any()).Return(nil)

		// act
		err = svc.Dispatch(context.Background(), batchID.Hex())

		// assert
		assert.NoError(t, err)
	})
}
