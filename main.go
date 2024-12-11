package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"zoo-inventory/database"
	"zoo-inventory/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func init() {
	// Load environment variables
	err = godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Warning: .env file not found, using default environment variables.")
	}
}

func connectDatabase() (*sql.DB, error) {
	// Database connection string
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
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

	// Run migrations
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
		log.Fatalf("Server failed to start: %v", err)
	}
}
