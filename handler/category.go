package handler

import (
	"net/http"
	"strconv"

	model "simple-crud/models"
	"simple-crud/service"
	"simple-crud/util"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

// ============================
// GET ALL
// ============================
//
// GetAll godoc
// @Summary Get all categories
// @Description Retrieve all categories
// @Tags categories
// @Produce json
// @Success 200 {object} util.JSONResponse{data=[]model.Category}
// @Failure 500 {object} util.JSONResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "categories retrieved",
		Data:    categories,
	})
}

// ============================
// GET BY ID
// ============================
//
// GetByID godoc
// @Summary Get category by ID
// @Description Get category detail by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} util.JSONResponse{data=model.Category}
// @Failure 400 {object} util.JSONResponse
// @Failure 404 {object} util.JSONResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "category retrieved",
		Data:    category,
	})
}

// ============================
// CREATE
// ============================
//
// Create godoc
// @Summary Create new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body model.Category true "Category payload"
// @Success 201 {object} util.JSONResponse{data=model.Category}
// @Failure 400 {object} util.JSONResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var payload model.Category
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	created, err := h.service.Create(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, util.JSONResponse{
		Message: "category created",
		Data:    created,
	})
}

// ============================
// UPDATE
// ============================
//
// Update godoc
// @Summary Update category
// @Description Update category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.Category true "Category payload"
// @Success 200 {object} util.JSONResponse
// @Failure 400 {object} util.JSONResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	var payload model.Category
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := h.service.Update(id, payload); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "category updated",
		Data: gin.H{
			"id":       id,
			"category": payload,
		},
	})
}

// ============================
// DELETE
// ============================
//
// Delete godoc
// @Summary Delete category
// @Description Delete category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} util.JSONResponse
// @Failure 400 {object} util.JSONResponse
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, util.JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, util.JSONResponse{
		Message: "category deleted",
		Data: gin.H{
			"id": id,
		},
	})
}
