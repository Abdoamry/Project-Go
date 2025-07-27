package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	var db *gorm.DB
	var err error
	maxRetries := 5

	// Retry connection with exponential backoff
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// Test the connection
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
			}

			if err == nil {
				// Set connection pool settings
				sqlDB.SetMaxIdleConns(10)
				sqlDB.SetMaxOpenConns(100)
				sqlDB.SetConnMaxLifetime(time.Hour)

				DB = db
				log.Println("Successfully connected to the database")
				return nil
			}
		}

		if i < maxRetries-1 {
			waitTime := time.Duration(i*i) * time.Second
			log.Printf("Failed to connect to database (attempt %d/%d), retrying in %v: %v", 
				i+1, maxRetries, waitTime, err)
			time.Sleep(waitTime)
		}
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}