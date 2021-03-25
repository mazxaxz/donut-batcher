package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Wait for termination signal then cancel the context
func Wait(cancel context.CancelFunc, l *logrus.Logger) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	<-termChan

	logrus.Info("Closing app...")
	cancel()
}
