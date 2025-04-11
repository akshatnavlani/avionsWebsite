package handlers

import (
	"net/http"

	"avions-club/backend/database"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Database connection error",
		})
		return
	}

	// Ping database
	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "connected",
	})
}
