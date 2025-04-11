package main

import (
	"log"
	"os"

	"avions-club/backend/database"
	"avions-club/backend/routes"
	"avions-club/backend/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Initialize Supabase storage
	if err := storage.InitStorage(); err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	// Initialize database
	database.InitDB()

	// Debug: Print environment variables
	log.Println("SUPABASE_URL:", os.Getenv("SUPABASE_URL"))
	log.Println("SUPABASE_SERVICE_KEY exists:", os.Getenv("SUPABASE_SERVICE_KEY") != "")

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = []string{allowedOrigins}
	} else {
		config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Setup routes
	routes.SetupRoutes(r)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
