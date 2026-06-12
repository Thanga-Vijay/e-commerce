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

	"github.com/ecommerce/cart-service/internal/cache"
	"github.com/ecommerce/cart-service/internal/client"
	"github.com/ecommerce/cart-service/internal/config"
	"github.com/ecommerce/cart-service/internal/handlers"
	"github.com/ecommerce/cart-service/internal/middleware"
	"github.com/ecommerce/cart-service/internal/models"
	"github.com/ecommerce/cart-service/internal/repository"
	"github.com/ecommerce/cart-service/internal/service"
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
	if err := db.AutoMigrate(&models.Cart{}, &models.CartItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize cache
	cartCache := cache.NewCartCache(redisClient)

	// Initialize clients
	productClient := client.NewProductClient(cfg.ProductService.URL)

	// Initialize repositories
	cartRepo := repository.NewCartRepository(db)

	// Initialize services
	cartService := service.NewCartService(
		cartRepo,
		cartCache,
		productClient,
		time.Duration(cfg.Cache.CartTTL)*time.Second,
	)

	// Initialize handlers
	cartHandler := handlers.NewCartHandler(cartService)
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
		// Cart routes (all protected)
		cart := v1.Group("/cart")
		cart.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("/items", cartHandler.AddToCart)
			cart.PUT("/items/:id", cartHandler.UpdateCartItem)
			cart.DELETE("/items/:id", cartHandler.RemoveFromCart)
			cart.DELETE("", cartHandler.ClearCart)
		}
	}

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Cart Service starting on port %s", cfg.Server.Port)
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
