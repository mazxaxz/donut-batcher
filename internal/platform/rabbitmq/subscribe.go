package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq/config"
	"github.com/mazxaxz/donut-batcher/pkg/logger"
	"github.com/mazxaxz/donut-batcher/pkg/requestid"
)

type Callback func(ctx context.Context, delivery amqp.Delivery) (bool, error)

func (c *Client) Subscribe(ctx context.Context, cfg config.Subscriber, cb Callback) {
	hostname, _ := os.Hostname()
	ch, err := c.connection.Channel()
	if err != nil {
		entry := logger.Log{
			Hostname:  hostname,
			Severity:  logrus.ErrorLevel.String(),
			Message:   errors.Wrap(err, "could not initialize channel").Error(),
			Timestamp: time.Now().UTC(),
		}
		c.logger.Error(entry)
		return
	}
	defer func() { _ = ch.Close() }()

	durable := !cfg.Exclusive
	autoDelete := cfg.Exclusive
	if _, err := ch.QueueDeclare(cfg.Queue, durable, autoDelete, cfg.Exclusive, false, nil); err != nil {
		entry := logger.Log{
			Hostname:  hostname,
			Severity:  logrus.ErrorLevel.String(),
			Message:   errors.Wrap(err, "could not declare queue").Error(),
			Timestamp: time.Now().UTC(),
		}
		c.logger.Error(entry)
		return
	}

	consumerID := uuid.NewString()
	ch.Qos(cfg.PrefetchCount, 0, false)
	msgs, err := ch.Consume(cfg.Queue, fmt.Sprintf("%s_%s", hostname, consumerID), false, cfg.Exclusive, false, false, nil)
	if err != nil {
		entry := logger.Log{
			Hostname:  hostname,
			Severity:  logrus.ErrorLevel.String(),
			Message:   errors.Wrap(err, "could not attach consumer").Error(),
			Timestamp: time.Now().UTC(),
		}
		c.logger.Error(entry)
		return
	}

	hold := make(chan bool)
	go func() {
		<-ctx.Done()
		hold <- true
	}()

	go func() {
		for d := range msgs {
			if d.CorrelationId == "" {
				d.CorrelationId = requestid.NewRequestID()
			}
			ctx = requestid.New(ctx, d.CorrelationId)

			start := time.Now()
			ack, err := cb(ctx, d)
			elapsed := time.Since(start)
			if err != nil {
				entry := logger.Log{
					Hostname:     hostname,
					Severity:     logrus.ErrorLevel.String(),
					RequestID:    d.CorrelationId,
					Message:      err.Error(),
					Timestamp:    time.Now().UTC(),
					Milliseconds: elapsed.Milliseconds(),
				}
				c.logger.Error(entry)
			}
			if ack {
				if err := d.Ack(false); err != nil {
					entry := logger.Log{
						Hostname:     hostname,
						Severity:     logrus.ErrorLevel.String(),
						RequestID:    d.CorrelationId,
						Message:      err.Error(),
						Timestamp:    time.Now().UTC(),
						Milliseconds: elapsed.Milliseconds(),
					}
					c.logger.Error(entry)
				}
			} else {
				if err := d.Nack(false, true); err != nil {
					entry := logger.Log{
						Hostname:     hostname,
						Severity:     logrus.ErrorLevel.String(),
						RequestID:    d.CorrelationId,
						Message:      err.Error(),
						Timestamp:    time.Now().UTC(),
						Milliseconds: elapsed.Milliseconds(),
					}
					c.logger.Error(entry)
				}
			}
			entry := logger.Log{
				Hostname:     hostname,
				Severity:     logrus.InfoLevel.String(),
				RequestID:    d.CorrelationId,
				Message:      "Processing finished",
				Timestamp:    time.Now().UTC(),
				Milliseconds: elapsed.Milliseconds(),
			}
			c.logger.Info(entry)
		}
	}()
	<-hold
}
