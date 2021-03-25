package logger

import "time"

type Log struct {
	Hostname     string    `json:"hostname"`
	Severity     string    `json:"severity"`
	RequestID    string    `json:"request_id"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	Milliseconds int64     `json:"milliseconds,omitempty"`
}
