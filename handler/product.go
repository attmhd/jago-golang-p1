package handler

import (
	"net/http"
	"strconv"

	model "simple-crud/models"
	"simple-crud/service"

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
// @Success 200 {object} JSONResponse{data=[]ProductResp}
// @Failure 500 {object} JSONResponse
// @Router /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse{
			Message: "Internal Server Error",
			Data:    nil,
		})
		return
	}

	resp := make([]ProductResp, 0, len(products))
	for _, p := range products {
		resp = append(resp, ProductResp{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Stock: p.Stock,
			Category: Category{
				ID:   p.CategoryID,
				Name: p.CategoryName,
			},
		})
	}

	c.JSON(http.StatusOK, JSONResponse{
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
// @Success 200 {object} JSONResponse{data=ProductResp}
// @Failure 400 {object} JSONResponse
// @Failure 404 {object} JSONResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp := ProductResp{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: Category{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}

	c.JSON(http.StatusOK, JSONResponse{
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
// @Success 201 {object} JSONResponse{data=ProductResp}
// @Failure 400 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var payload model.Product
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "Invalid payload",
			Data:    nil,
		})
		return
	}

	product, err := h.service.Create(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse{
			Message: "Failed to create product",
			Data:    nil,
		})
		return
	}

	resp := ProductResp{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: Category{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}

	c.JSON(http.StatusCreated, JSONResponse{
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
// @Success 200 {object} JSONResponse{data=ProductResp}
// @Failure 400 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	var payload model.Product
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "Invalid payload",
			Data:    nil,
		})
		return
	}

	payload.ID = id

	product, err := h.service.Update(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse{
			Message: "Failed to update product",
			Data:    nil,
		})
		return
	}

	resp := ProductResp{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: Category{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
	}

	c.JSON(http.StatusOK, JSONResponse{
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
// @Success 200 {object} JSONResponse
// @Failure 400 {object} JSONResponse
// @Failure 404 {object} JSONResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, JSONResponse{
		Message: "Product deleted successfully",
		Data:    nil,
	})
}
