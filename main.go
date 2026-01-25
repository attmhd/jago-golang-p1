package main

import (
	"net/http"
	"simple-crud/handler"
	"simple-crud/repository"
	"simple-crud/service"

	"github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/gin-gonic/gin"
)

func main() {

	repo := repository.NewCategoryRepository()
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	cat := router.Group("/categories")
	{
		cat.GET("", h.GetAll)
		cat.GET("/:id", h.GetByID)
		cat.POST("", h.Create)
		cat.PUT("/:id", h.Update)
		cat.DELETE("/:id", h.Delete)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Simple CRUD API"})
	})

	router.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: "./docs/swagger.json",
		Title:        "Example API",
		Theme:        "light", // or "dark"
	}))

	router.Run(":8080")
}
