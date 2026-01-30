package main

import (
	"log"
	"net/http"
	"simple-crud/config"
	"simple-crud/database"
	"simple-crud/handler"
	"simple-crud/repository"
	"simple-crud/service"

	"github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/gin-gonic/gin"
)

func main() {

	config := config.Load()

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database : ", err)
	}

	defer db.Close()

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(*categoryRepo)
	categoryHandler := handler.NewCategoryHandler(*categoryService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Health check successful"})
	})

	cat := router.Group("/categories")
	{
		cat.GET("", categoryHandler.GetAll)
		cat.GET("/:id", categoryHandler.GetByID)
		cat.POST("", categoryHandler.Create)
		cat.PUT("/:id", categoryHandler.Update)
		cat.DELETE("/:id", categoryHandler.Delete)
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

	router.Run(":" + config.Port)
}
