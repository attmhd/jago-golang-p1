package handler

import (
	"fmt"
	"net/http"
	"strconv"

	model "simple-crud/models"
	"simple-crud/service"
	"simple-crud/util"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: svc,
	}
}

// ============================
// GET ALL PRODUCTS
// ============================
//
// GetAll godoc
// @Summary Get all products
// @Description Get list of products with category
// @Tags products
// @Produce json
// @Success 200 {object} util.JSONResponse{data=[]util.ProductResp}
// @Failure 500 {object} util.JSONResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	name := c.Request.URL.Query().Get("name")
	fmt.Println(name)
	products, err := h.service.GetAll(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: "Internal Server Error",
			Data:    nil,
		})
		return
	}

	resp := make([]util.ProductResp, 0, len(products))
	for _, p := range products {
		resp = append(resp, util.ProductResp{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Stock: p.Stock,
			Category: util.Category{
				ID:   p.CategoryID,
				Name: p.CategoryName,
			},
		})
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "Success",
		Data:    resp,
	})
}

// ============================
// GET PRODUCT BY ID
// ============================
//
// GetById godoc
// @Summary Get product by ID
// @Description Get product detail with category
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} util.JSONResponse{data=util.ProductResp}
// @Failure 400 {object} util.JSONResponse
// @Failure 404 {object} util.JSONResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp := util.ProductResp{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: util.Category{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "Success",
		Data:    resp,
	})
}

// ============================
// CREATE PRODUCT
// ============================
//
// Create godoc
// @Summary Create product
// @Description Create new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.Product true "Product payload"
// @Success 201 {object} util.JSONResponse{data=util.ProductResp}
// @Failure 400 {object} util.JSONResponse
// @Failure 500 {object} util.JSONResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var payload model.Product
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "Invalid payload",
			Data:    nil,
		})
		return
	}

	product, err := h.service.Create(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: "Failed to create product",
			Data:    nil,
		})
		return
	}

	resp := util.ProductResp{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: util.Category{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}

	c.JSON(http.StatusCreated, util.JSONResponse{
		Message: "Product created successfully",
		Data:    resp,
	})
}

// ============================
// UPDATE PRODUCT
// ============================
//
// Update godoc
// @Summary Update product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.Product true "Product payload"
// @Success 200 {object} util.JSONResponse{data=util.ProductResp}
// @Failure 400 {object} util.JSONResponse
// @Failure 500 {object} util.JSONResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	var payload model.Product
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "Invalid payload",
			Data:    nil,
		})
		return
	}

	payload.ID = id

	product, err := h.service.Update(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: "Failed to update product",
			Data:    nil,
		})
		return
	}

	resp := util.ProductResp{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: util.Category{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "Product updated successfully",
		Data:    resp,
	})
}

// ============================
// DELETE PRODUCT
// ============================
//
// Delete godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} util.JSONResponse
// @Failure 400 {object} util.JSONResponse
// @Failure 404 {object} util.JSONResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "Product deleted successfully",
		Data:    nil,
	})
}
