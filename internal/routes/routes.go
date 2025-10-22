package routes

import (
	"gapstack-api/internal/controller"
	"gapstack-api/internal/repository"
	"gapstack-api/internal/service"
	"net/http"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	repo := repository.NewTransactionRepository(db)
	svc := service.NewTransactionService(repo)
	ctrl := controller.NewTransactionController(svc)

	// Health Check Endpoint
	router.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"message": "database unreachable",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "service healthy",
		})
	})

	api := router.Group("/api/v1")
	trx := api.Group("/transactions")
	{
		trx.POST("/", ctrl.CreateTransaction)
		trx.GET("/:id", ctrl.GetTransactionByID)
		trx.GET("/", ctrl.ListTransactions)
		trx.PUT("/:id/status", ctrl.UpdateTransactionStatus)
	}
}
