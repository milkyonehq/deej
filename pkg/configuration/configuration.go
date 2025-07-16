package configuration

import (
	"errors"
	"fmt"
	"github.com/Netflix/go-env"
)

var (
	ErrLoadConfiguration = errors.New("failed to load configuration")
)

type Config struct {
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN,required=true"`
	LogLevel        string `env:"LOG_LEVEL,default=info"`
}

func New() (*Config, error) {
	var config Config

	if _, err := env.UnmarshalFromEnviron(&config); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrLoadConfiguration, err)
	}

	return &config, nil
}
