package main

import (
	"beer/internal/handler"
	"beer/internal/repository/client"
	"beer/internal/repository/position"
	"beer/internal/repository/seller"
	"beer/internal/usecase/seller"
	"beer/pkg/orderclient"
	"beer/pkg/pgprovider"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://beer_user:beer_password@localhost:5432/beer?sslmode=disable"
	}
	orderServiceAddr := os.Getenv("ORDER_SERVICE_ADDR")
	if orderServiceAddr == "" {
		orderServiceAddr = "localhost:50051"
	}

	pool, err := pgprovider.NewPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	orderConn, orderClient, err := orderclient.Dial(orderServiceAddr)
	if err != nil {
		log.Fatalf("order-service dial failed: %v", err)
	}
	defer orderConn.Close()

	clientRepo := client.NewRepository(pool)
	sellerRepo := seller.NewRepository(pool)
	positionRepo := position.NewRepository(pool)
	sellerUC := sellerusecase.NewUsecase(sellerRepo)
	hdl := handler.NewHandler(clientRepo, positionRepo, sellerUC, orderClient)

	router := gin.Default()

	router.GET("/positions", hdl.GetPositions)
	router.GET("/positions/:id", hdl.GetPositionByID)
	router.POST("/positions", hdl.CreatePosition)
	router.PATCH("/positions/:id", hdl.PatchPositionByID)
	router.DELETE("/positions/:id", hdl.DeletePositionByID)

	router.GET("/clients", hdl.GetClients)
	router.GET("/clients/:id", hdl.GetClientByID)
	router.POST("/clients", hdl.CreateClient)
	router.PATCH("/clients/:id", hdl.PatchClientByID)
	router.DELETE("/clients/:id", hdl.DeleteClientByID)

	router.GET("/sellers", hdl.GetSellers)
	router.GET("/sellers/:id", hdl.GetSellerByID)
	router.POST("/sellers", hdl.CreateSeller)
	router.PATCH("/sellers/:id", hdl.PatchSellerByID)
	router.DELETE("/sellers/:id", hdl.DeleteSellerByID)

	router.GET("/admins", hdl.GetSellers)
	router.GET("/admins/:id", hdl.GetSellerByID)
	router.POST("/admins", hdl.CreateSeller)
	router.PATCH("/admins/:id", hdl.PatchSellerByID)
	router.DELETE("/admins/:id", hdl.DeleteSellerByID)

	router.POST("/orders", hdl.CreateOrder)
	router.GET("/orders", hdl.ListClientOrders)
	router.GET("/orders/:id", hdl.GetOrderByID)
	router.DELETE("/orders/:id", hdl.CancelOrderByClient)

	router.GET("/seller/orders", hdl.ListSellerOrders)
	router.PATCH("/seller/orders/:id/status", hdl.UpdateOrderStatusBySeller)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
}
