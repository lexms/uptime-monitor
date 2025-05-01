package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tou01/uptime-monitor/pkg/handler"
	"github.com/tou01/uptime-monitor/pkg/middleware"
	"github.com/tou01/uptime-monitor/pkg/service"
)

const (
	jwtSecretKey = "your-secret-key-replace-in-production" // Change this in production
	tokenTTL     = 24 * time.Hour                          // Token valid for 24 hours
	serverPort   = ":8080"
)

func main() {
	// Set up logger
	logger := log.New(os.Stdout, "UPTIME-MONITOR: ", log.LstdFlags)

	// Initialize router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create services
	monitorService := service.NewMonitorService()

	// Create handlers
	authHandler := handler.NewAuthHandler(jwtSecretKey, tokenTTL)
	endpointHandler := handler.NewEndpointHandler(monitorService)

	// Set up API routes
	api := router.Group("/api")

	// Public routes (no auth required)
	authHandler.RegisterRoutes(api)

	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Protected routes (auth required)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(jwtSecretKey))
	{
		endpointHandler.RegisterRoutes(protected)
	}

	// Welcome route
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Uptime Monitor API",
		})
	})

	// Configure server
	server := &http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		logger.Printf("Server starting on port %s", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()

	// Start a goroutine to handle monitoring results
	go handleMonitoringResults(monitorService, logger)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}

// handleMonitoringResults processes the results from endpoint checks
func handleMonitoringResults(monitorService *service.MonitorService, logger *log.Logger) {
	resultChan := monitorService.ListenForResults()

	for result := range resultChan {
		if !result.Success {
			logger.Printf("Endpoint %s is DOWN: %s - %s",
				result.EndpointID,
				result.StatusText,
				result.Error,
			)

			// Here you would typically:
			// 1. Store the result in a database
			// 2. Check if an alert should be triggered
			// 3. Send notifications if needed
		} else {
			logger.Printf("Endpoint %s is UP: %s - %dms",
				result.EndpointID,
				result.StatusText,
				result.ResponseTime,
			)

			// Store successful result in database
		}
	}
}
