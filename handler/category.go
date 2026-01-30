package handler

import (
	"net/http"
	"strconv"

	model "simple-crud/models"
	"simple-crud/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, JSONResponse{
		Message: "categories retrieved",
		Data:    categories,
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, JSONResponse{
		Message: "category retrieved",
		Data:    category,
	})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var payload model.Category
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	created, err := h.service.Create(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, JSONResponse{
		Message: "category created",
		Data:    created,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	var payload model.Category
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := h.service.Update(id, payload); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, JSONResponse{
		Message: "category updated",
		Data: gin.H{
			"id":       id,
			"category": payload,
		},
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, JSONResponse{
		Message: "category deleted",
		Data: gin.H{
			"id": id,
		},
	})
}
