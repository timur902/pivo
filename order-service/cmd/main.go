package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"beer/order-service/internal/config"
	"beer/order-service/internal/events"
	"beer/order-service/internal/repository"
	"beer/order-service/internal/server"
	"beer/order-service/internal/usecase"
	"beer/pkg/grpcmiddleware"
	"beer/pkg/pgprovider"
	"beer/proto/order"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	pool, err := pgprovider.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	producer := events.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopicOrdersReady)
	defer producer.Close()

	repo := repository.NewRepository(pool)
	uc := usecase.NewUsecase(repo, producer)
	srv := server.NewServer(uc)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.UnaryServerLogger()),
	)
	orderpb.RegisterOrderServiceServer(grpcServer, srv)

	listener, err := net.Listen("tcp", cfg.GRPCListenAddr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	go func() {
		log.Printf("order-service gRPC listening on %s", cfg.GRPCListenAddr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("grpc server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("order-service shutting down")
	grpcServer.GracefulStop()
}
