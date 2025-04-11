package handlers

import (
	"net/http"

	"avions-club/backend/database"
	"avions-club/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetProjects returns all projects
func GetProjects(c *gin.Context) {
	var projects []models.Project
	result := database.DB.Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject returns a specific project
func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := database.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = uuid.New()
	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// UpdateProject updates an existing project
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := database.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject deletes a project
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := database.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := database.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
