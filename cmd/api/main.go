package main

import (
	"beer/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
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

	router.Run(":8080")
}
