package handler

import (
	"net/http"

	model "simple-crud/models"
	"simple-crud/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: svc,
	}
}

// ============================
// CHECKOUT
// ============================
//
// Checkout godoc
// @Summary Checkout transaction
// @Description Create transaction from cart items and update product stock
// @Tags transactions
// @Accept json
// @Produce json
// @Param checkout body model.CheckoutRequest true "Checkout payload"
// @Success 200 {object} model.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/checkout [post]
func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req model.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	transaction, err := h.service.Checkout(req.Items, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
