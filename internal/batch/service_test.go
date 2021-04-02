package batch

import (
	"context"
	"github.com/mazxaxz/donut-batcher/pkg/banksdk"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	mockMongodb "github.com/mazxaxz/donut-batcher/internal/platform/mongodb/mock"
)

func TestIndex(t *testing.T) {
	t.Run("should call for index creation", func(t *testing.T) {
		// arrange
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockMongoClient := mockMongodb.NewMockClienter(mockCtrl)
		svc, err := New(mockMongoClient, banksdk.New(), logrus.New(), map[string]string{})
		assert.NoError(t, err)

		// expected calls
		idx := mongo.IndexModel{Keys: bson.D{{"status", 1}}}
		mockMongoClient.EXPECT().CreateIndex(gomock.Any(), _collectionName, idx).Return(nil)
		idx = mongo.IndexModel{Keys: bson.D{{"createdDate", -1}}}
		mockMongoClient.EXPECT().CreateIndex(gomock.Any(), _collectionName, idx).Return(nil)
		idx = mongo.IndexModel{Keys: bson.D{{"_id", 1}, {"status", 1}}}
		mockMongoClient.EXPECT().CreateIndex(gomock.Any(), _collectionName, idx).Return(nil)
		idx = mongo.IndexModel{Keys: bson.D{{"userId", 1}, {"status", 1}, {"currency", 1}}}
		mockMongoClient.EXPECT().CreateIndex(gomock.Any(), _collectionName, idx).Return(nil)

		// act
		svc.Index(context.Background())

		// assert
	})
}
