package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"github.com/abdoamry/Project-go/database"
	"github.com/abdoamry/Project-go/routes"
	logger "github.com/abdoamry/Project-go/utils"
)



func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger.InitLogger()

	// Connect to database
	if err := database.ConnectDB(); err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize Redis
	if err := database.InitRedis(); err != nil {
		zap.L().Fatal("Failed to initialize Redis", zap.Error(err))
	}

	// Create new Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Project Go API",
	})

	// Setup routes
	routes.SetupRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			zap.L().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	zap.L().Info("Server started", zap.String("port", port))

	// Wait for interrupt signal
	<-done
	zap.L().Info("Shutting down server...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server stopped")
}
