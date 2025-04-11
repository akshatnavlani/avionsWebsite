package handlers

import (
	"avions-club/backend/database"
	"avions-club/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	Members  []models.Member  `json:"members"`
	Projects []models.Project `json:"projects"`
	Blogs    []models.Blog    `json:"blogs"`
}

func Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var response SearchResponse

	// Search in members
	database.DB.Where("name ILIKE ? OR position ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&response.Members)

	// Search in projects
	database.DB.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&response.Projects)

	// Search in blogs
	database.DB.Preload("Author").
		Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&response.Blogs)

	c.JSON(http.StatusOK, response)
}
