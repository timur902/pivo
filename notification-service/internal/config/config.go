package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	KafkaBrokers          []string `env:"KAFKA_BROKERS,required" envSeparator:","`
	KafkaTopicOrdersReady string   `env:"KAFKA_TOPIC_ORDERS_READY" envDefault:"orders.ready"`
	KafkaGroupID          string   `env:"KAFKA_GROUP_ID" envDefault:"notification-service"`
}

func Load() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
