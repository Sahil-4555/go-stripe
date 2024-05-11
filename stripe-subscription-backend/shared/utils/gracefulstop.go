package utils

import (
	"context"
	"os"
	"os/signal"
	"stripe-subscription/shared/log"
	"syscall"
	"time"
)

func GracefulStop(log log.ILogger, callback func(context.Context) error) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulStop

	log.Info("", "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := callback(ctx); err != nil {
		log.Fatal("", "Server forced to shutdown:", err)
	}

	log.Info("", "Server exiting")
}
