package handler

import (
	"net/http"

	"simple-crud/models"
	"simple-crud/service"
	"simple-crud/util"

	"github.com/gin-gonic/gin"
)

type ProdukTerlarisResp struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type SalesSummaryResp struct {
	TotalRevenue   int                `json:"total_revenue"`
	TotalTransaksi int                `json:"total_transaksi"`
	ProdukTerlaris ProdukTerlarisResp `json:"produk_terlaris"`
}

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
// @Success 200 {object} util.JSONResponse{data=model.Transaction}
// @Failure 400 {object} util.JSONResponse
// @Failure 500 {object} util.JSONResponse
// @Router /api/v1/checkout [post]
func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	transaction, err := h.service.Checkout(req.Items, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "Checkout berhasil",
		Data:    transaction,
	})
}

// GetSalesSummary godoc
// @Summary Get sales summary
// @Description Get sales summary for today or within a date range if start_date and end_date are provided
// @Tags transactions
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} util.JSONResponse{data=handler.SalesSummaryResp}
// @Failure 400 {object} util.JSONResponse
// @Failure 500 {object} util.JSONResponse
// @Router /api/v1/report/hari-ini [get]
// @Router /api/v1/report [get]
func (h *TransactionHandler) GetSalesSummary(c *gin.Context) {
	// Parse optional query params
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var (
		summary *util.SalesSummary
		err     error
	)

	// Jika ada rentang tanggal, gunakan summary by range
	if startDate != "" && endDate != "" {
		summary, err = h.service.GetSalesSummaryByRange(startDate, endDate)
	} else {
		// Jika tidak ada params, tampilkan ringkasan hari ini
		summary, err = h.service.GetSalesSummary()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Produk terlaris mengikuti rentang/harian yang sama
	var topSellingProduct *models.TopSellingProduct
	if startDate != "" && endDate != "" {
		topSellingProduct, err = h.service.GetTopSellingProductByRange(startDate, endDate)
	} else {
		topSellingProduct, err = h.service.GetTopSellingProduct()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp := SalesSummaryResp{
		TotalRevenue:   summary.TotalRevenue,
		TotalTransaksi: summary.TotalTransaksi,
		ProdukTerlaris: ProdukTerlarisResp{
			Nama:       topSellingProduct.Name,
			QtyTerjual: topSellingProduct.QtySold,
		},
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "Sales summary",
		Data:    resp,
	})
}
