package server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address         string        `envconfig:"ADDRESS" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
}

func NewConfigMust() *Config {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		panic(fmt.Errorf("get Server config: %w", err))
	}

	return &config
}
