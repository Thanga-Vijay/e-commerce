package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ecommerce/product-service/internal/cache"
	"github.com/ecommerce/product-service/internal/config"
	"github.com/ecommerce/product-service/internal/handlers"
	"github.com/ecommerce/product-service/internal/middleware"
	"github.com/ecommerce/product-service/internal/models"
	"github.com/ecommerce/product-service/internal/repository"
	"github.com/ecommerce/product-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	redisClient := initRedis(cfg)

	// Auto migrate database
	if err := db.AutoMigrate(&models.Product{}, &models.Category{}, &models.Review{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize cache
	productCache := cache.NewProductCache(redisClient)

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo, productCache, time.Duration(cfg.Cache.ProductTTL)*time.Second)
	categoryService := service.NewCategoryService(categoryRepo, productCache, time.Duration(cfg.Cache.CategoryTTL)*time.Second)
	reviewService := service.NewReviewService(reviewRepo, productRepo, productCache)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	healthHandler := handlers.NewHealthHandler()

	// Setup Gin router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORSMiddleware())

	// Health check routes
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			// Public routes
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)

			// Protected routes (admin only)
			protected := products.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			protected.Use(middleware.RoleMiddleware("admin"))
			{
				protected.POST("", productHandler.CreateProduct)
				protected.PUT("/:id", productHandler.UpdateProduct)
				protected.DELETE("/:id", productHandler.DeleteProduct)
			}
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			// Public routes
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)

			// Protected routes (admin only)
			protected := categories.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			protected.Use(middleware.RoleMiddleware("admin"))
			{
				protected.POST("", categoryHandler.CreateCategory)
				protected.PUT("/:id", categoryHandler.UpdateCategory)
				protected.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}

		// Review routes
		reviews := v1.Group("/reviews")
		{
			// Public routes
			reviews.GET("/product/:productId", reviewHandler.GetProductReviews)

			// Protected routes (authenticated users)
			protected := reviews.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			{
				protected.POST("", reviewHandler.CreateReview)
				protected.PUT("/:id", reviewHandler.UpdateReview)
				protected.DELETE("/:id", reviewHandler.DeleteReview)
			}
		}
	}

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Product Service starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func initRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	return client
}
