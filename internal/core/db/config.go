package db

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     int           `envconfig:"PORT" required:"true"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfigMust() *Config {
	var config Config

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		panic(fmt.Errorf("get Server config: %w", err))
	}

	return &config
}
