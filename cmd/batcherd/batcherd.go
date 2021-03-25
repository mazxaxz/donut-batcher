package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/donut-batcher/cmd/batcherd/config"
	"github.com/mazxaxz/donut-batcher/cmd/batcherd/dispatchmessagehandler"
	"github.com/mazxaxz/donut-batcher/cmd/batcherd/transactionhttphandler"
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

	// Clients
	rabbitClient, err := rabbitmq.NewClient(ctx, cfg.MQClient, log)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient, err := mongodb.New(ctx, cfg.MongoClient, log)
	if err != nil {
		log.Fatal(err)
	}

	// Services/Publishers
	bankSDK := banksdk.New()

	transactionPublisher, err := rabbitmq.NewPublisher(ctx, rabbitClient, cfg.MQTransactionPublisher)
	if err != nil {
		log.Fatal(err)
	}
	dispatchPublisher, err := rabbitmq.NewPublisher(ctx, rabbitClient, cfg.MQDispatchPublisher)
	if err != nil {
		log.Fatal(err)
	}

	thresholds := map[string]string{"USD": cfg.ThresholdUSD}
	batchService, err := batch.New(mongoClient, bankSDK, log, thresholds)
	if err != nil {
		log.Fatal(err)
	}

	go index(ctx, batchService)

	// Message handlers
	transactionMessageHandler := transactionmessagehandler.New(batchService, dispatchPublisher, log)
	dispatchMessageHandler := dispatchmessagehandler.New(batchService, log)

	go rabbitClient.Subscribe(ctx, cfg.MQTransactionSubscriber, transactionMessageHandler.Handle)
	go rabbitClient.Subscribe(ctx, cfg.MQDispatchSubscriber, dispatchMessageHandler.Handle)

	/*
		Dispatching leftovers every n hours using cron should be handled here.
		I see no point of doing that in here, it's just a function invocation.
	*/

	// HTTP Handlers
	transactionHTTPHandler := transactionhttphandler.New(batchService, transactionPublisher, log)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      setupRouting(transactionHTTPHandler),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	go func(ctx context.Context, srv http.Server) {
		log.Info(fmt.Sprintf("Starting server on port: %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil {
			log.Info("Closing server...")
		}
	}(ctx, srv)

	shutdown.Wait(cancel, log)
}

func index(ctx context.Context, indexers ...mongodb.Indexer) {
	for _, idx := range indexers {
		idx.Index(ctx)
	}
}
