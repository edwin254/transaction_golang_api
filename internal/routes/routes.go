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
	trx := api.Group("/transactions")
	{
		trx.POST("/", ctrl.CreateTransaction)
		trx.GET("/:id", ctrl.GetTransactionByID)
		trx.GET("/", ctrl.ListTransactions)
		trx.PUT("/:id/status", ctrl.UpdateTransactionStatus)
	}
}
