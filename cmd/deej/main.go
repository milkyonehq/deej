package main

import (
	"context"
	"github.com/milkyonehq/deej/pkg/configuration"
	"github.com/milkyonehq/deej/pkg/logger"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	logger.Init("info")

	config, err := configuration.New()
	if err != nil {
		log.WithError(err).Fatalln("Failed to load configuration.")
	}
	log.Infoln("Configuration successfully loaded.")

	logger.Init(config.LogLevel)

	log.Infoln("Bot is starting...")

	log.Infoln("Bot is running. Press CTRL+C to exit.")

	<-ctx.Done()

	log.Infoln("Received shutdown signal. DeeJ bot is shutting down...")
}
