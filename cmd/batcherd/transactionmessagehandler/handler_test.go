package transactionmessagehandler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mazxaxz/donut-batcher/internal/batch"
	mockBatch "github.com/mazxaxz/donut-batcher/internal/batch/mock"
	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	mockRabbitmq "github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/mock"
	"github.com/mazxaxz/donut-batcher/pkg/message/dispatch"
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
	"github.com/mazxaxz/donut-batcher/pkg/money"
)

func TestHandle(t *testing.T) {
	t.Run("should return error, invalid message type", func(t *testing.T) {
		// arrange
		d := amqp.Delivery{Type: "invalid", Body: nil}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err, rabbitmq.ErrUnknownMessageType)
	})

	t.Run("should return error, invalid payload", func(t *testing.T) {
		// arrange
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: []byte("}invalid{{{")}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err)
	})

	t.Run("should return no transaction id error", func(t *testing.T) {
		// arrange
		msg := transaction.Transaction{
			UserID:   "11",
			Amount:   "1.11",
			Currency: "USD",
		}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Batch(gomock.Any(), msg).Return(batch.BatchResult{}, batch.ErrNoTransactionID)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err, batch.ErrNoTransactionID)
	})

	t.Run("should return no user id error", func(t *testing.T) {
		// arrange
		msg := transaction.Transaction{
			ID:       "1",
			Amount:   "1.11",
			Currency: "USD",
		}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Batch(gomock.Any(), msg).Return(batch.BatchResult{}, batch.ErrNoUserID)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err, batch.ErrNoUserID)
	})

	t.Run("should return invalid currency", func(t *testing.T) {
		// arrange
		msg := transaction.Transaction{
			ID:       "1",
			UserID:   "11",
			Amount:   "1.11",
			Currency: "AAAAAA",
		}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Batch(gomock.Any(), msg).Return(batch.BatchResult{}, money.ErrInvalidCurrencyCode)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err, money.ErrInvalidCurrencyCode)
	})

	t.Run("should return batch undispatched result", func(t *testing.T) {
		// arrange
		msg := transaction.Transaction{
			ID:       "1",
			UserID:   "11",
			Amount:   "1.11",
			Currency: "AAAAAA",
		}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Batch(gomock.Any(), msg).Return(batch.BatchResult{
			ID:     primitive.NewObjectID(),
			Status: batch.StatusUndispatched,
		}, nil)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.NoError(t, err)
	})

	t.Run("should successfully publish dispatch event", func(t *testing.T) {
		// arrange
		msg := transaction.Transaction{
			ID:       "1",
			UserID:   "11",
			Amount:   "1.11",
			Currency: "AAAAAA",
		}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls
		oid := primitive.NewObjectID()
		mockBatchSvc.EXPECT().Batch(gomock.Any(), msg).Return(batch.BatchResult{
			ID:     oid,
			Status: batch.StatusReadyToDispatch,
		}, nil)
		mockPublisher.EXPECT().Publish(gomock.Any(), dispatch.Dispatch{BatchID: oid.Hex()}, dispatch.MessageTypeDispatch).Return(nil)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.NoError(t, err)
	})

	t.Run("should ack message after publish failure", func(t *testing.T) {
		// arrange
		msg := transaction.Transaction{
			ID:       "1",
			UserID:   "11",
			Amount:   "1.11",
			Currency: "AAAAAA",
		}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: transaction.MessageTypeTransaction, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		mockPublisher := mockRabbitmq.NewMockPublisher(mockCtrl)
		handler := New(mockBatchSvc, mockPublisher, logrus.New())

		// expected calls
		oid := primitive.NewObjectID()
		mockBatchSvc.EXPECT().Batch(gomock.Any(), msg).Return(batch.BatchResult{
			ID:     oid,
			Status: batch.StatusReadyToDispatch,
		}, nil)
		mockPublisher.EXPECT().Publish(gomock.Any(), dispatch.Dispatch{
			BatchID: oid.Hex(),
		}, dispatch.MessageTypeDispatch).Return(errors.New("random error"))

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err)
	})
}
