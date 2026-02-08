package main

import (
	"log"
	"net/http"

	"simple-crud/config"
	"simple-crud/database"
	"simple-crud/handler"
	"simple-crud/repository"
	"simple-crud/service"

	_ "simple-crud/docs"

	"github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/gin-gonic/gin"
)

// @title Simple CRUD API
// @version 1.0
// @description REST API for product and category
// @BasePath /
func main() {

	cfg := config.Load()

	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// === Dependency Injection ===
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(*categoryRepo)
	categoryHandler := handler.NewCategoryHandler(*categoryService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(*productRepo)
	productHandler := handler.NewProductHandler(*productService)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(*transactionRepo)
	transactionHandler := handler.NewTransactionHandler(*transactionService)

	// === Gin Router ===
	router := gin.Default()

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Simple CRUD API"})
	})

	// Healthcheck
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// EXPOSE swagger.json via HTTP (ini yang dibaca Scalar)
	router.GET("/openapi.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	// Scalar UI
	router.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL: "/openapi.json",
		Title:   "Simple CRUD API",
		Theme:   "light",
	}))

	// === Routes ===
	api := router.Group("/api/v1")
	{
		cat := api.Group("/categories")
		{
			cat.GET("", categoryHandler.GetAll)
			cat.GET("/:id", categoryHandler.GetByID)
			cat.POST("", categoryHandler.Create)
			cat.PUT("/:id", categoryHandler.Update)
			cat.DELETE("/:id", categoryHandler.Delete)
		}

		product := api.Group("/products")
		{
			product.GET("", productHandler.GetAll)
			product.GET("/:id", productHandler.GetById)
			product.POST("", productHandler.Create)
			product.PUT("/:id", productHandler.Update)
			product.DELETE("/:id", productHandler.Delete)
		}

		api.POST("/checkout", transactionHandler.Checkout)

		report := api.Group("/report")
		{
			report.GET("/hari-ini", transactionHandler.GetSalesSummary)
			report.GET("", transactionHandler.GetSalesSummary)
		}

	}

	log.Println("Server running on port", cfg.Port)
	router.Run(":" + cfg.Port)
}
