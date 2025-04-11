package handlers

import (
	"net/http"

	"avions-club/backend/database"
	"avions-club/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetMembers returns all members
func GetMembers(c *gin.Context) {
	var members []models.Member
	result := database.DB.Find(&members)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetMember returns a specific member
func GetMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member

	if err := database.DB.First(&member, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	c.JSON(http.StatusOK, member)
}

// CreateMember creates a new member
func CreateMember(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member.ID = uuid.New()
	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating member"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// UpdateMember updates an existing member
func UpdateMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member

	if err := database.DB.First(&member, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating member"})
		return
	}

	c.JSON(http.StatusOK, member)
}

// DeleteMember deletes a member
func DeleteMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member

	if err := database.DB.First(&member, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	if err := database.DB.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}
