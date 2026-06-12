package main

import (
	"fmt"
	"log"
	"notification-service/config"
	"notification-service/handlers"
	"notification-service/internal/mailer"
	"notification-service/internal/models"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate SMTP configuration
	if cfg.SMTP.Username == "" || cfg.SMTP.Password == "" {
		log.Println("Warning: SMTP credentials not configured. Email sending will fail.")
	}

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Notification{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
	}

	// Initialize mailer
	m := mailer.NewMailer(cfg)

	// Initialize repositories
	notificationRepo := repository.NewNotificationRepository(db)

	// Initialize services
	notificationService := service.NewNotificationService(notificationRepo, m)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Setup router
	router := setupRouter(cfg, notificationHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Notification Service starting on port %s", cfg.Server.Port)
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

func setupRouter(cfg *config.Config, notificationHandler *handlers.NotificationHandler) *gin.Engine {
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
	api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// Notification routes
		notifications := api.Group("/notifications")
		{
			// Authenticated routes
			notifications.POST("", notificationHandler.SendNotification)
			notifications.GET("/my", notificationHandler.GetNotificationsByUserID)
			notifications.GET("/:id", notificationHandler.GetNotificationByID)

			// Admin only routes (for processing/retry)
			admin := notifications.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.POST("/process-pending", notificationHandler.ProcessPending)
				admin.POST("/retry-failed", notificationHandler.RetryFailed)
			}
		}
	}

	return router
}
