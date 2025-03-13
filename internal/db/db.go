package db

import (
	"fmt"
	"log"
	"time"

	"loan-billing-system/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establishes a connection to the database
func Connect(config Config) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate performs database migrations
func Migrate(db *gorm.DB) error {
	// Migrate tables in the correct order to avoid foreign key constraint issues
	if err := db.AutoMigrate(&models.Borrower{}); err != nil {
		return fmt.Errorf("failed to migrate borrowers table: %w", err)
	}

	if err := db.AutoMigrate(&models.Loan{}); err != nil {
		return fmt.Errorf("failed to migrate loans table: %w", err)
	}

	if err := db.AutoMigrate(&models.Schedule{}); err != nil {
		return fmt.Errorf("failed to migrate schedules table: %w", err)
	}

	if err := db.AutoMigrate(&models.Payment{}); err != nil {
		return fmt.Errorf("failed to migrate payments table: %w", err)
	}

	return nil
}
