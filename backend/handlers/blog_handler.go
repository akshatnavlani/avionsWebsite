package handlers

import (
	"fmt"
	"net/http"

	"avions-club/backend/database"
	"avions-club/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetBlogs returns all blogs with their authors
func GetBlogs(c *gin.Context) {
	var blogs []models.Blog
	result := database.DB.Preload("Author").Find(&blogs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching blogs"})
		return
	}

	c.JSON(http.StatusOK, blogs)
}

// GetBlog returns a specific blog with its author
func GetBlog(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	blogID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid blog ID format: %s", id),
		})
		return
	}

	var blog models.Blog
	// Use First to get a single record, respecting soft deletes
	if err := database.DB.Preload("Author").First(&blog, "id = ?", blogID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Blog not found with ID: %s", id),
		})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// CreateBlog creates a new blog
func CreateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog.ID = uuid.New()
	if err := database.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating blog"})
		return
	}

	// Fetch the complete blog with author details
	if err := database.DB.Preload("Author").First(&blog, "id = ?", blog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching created blog"})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

// UpdateBlog updates an existing blog
func UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog

	if err := database.DB.First(&blog, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating blog"})
		return
	}

	// Fetch the updated blog with author details
	if err := database.DB.Preload("Author").First(&blog, "id = ?", blog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated blog"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// DeleteBlog deletes a blog
func DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	blogID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid blog ID format: %s", id),
		})
		return
	}

	// First check if the blog exists
	var blog models.Blog
	if err := database.DB.First(&blog, "id = ?", blogID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Blog not found with ID: %s", id),
		})
		return
	}

	// Hard delete the blog
	if err := database.DB.Unscoped().Delete(&blog, "id = ?", blogID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error deleting blog: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Blog deleted successfully",
	})
}
