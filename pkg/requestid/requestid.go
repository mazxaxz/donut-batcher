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

func Context(parent context.Context) context.Context {
	var rid string
	uid, err := uuid.NewUUID()
	if err != nil {
		rid = unableToCorrelate
	} else {
		rid = "|:" + uid.String()
	}
	return context.WithValue(parent, contextKeyRequestID, rid)
}

func From(ctx context.Context) (rid string, exists bool) {
	rid, ok := ctx.Value(contextKeyRequestID).(string)
	if !ok {
		return unableToCorrelate, false
	}
	return rid, true
}
