package main

import (
	"context"
	"github.com/milkyonehq/deej/pkg/configuration"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/bot"
	"github.com/milkyonehq/deej/pkg/logger"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func initProviders(providerRegistry *provider.Registry) {
	pvrs := []provider.Provider{
		provider.NewRaw(),
		provider.NewYoutube(),
	}
	var pvrNames []string
	for _, pvr := range pvrs {
		providerRegistry.Register(pvr)
		pvrNames = append(pvrNames, pvr.Name())
	}
	log.WithFields(log.Fields{
		"count":     len(pvrNames),
		"providers": strings.Join(pvrNames, ","),
	}).Infoln("Providers successfully registered.")
}

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

	playerRegistry := player.NewRegistry()
	providerRegistry := provider.NewRegistry()

	initProviders(providerRegistry)

	b, err := bot.New(config.DiscordBotToken, playerRegistry, providerRegistry)
	if err != nil {
		log.WithError(err).Fatalln("Failed to create bot.")
	}

	if err = b.Start(); err != nil {
		log.WithError(err).Fatalln("Failed to start bot.")
	}
	defer b.Stop()

	log.Infoln("Bot is running. Press CTRL+C to exit.")

	<-ctx.Done()

	log.Infoln("Received shutdown signal. DeeJ bot is shutting down...")
}
