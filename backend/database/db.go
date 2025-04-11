package database

import (
	"fmt"
	"log"
	"os"

	"avions-club/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	// Build connection string for the connection pooler
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	log.Printf("Attempting to connect to database...")

	// Configure GORM with connection pool settings
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements
	}), &gorm.Config{
		PrepareStmt: false, // Disable prepare statements as they're not supported by the transaction pooler
	})

	if err != nil {
		log.Fatal("Failed to create GORM instance:", err)
	}

	// Configure connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)

	// Test the connection
	var version string
	if err := gormDB.Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	log.Printf("Connected to database: %s", version)

	// Auto Migrate the schemas
	err = gormDB.AutoMigrate(
		&models.Member{},
		&models.Project{},
		&models.Blog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = gormDB
	log.Println("Database connection and migrations completed")
}
