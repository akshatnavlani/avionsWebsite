package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"avions-club/backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ImageUploadResponse represents the response for image uploads
type ImageUploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

// ProcessedContent represents the processed markdown content with image URLs
type ProcessedContent struct {
	Content string            `json:"content"`
	Images  map[string]string `json:"images"`
}

// allowedImageTypes defines the allowed image MIME types
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// allowedMarkdownTypes defines the allowed markdown MIME types
var allowedMarkdownTypes = map[string]bool{
	"text/markdown":          true,
	"text/x-markdown":        true,
	"application/x-markdown": true,
}

// UploadFile handles file uploads
func UploadFile(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Get file type from form
	fileType := c.PostForm("type")
	if fileType == "" {
		// Try to determine type from file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			fileType = "image"
		case ".md":
			fileType = "markdown"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}
	}

	// Check file size
	if file.Size > 5<<20 { // 5MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Upload to storage
	url, err := storage.UploadFile(file, filename)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, ImageUploadResponse{
		URL:      url,
		Filename: filename,
	})
}

// ProcessMarkdownContent processes markdown content and handles image uploads
func ProcessMarkdownContent(c *gin.Context) {
	var content struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create a map to store image references and their URLs
	imageMap := make(map[string]string)

	// Find all image references in markdown
	imageRegex := regexp.MustCompile(`!\[.*?\]\((image/.*?)\)`)
	matches := imageRegex.FindAllStringSubmatch(content.Content, -1)

	// Process each image reference
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		imagePath := match[1]
		imageName := filepath.Base(imagePath)

		// Check if image exists in test files
		testFilePath := filepath.Join("test_files", imageName)
		if _, err := os.Stat(testFilePath); err == nil {
			// Open the image file
			file, err := os.Open(testFilePath)
			if err != nil {
				log.Printf("Error opening image file: %v", err)
				continue
			}

			// Create a multipart file header
			fileHeader := &multipart.FileHeader{
				Filename: imageName,
				Size:     0, // We'll get this from the file
			}

			// Upload the image
			url, err := storage.UploadFile(fileHeader, imageName)
			if err != nil {
				log.Printf("Error uploading image: %v", err)
				file.Close()
				continue
			}

			// Store the URL mapping
			imageMap[imagePath] = url
			file.Close()
		}
	}

	// Replace image references with URLs in content
	processedContent := content.Content
	for oldPath, newURL := range imageMap {
		processedContent = strings.ReplaceAll(processedContent, oldPath, newURL)
	}

	// Return processed content and image mappings
	c.JSON(http.StatusOK, ProcessedContent{
		Content: processedContent,
		Images:  imageMap,
	})
}

// DeleteFile handles file deletion from Supabase Storage
func DeleteFile(c *gin.Context) {
	bucket := c.Param("bucket")
	filename := c.Param("filename")

	if bucket != "images" && bucket != "markdown" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bucket"})
		return
	}

	err := storage.DeleteFile(bucket, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
