package controller

import (
	"gapstack-api/internal/model"
	"gapstack-api/internal/service"
	"net/http"

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
	if err := ctx.ShouldBindJSON(&trx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := c.service.CreateTransaction(&trx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	trx, err := c.service.GetTransactionByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, trx)
}

func (c *TransactionController) ListTransactions(ctx *gin.Context) {
	transactions, err := c.service.ListTransactions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transactions)
}

func (c *TransactionController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var body struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trx, err := c.service.UpdateTransactionStatus(id, body.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, trx)
}
