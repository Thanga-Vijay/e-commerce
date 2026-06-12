package main

import (
	"fmt"
	"log"
	"order-service/config"
	"order-service/handlers"
	"order-service/internal/client"
	"order-service/internal/models"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
		&models.OrderStatus{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
	}

	// Initialize clients
	cartClient := client.NewCartClient(cfg.Services.CartServiceURL)

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	orderService := service.NewOrderService(orderRepo, cartClient)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router
	router := setupRouter(cfg, orderHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Order Service starting on port %s", cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

	return db, nil
}

func setupRouter(cfg *config.Config, orderHandler *handlers.OrderHandler) *gin.Engine {
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Order routes (authenticated)
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			orders.GET("/:id/track", orderHandler.TrackOrder)

			// Admin only routes
			admin := orders.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			}
		}
	}

	return router
}
