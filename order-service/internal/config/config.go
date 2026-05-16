package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL           string   `env:"DATABASE_URL,required"`
	GRPCListenAddr        string   `env:"GRPC_LISTEN_ADDR" envDefault:":50051"`
	KafkaBrokers          []string `env:"KAFKA_BROKERS,required" envSeparator:","`
	KafkaTopicOrdersReady string   `env:"KAFKA_TOPIC_ORDERS_READY" envDefault:"orders.ready"`
}

func Load() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
