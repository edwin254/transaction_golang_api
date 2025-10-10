package controller

import (
	"gapstack-api/internal/model"
	"gapstack-api/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service *service.TransactionService
}

func NewTransactionController(svc *service.TransactionService) *TransactionController {
	return &TransactionController{service: svc}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var trx model.Transaction
	log.Println("Received transaction creation request")

	if err := ctx.ShouldBindJSON(&trx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := c.service.Create(&trx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	trx, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, trx)
}

func (c *TransactionController) ListTransactions(ctx *gin.Context) {
	// Default values
	limit := 10
	offset := 0

	// Parse optional query parameters
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Fetch paginated list
	transactions, err := c.service.List(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"offset":       offset,
		"count":        len(transactions),
		"transactions": transactions,
	})
}

func (c *TransactionController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	statusInput := model.Status(req.Status)
	if statusInput != model.StatusCompleted && statusInput != model.StatusFailed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status must be either 'completed' or 'failed'"})
		return
	}
	updated, err := c.service.UpdateStatus(id, statusInput)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}
