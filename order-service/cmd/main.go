package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"beer/order-service/internal/repository"
	"beer/order-service/internal/server"
	"beer/order-service/internal/usecase"
	"beer/pkg/grpcmiddleware"
	"beer/pkg/pgprovider"
	"beer/proto/order"

	"google.golang.org/grpc"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://beer_user:beer_password@localhost:5432/beer?sslmode=disable"
	}
	listenAddr := os.Getenv("GRPC_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":50051"
	}

	pool, err := pgprovider.NewPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	uc := usecase.NewUsecase(repo)
	srv := server.NewServer(uc)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.UnaryServerLogger()),
	)
	orderpb.RegisterOrderServiceServer(grpcServer, srv)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	go func() {
		log.Printf("order-service gRPC listening on %s", listenAddr)
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
