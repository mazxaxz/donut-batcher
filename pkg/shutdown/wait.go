package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Wait for termination signal then cancel the context
func Wait(cancel context.CancelFunc) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	<-termChan

	// TODO LOG
	cancel()
}
