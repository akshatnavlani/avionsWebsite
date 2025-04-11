package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

var supabaseClient *storage_go.Client

// BucketInfo represents the structure of a bucket in Supabase
type BucketInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Public    bool   `json:"public"`
}

// verifyServiceKey checks if the service key is valid by making a simple API call
func verifyServiceKey(url, key string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/v1/", url), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("invalid service key (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// listBuckets makes a direct API call to list buckets
func listBuckets(url, key string) ([]BucketInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/storage/v1/bucket", url), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list buckets (status %d): %s", resp.StatusCode, string(body))
	}

	var buckets []BucketInfo
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return buckets, nil
}

// InitStorage initializes the Supabase storage client
func InitStorage() error {
	// Debug logging
	fmt.Println("Initializing Supabase storage...")
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")

	fmt.Printf("SUPABASE_URL: %s\n", supabaseURL)
	fmt.Printf("SUPABASE_SERVICE_KEY length: %d\n", len(supabaseKey))

	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("supabase credentials not found in environment variables")
	}

	// Verify service key first
	fmt.Println("Verifying service key...")
	if err := verifyServiceKey(supabaseURL, supabaseKey); err != nil {
		return fmt.Errorf("invalid service key: %v", err)
	}
	fmt.Println("Service key verified successfully")

	fmt.Println("Creating Supabase client...")
	supabaseClient = storage_go.NewClient(supabaseURL, supabaseKey, nil)

	// First, try to list buckets to verify permissions
	fmt.Println("Listing buckets to verify permissions...")
	buckets, err := listBuckets(supabaseURL, supabaseKey)
	if err != nil {
		fmt.Printf("Error listing buckets: %v\n", err)
		return fmt.Errorf("failed to list buckets: %v", err)
	}
	fmt.Printf("Successfully listed buckets: %+v\n", buckets)

	// Verify required buckets exist
	requiredBuckets := []string{"images", "markdown"}
	bucketMap := make(map[string]bool)
	for _, bucket := range buckets {
		bucketMap[bucket.Name] = true
	}

	for _, requiredBucket := range requiredBuckets {
		if !bucketMap[requiredBucket] {
			return fmt.Errorf("required bucket '%s' not found", requiredBucket)
		}
		fmt.Printf("Verified bucket '%s' exists and is public\n", requiredBucket)
	}

	return nil
}

func UploadFile(file *multipart.FileHeader, filename string) (string, error) {
	if supabaseClient == nil {
		return "", fmt.Errorf("supabase client is not initialized")
	}

	fmt.Printf("Starting upload for file: %s\n", filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer src.Close()

	// Read file content
	content, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	fmt.Printf("Read %d bytes from file\n", len(content))

	// Get bucket based on file type
	bucket := getBucketFromFilename(filename)
	fmt.Printf("Using bucket: %s\n", bucket)

	// Clean up the filename - remove any directory prefixes
	cleanFilename := filepath.Base(filename)
	fmt.Printf("Cleaned filename: %s\n", cleanFilename)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", cleanFilename)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		return "", fmt.Errorf("error copying file content: %v", err)
	}

	// Add bucket name
	if err := writer.WriteField("bucketName", bucket); err != nil {
		return "", fmt.Errorf("error writing bucket name: %v", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	// Create request
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", os.Getenv("SUPABASE_URL"), bucket, cleanFilename)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_KEY"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SUPABASE_SERVICE_KEY")))

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		os.Getenv("SUPABASE_URL"),
		bucket,
		cleanFilename,
	)

	fmt.Printf("Generated URL: %s\n", publicURL)
	return publicURL, nil
}

func DeleteFile(bucket, filename string) error {
	if bucket != "images" && bucket != "markdown" {
		return fmt.Errorf("invalid bucket: %s", bucket)
	}

	_, err := supabaseClient.RemoveFile(bucket, []string{filename})
	if err != nil {
		return fmt.Errorf("error deleting file from Supabase: %v", err)
	}

	return nil
}

func getBucketFromFilename(filename string) string {
	// Extract file type from path (e.g., "images/file.jpg" -> "images")
	parts := strings.Split(filename, "/")
	if len(parts) > 1 {
		switch parts[0] {
		case "images", "markdown":
			return parts[0]
		}
	}

	// Determine bucket by file extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "images"
	case ".md":
		return "markdown"
	default:
		return "images" // Default to images bucket
	}
}
