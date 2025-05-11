package main

import (
	"backend/cache"
	"backend/config"
	"backend/database"
	"backend/handler"
	"backend/repository"
	"backend/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to write database
	writeDB, err := database.ConnectWriteDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to write database: %v", err)
	}
	defer writeDB.Close()

	// Connect to read database
	readDB, err := database.ConnectReadDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to read database: %v", err)
	}
	defer readDB.Close()

	// Connect to Redis if configured
	redisClient, err := cache.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	if redisClient != nil {
		defer redisClient.Close()
		log.Println("Redis connection established")
	} else {
		log.Println("Redis not configured, using database for likes")
	}

	// Initialize dependencies
	productRepo := repository.NewProductRepository(writeDB, readDB, redisClient)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Initialize router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = false // credentials를 허용하지 않음
	corsConfig.MaxAge = 12 * 60 * 60    // 12 hours

	// Apply CORS middleware
	router.Use(cors.New(corsConfig))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProductByID)
		api.POST("/products/:id/like", productHandler.LikeProduct)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
