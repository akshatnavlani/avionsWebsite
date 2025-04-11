package handlers

import (
	"net/http"
	"os"

	"avions-club/backend/middleware"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// Login handles admin authentication
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if password matches
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if req.Password != adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}
