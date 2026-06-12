package main

import (
	"fmt"
	"log"
	"payment-service/config"
	"payment-service/handlers"
	"payment-service/internal/models"
	"payment-service/internal/repository"
	"payment-service/internal/service"
	"payment-service/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate Stripe configuration
	if cfg.Stripe.SecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is required")
	}

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Refund{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
	}

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	paymentService := service.NewPaymentService(paymentRepo, cfg)

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup router
	router := setupRouter(cfg, paymentHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Payment Service starting on port %s", cfg.Server.Port)
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

func setupRouter(cfg *config.Config, paymentHandler *handlers.PaymentHandler) *gin.Engine {
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

	// Webhook endpoint (no auth)
	router.POST("/api/v1/webhooks/stripe", paymentHandler.HandleWebhook)

	// API routes
	api := router.Group("/api/v1")
	{
		// Payment routes (authenticated)
		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			payments.POST("/intent", paymentHandler.CreatePaymentIntent)
			payments.GET("/:id", paymentHandler.GetPaymentByID)
			payments.GET("/order/:orderId", paymentHandler.GetPaymentByOrderID)

			// Admin only routes
			admin := payments.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.POST("/:id/refund", paymentHandler.ProcessRefund)
			}
		}
	}

	return router
}
