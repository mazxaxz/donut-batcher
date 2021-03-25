package batch

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mazxaxz/donut-batcher/pkg/money"
)

type Status string

const (
	StatusUndispatched    = "undispatched"
	StatusReadyToDispatch = "ready-to-dispatch"
	StatusDispatched      = "dispatched"
)

type Batch struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID         string               `bson:"userId" json:"userId"`
	Amount         primitive.Decimal128 `bson:"amount", json:"amount"`
	Currency       money.Currency       `bson:"currency", json:"currency"`
	TransactionIDs []string             `bson:"transactionIds" json:"transactionIds"`
	Status         Status               `bson:"status" json:"status"`
	CreatedDate    time.Time            `bson:"createdDate" json:"createdDate"`
	UpdatedDate    time.Time            `bson:"updatedDate" json:"updatedDate"`
	DispatchedDate time.Time            `bson:"dispatchedDate,omitempty" json:"dispatchedDate"`
}

func NewBatch(userID string, currency money.Currency) Batch {
	defaultAmount, _ := primitive.ParseDecimal128("0")
	b := Batch{
		UserID:         userID,
		Amount:         defaultAmount,
		Currency:       currency,
		TransactionIDs: make([]string, 0),
		Status:         StatusUndispatched,
		CreatedDate:    time.Now().UTC(),
	}
	return b
}
