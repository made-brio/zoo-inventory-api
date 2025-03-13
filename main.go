package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"zoo-inventory/database"
	"zoo-inventory/internal/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	DB  *sql.DB
	err error
)

func init() {
	// Load environment variables
	if err = godotenv.Load("config/.env"); err != nil {
		fmt.Println("Warning: .env file not found, using default environment variables.")
	}
}

func connectDatabase() (*sql.DB, error) {
	// Construct DSN for MySQL connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	// Establish database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Verify database connectivity
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Create database if it does not exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS zoo_inventory")
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Update DSN to use the newly created database
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection with database: %w", err)
	}

	return db, nil
}

func main() {
	// Connect to the database
	DB, err = connectDatabase()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer DB.Close()

	// Execute database migrations
	if err := database.DBMigrate(DB); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()
	routes.RegisterRoutes(router, DB)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085" // Default port
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
