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

	"github.com/ecommerce/wishlist-service/internal/client"
	"github.com/ecommerce/wishlist-service/internal/config"
	"github.com/ecommerce/wishlist-service/internal/handlers"
	"github.com/ecommerce/wishlist-service/internal/middleware"
	"github.com/ecommerce/wishlist-service/internal/models"
	"github.com/ecommerce/wishlist-service/internal/repository"
	"github.com/ecommerce/wishlist-service/internal/service"
	"github.com/gin-gonic/gin"
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

	// Auto migrate database
	if err := db.AutoMigrate(&models.Wishlist{}, &models.WishlistItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize clients
	productClient := client.NewProductClient(cfg.ProductService.URL)
	cartClient := client.NewCartClient(cfg.CartService.URL)

	// Initialize repositories
	wishlistRepo := repository.NewWishlistRepository(db)

	// Initialize services
	wishlistService := service.NewWishlistService(wishlistRepo, productClient, cartClient)

	// Initialize handlers
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
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
		// Wishlist routes (all protected)
		wishlist := v1.Group("/wishlist")
		wishlist.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			wishlist.GET("", wishlistHandler.GetWishlist)
			wishlist.POST("/items", wishlistHandler.AddToWishlist)
			wishlist.DELETE("/items/:id", wishlistHandler.RemoveFromWishlist)
			wishlist.DELETE("", wishlistHandler.ClearWishlist)
			wishlist.POST("/items/:id/move-to-cart", wishlistHandler.MoveToCart)
		}
	}

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Wishlist Service starting on port %s", cfg.Server.Port)
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
