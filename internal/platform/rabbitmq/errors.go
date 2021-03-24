package rabbitmq

import "errors"

var (
	ErrClientNotProvided  = errors.New("client is nil")
	ErrUnknownMessageType = errors.New("unknown message type")
)
