package main

import (
	"beer/internal/handler"
	clientrepo "beer/internal/repository/client"
	positionrepo "beer/internal/repository/position"
	sellerrepo "beer/internal/repository/seller"
	"beer/pkg/pgprovider"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://beer_user:beer_password@localhost:5432/beer?sslmode=disable"
	}
	pool, err := pgprovider.NewPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()
	clientrepo.SetPool(pool)
	sellerrepo.SetPool(pool)
	positionrepo.SetPool(pool)

	router := gin.Default()

	router.GET("/positions", handler.GetPositions)
	router.GET("/positions/:id", handler.GetPositionByID)
	router.POST("/positions", handler.CreatePosition)
	router.PATCH("/positions/:id", handler.PatchPositionByID)
	router.DELETE("/positions/:id", handler.DeletePositionByID)

	router.GET("/clients", handler.GetClients)
	router.GET("/clients/:id", handler.GetClientByID)
	router.POST("/clients", handler.CreateClient)
	router.PATCH("/clients/:id", handler.PatchClientByID)
	router.DELETE("/clients/:id", handler.DeleteClientByID)

	router.GET("/sellers", handler.GetSellers)
	router.GET("/sellers/:id", handler.GetSellerByID)
	router.POST("/sellers", handler.CreateSeller)
	router.PATCH("/sellers/:id", handler.PatchSellerByID)
	router.DELETE("/sellers/:id", handler.DeleteSellerByID)

	router.GET("/admins", handler.GetSellers)
	router.GET("/admins/:id", handler.GetSellerByID)
	router.POST("/admins", handler.CreateSeller)
	router.PATCH("/admins/:id", handler.PatchSellerByID)
	router.DELETE("/admins/:id", handler.DeleteSellerByID)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
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
