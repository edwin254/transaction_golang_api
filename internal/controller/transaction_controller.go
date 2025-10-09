package controller

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "gapstack-api/internal/model"
    "gapstack-api/internal/service"
)

type TransactionController struct {
    service *service.TransactionService
}

func NewTransactionController(svc *service.TransactionService) *TransactionController {
    return &TransactionController{service: svc}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
    var tx model.Transaction
    if err := ctx.ShouldBindJSON(&tx); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    created, err := c.service.CreateTransaction(&tx)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, created)
}

func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
    id := ctx.Param("id")
    tx, err := c.service.GetTransactionByID(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, tx)
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
    tx, err := c.service.UpdateTransactionStatus(id, body.Status)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, tx)
}
