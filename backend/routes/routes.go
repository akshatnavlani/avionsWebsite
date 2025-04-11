package routes

import (
	"avions-club/backend/handlers"
	"avions-club/backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for our application
func SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Public routes
	r.GET("/api/members", handlers.GetMembers)
	r.GET("/api/members/:id", handlers.GetMember)
	r.GET("/api/projects", handlers.GetProjects)
	r.GET("/api/projects/:id", handlers.GetProject)
	r.GET("/api/blogs", handlers.GetBlogs)
	r.GET("/api/blogs/:id", handlers.GetBlog)
	r.GET("/api/search", handlers.Search)

	// Auth routes (public)
	r.POST("/api/auth/login", handlers.Login)

	// Protected routes
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Members
		protected.POST("/api/members", handlers.CreateMember)
		protected.PUT("/api/members/:id", handlers.UpdateMember)
		protected.DELETE("/api/members/:id", handlers.DeleteMember)

		// Projects
		protected.POST("/api/projects", handlers.CreateProject)
		protected.PUT("/api/projects/:id", handlers.UpdateProject)
		protected.DELETE("/api/projects/:id", handlers.DeleteProject)

		// Blogs
		protected.POST("/api/blogs", handlers.CreateBlog)
		protected.PUT("/api/blogs/:id", handlers.UpdateBlog)
		protected.DELETE("/api/blogs/:id", handlers.DeleteBlog)

		// Storage
		protected.POST("/api/storage/upload", handlers.UploadFile)
		protected.DELETE("/api/storage/:bucket/:filename", handlers.DeleteFile)
	}
}
