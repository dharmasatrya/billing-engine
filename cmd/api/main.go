package main

import (
	"fmt"
	"log"
	"net/http"

	"loan-billing-system/config"
	"loan-billing-system/internal/api"
	"loan-billing-system/internal/db"
	"loan-billing-system/internal/repositories"
	"loan-billing-system/internal/scheduler"
	"loan-billing-system/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "loan-billing-system/docs"
)

// @title Loan Billing System API
// @version 1.0
// @description API for managing loans, borrowers, payments, and delinquency status
// @host localhost:8080
// @BasePath /api

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	database, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize repositories
	repoManager := repositories.NewGormRepositoryManager(database)

	// Initialize services
	loanService := services.NewLoanService(repoManager)
	borrowerService := services.NewBorrowerService(repoManager)

	// Set up scheduler
	scheduler := scheduler.NewScheduler(database, loanService)
	scheduler.Start()
	defer scheduler.Stop()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})

	// Set up API routes
	api.SetupRoutes(e, database, borrowerService, loanService)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
