package main

import (
	"context"

	"github.com/mazxaxz/donut-batcher/cmd/batcherd/config"
	"github.com/mazxaxz/donut-batcher/cmd/batcherd/dispatchmessagehandlers"
	"github.com/mazxaxz/donut-batcher/cmd/batcherd/transactionmessagehandlers"
	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/shutdown"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		// TODO panic
	}
	ctx, cancel := context.WithCancel(context.Background())

	rabbitClient, err := rabbitmq.NewClient(ctx, cfg.MQClient)
	if err != nil {
		// TODO panic
	}

	dispatchPublisher, err := rabbitmq.NewPublisher(ctx, rabbitClient, cfg.MQDispatchPublisher)
	if err != nil {
		// TODO panic
	}

	transactionMessageHandler := transactionmessagehandlers.New(dispatchPublisher)
	dispatchMessageHandler := dispatchmessagehandlers.New()

	go rabbitClient.Subscribe(ctx, cfg.MQTransactionSubscriber, transactionMessageHandler.Handle)
	go rabbitClient.Subscribe(ctx, cfg.MQDispatchSubscriber, dispatchMessageHandler.Handle)

	// TODO remember about timing fallback

	shutdown.Wait(cancel)
}
