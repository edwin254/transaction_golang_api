package routes

import (
	"gapstack-api/internal/controller"
	"gapstack-api/internal/repository"
	"gapstack/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	repo := repository.NewTransactionRepository(db)
	svc := service.NewTransactionService(repo)
	ctrl := controller.NewTransactionController(svc)

	api := router.Group("/api/v1")
	tx := api.Group("/transactions")
	{
		tx.POST("/", ctrl.CreateTransaction)
		tx.GET("/:id", ctrl.GetTransactionByID)
		tx.GET("/", ctrl.ListTransactions)
		tx.PUT("/:id/status", ctrl.UpdateTransactionStatus)
	}
}
