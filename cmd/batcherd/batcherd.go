package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/donut-batcher/cmd/batcherd/config"
	"github.com/mazxaxz/donut-batcher/cmd/batcherd/dispatchmessagehandler"
	"github.com/mazxaxz/donut-batcher/cmd/batcherd/transactionmessagehandler"
	"github.com/mazxaxz/donut-batcher/internal/batch"
	"github.com/mazxaxz/donut-batcher/internal/platform/mongodb"
	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/banksdk"
	"github.com/mazxaxz/donut-batcher/pkg/logger"
	"github.com/mazxaxz/donut-batcher/pkg/shutdown"
)

var log = logrus.New()

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Configure(log, cfg.Logger); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	rabbitClient, err := rabbitmq.NewClient(ctx, cfg.MQClient, log)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient, err := mongodb.New(ctx, cfg.MongoClient, log)
	if err != nil {
		log.Fatal(err)
	}

	bankSDK := banksdk.New()

	dispatchPublisher, err := rabbitmq.NewPublisher(ctx, rabbitClient, cfg.MQDispatchPublisher)
	if err != nil {
		log.Fatal(err)
	}

	thresholds := map[string]string{"USD": cfg.ThresholdUSD}
	batchService, err := batch.New(mongoClient, bankSDK, thresholds)
	if err != nil {
		log.Fatal(err)
	}

	transactionMessageHandler := transactionmessagehandler.New(batchService, dispatchPublisher, log)
	dispatchMessageHandler := dispatchmessagehandler.New(log)

	go rabbitClient.Subscribe(ctx, cfg.MQTransactionSubscriber, transactionMessageHandler.Handle)
	go rabbitClient.Subscribe(ctx, cfg.MQDispatchSubscriber, dispatchMessageHandler.Handle)

	go index(ctx, batchService)

	/*
		Dispatching leftovers every n hours using cron should be handled here.
		I see no point of doing that in here, it's just a function invocation.
	*/

	log.Info("App started...")
	shutdown.Wait(cancel, log)
}

func index(ctx context.Context, indexers ...mongodb.Indexer) {
	for _, idx := range indexers {
		if err := idx.Index(ctx); err != nil {
			log.Error(err)
		}
	}
}
