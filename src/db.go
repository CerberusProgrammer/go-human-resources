package src

import (
	"fmt"
	"log"

	"go-human-resources/src/employee"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database holds the connection to the database
type Database struct {
	DB *gorm.DB
}

// InitDB initializes the database connection
func InitDB(config *Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("Connected to database successfully")

	// Run migrations
	err = db.AutoMigrate(&employee.Employee{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Database migrations completed")

	return &Database{DB: db}, nil
}
