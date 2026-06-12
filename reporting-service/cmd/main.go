package main

import (
	"fmt"
	"log"
	"reporting-service/config"
	"reporting-service/handlers"
	"reporting-service/internal/cache"
	"reporting-service/internal/client"
	"reporting-service/internal/export"
	"reporting-service/internal/models"
	"reporting-service/internal/repository"
	"reporting-service/internal/service"
	"reporting-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	if err := db.AutoMigrate(&models.Report{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Initialize services
	cacheService := cache.NewCacheService(redisClient)
	serviceClient := client.NewServiceClient(
		cfg.Services.AuthServiceURL,
		cfg.Services.ProductServiceURL,
		cfg.Services.OrderServiceURL,
		cfg.Services.PaymentServiceURL,
	)
	reportRepo := repository.NewReportRepository(db)
	metricsService := service.NewMetricsService(serviceClient, cacheService)
	reportService := service.NewReportService(reportRepo)
	exportService := export.NewExportService()

	// Initialize handlers
	metricsHandler := handlers.NewMetricsHandler(metricsService, reportService, exportService)

	// Setup router
	router := setupRouter(cfg, metricsHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Reporting Service starting on port %s", cfg.Server.Port)
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

func setupRouter(cfg *config.Config, metricsHandler *handlers.MetricsHandler) *gin.Engine {
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
		// Dashboard metrics
		api.GET("/dashboard", metricsHandler.GetDashboard)
		api.GET("/dashboard/export", metricsHandler.ExportDashboard)

		// Revenue reports
		api.GET("/reports/revenue", metricsHandler.GetRevenueReport)
		api.GET("/reports/revenue/export", metricsHandler.ExportRevenueReport)

		// Product reports
		api.GET("/reports/products", metricsHandler.GetTopProducts)
		api.GET("/reports/products/export", metricsHandler.ExportTopProducts)

		// Customer reports
		api.GET("/reports/customers", metricsHandler.GetCustomerReport)
		api.GET("/reports/customers/export", metricsHandler.ExportCustomerReport)

		// Saved reports
		api.POST("/reports", metricsHandler.SaveReport)
		api.GET("/reports", metricsHandler.GetSavedReports)
		api.GET("/reports/:id", metricsHandler.GetReportByID)
		api.DELETE("/reports/:id", metricsHandler.DeleteReport)

		// Admin routes
		admin := api.Group("")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/cache/invalidate", metricsHandler.InvalidateCache)
		}
	}

	return router
}
