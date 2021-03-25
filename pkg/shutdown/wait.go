package shutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Wait for termination signal then cancel the context
func Wait(cancel context.CancelFunc, l *logrus.Logger, srvs ...http.Server) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	<-termChan

	for i := range srvs {
		if err := srvs[i].Close(); err != nil {
			l.Warn(err)
		}
	}

	l.Info("Closing app...")
	cancel()
}
