package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL      string `env:"DATABASE_URL,required"`
	HTTPListenAddr   string `env:"HTTP_LISTEN_ADDR" envDefault:":8080"`
	OrderServiceAddr string `env:"ORDER_SERVICE_ADDR,required"`
}

func Load() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
