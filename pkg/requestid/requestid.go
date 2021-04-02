package requestid

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

func (c contextKey) String() string {
	return "RequestID_" + string(c)
}

var (
	contextKeyRequestID = contextKey("request-id")
)

const (
	unableToCorrelate = "UNABLE_TO_CORRELATE"
)

func New(parent context.Context, rid string) context.Context {
	if rid == "" {
		return Context(parent)
	}
	return context.WithValue(parent, contextKeyRequestID, rid)
}

func Context(parent context.Context) context.Context {
	return context.WithValue(parent, contextKeyRequestID, NewRequestID())
}

func NewRequestID() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		return unableToCorrelate
	}
	return "|:" + uid.String()
}

func From(ctx context.Context) (rid string, exists bool) {
	rid, ok := ctx.Value(contextKeyRequestID).(string)
	if !ok {
		return unableToCorrelate, false
	}
	return rid, true
}
