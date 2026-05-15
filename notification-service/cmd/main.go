package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"beer/notification-service/internal/config"
	"beer/notification-service/internal/consumer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	c := consumer.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopicOrdersReady, cfg.KafkaGroupID)
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.Run(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		log.Println("notification-service shutting down")
		cancel()
		<-errCh
	case err := <-errCh:
		if err != nil {
			log.Fatalf("consumer failed: %v", err)
		}
	}
}
