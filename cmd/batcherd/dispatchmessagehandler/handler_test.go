package dispatchmessagehandler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"

	"github.com/mazxaxz/donut-batcher/internal/batch"
	mockBatch "github.com/mazxaxz/donut-batcher/internal/batch/mock"
	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/message/dispatch"
)

func TestHandle(t *testing.T) {
	t.Run("should return error, invalid message type", func(t *testing.T) {
		// arrange
		d := amqp.Delivery{Type: "invalid", Body: nil}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		handler := New(mockBatchSvc, logrus.New())

		// expected calls

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err, rabbitmq.ErrUnknownMessageType)
	})

	t.Run("should return error, invalid payload", func(t *testing.T) {
		// arrange
		d := amqp.Delivery{Type: dispatch.MessageTypeDispatch, Body: []byte("}invalid{{{")}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		handler := New(mockBatchSvc, logrus.New())

		// expected calls

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err)
	})

	t.Run("should return error, dispatching resulted in error", func(t *testing.T) {
		// arrange
		msg := dispatch.Dispatch{BatchID: "11111"}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: dispatch.MessageTypeDispatch, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		handler := New(mockBatchSvc, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Dispatch(gomock.Any(), msg.BatchID).Return(errors.New("random error"))

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.False(t, ack)
		assert.Error(t, err)
	})

	t.Run("should return error, no batch id provided", func(t *testing.T) {
		// arrange
		msg := dispatch.Dispatch{}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: dispatch.MessageTypeDispatch, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		handler := New(mockBatchSvc, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Dispatch(gomock.Any(), msg.BatchID).Return(batch.ErrNoBatchID)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.Error(t, err, batch.ErrNoBatchID)
	})

	t.Run("should ack message and return no error", func(t *testing.T) {
		// arrange
		msg := dispatch.Dispatch{BatchID: "11111"}
		body, err := json.Marshal(msg)
		assert.NoError(t, err)
		d := amqp.Delivery{Type: dispatch.MessageTypeDispatch, Body: body}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBatchSvc := mockBatch.NewMockService(mockCtrl)
		handler := New(mockBatchSvc, logrus.New())

		// expected calls
		mockBatchSvc.EXPECT().Dispatch(gomock.Any(), msg.BatchID).Return(nil)

		// act
		ack, err := handler.Handle(context.Background(), d)

		// assert
		assert.True(t, ack)
		assert.NoError(t, err)
	})
}
