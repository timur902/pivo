package main

import (
	"beer/beer-api/internal/config"
	"beer/beer-api/internal/handler"
	"beer/beer-api/internal/repository/client"
	"beer/beer-api/internal/repository/position"
	"beer/beer-api/internal/repository/seller"
	"beer/beer-api/internal/usecase/seller"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	pool, err := pgprovider.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	orderConn, orderClient, err := orderclient.Dial(cfg.OrderServiceAddr)
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
		Addr:              cfg.HTTPListenAddr,
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
