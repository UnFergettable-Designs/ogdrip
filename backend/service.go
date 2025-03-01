package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// APIResponse represents the structure of the API response
type APIResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ImageURL     string `json:"image_url,omitempty"`
	MetaTagsURL  string `json:"meta_tags_url,omitempty"`
	PreviewURL   string `json:"preview_url,omitempty"`
	ZipURL       string `json:"zip_url,omitempty"`      // URL to download files as zip
	HtmlContent  string `json:"html_content,omitempty"` // HTML content for direct display
}

// Config holds the service configuration
type Config struct {
	Port         string
	BaseURL      string
	OutputDir    string
	EnableCORS   bool
	MaxQueueSize int
	ChromePath   string
}

// Default configuration
var config = Config{
	Port:         "8888",
	BaseURL:      "http://localhost:8888",
	OutputDir:    "outputs",
	EnableCORS:   true,
	MaxQueueSize: 10,
	ChromePath:   "", // Will use system default if empty
}

// loadConfig loads configuration from environment variables
func loadConfig() {
	// Load configuration from environment variables
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
		log.Printf("Using PORT from environment: %s", port)
	}
	
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
		log.Printf("Using BASE_URL from environment: %s", baseURL)
	}
	
	if outputDir := os.Getenv("OUTPUT_DIR"); outputDir != "" {
		config.OutputDir = outputDir
		log.Printf("Using OUTPUT_DIR from environment: %s", outputDir)
	}
	
	if enableCORS := os.Getenv("ENABLE_CORS"); enableCORS != "" {
		config.EnableCORS = enableCORS == "true" || enableCORS == "1" || enableCORS == "yes"
		log.Printf("CORS %s based on environment setting", map[bool]string{true: "enabled", false: "disabled"}[config.EnableCORS])
	}
	
	if maxQueue := os.Getenv("MAX_QUEUE_SIZE"); maxQueue != "" {
		if val, err := strconv.Atoi(maxQueue); err == nil && val > 0 {
			config.MaxQueueSize = val
			log.Printf("Using MAX_QUEUE_SIZE from environment: %d", val)
		} else {
			log.Printf("Invalid MAX_QUEUE_SIZE value: %s, using default: %d", maxQueue, config.MaxQueueSize)
		}
	}
	
	if chromePath := os.Getenv("CHROME_PATH"); chromePath != "" {
		config.ChromePath = chromePath
		log.Printf("Using custom Chrome path from environment: %s", chromePath)
	}
	
	// Set logging level based on environment
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		switch strings.ToLower(logLevel) {
		case "debug":
			log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
			log.Printf("Debug logging enabled")
		case "info":
			log.SetFlags(log.Ldate | log.Ltime)
		case "warn", "warning":
			// Just use default logging format
		case "error":
			// Minimal logging
			log.SetFlags(log.Ltime)
		}
	}
	
	log.Printf("Configuration loaded: Port=%s, BaseURL=%s, OutputDir=%s, EnableCORS=%v, MaxQueueSize=%d",
		config.Port, config.BaseURL, config.OutputDir, config.EnableCORS, config.MaxQueueSize)
	
	// Make sure base URL doesn't end with a slash
	config.BaseURL = strings.TrimSuffix(config.BaseURL, "/")
	
	// Create OUTPUT_DIR if it doesn't exist
	if _, err := os.Stat(config.OutputDir); os.IsNotExist(err) {
		log.Printf("Creating output directory: %s", config.OutputDir)
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			log.Printf("WARNING: Failed to create output directory: %v", err)
		}
	}
}

// ServiceMain is the entry point for the API service
func ServiceMain() {
	// Set environment variable to indicate we're running in API service mode
	os.Setenv("OG_API_SERVICE", "true")
	
	// Load configuration from environment variables
	loadConfig()
	
	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Global CORS middleware applied to all requests
	globalCorsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the incoming request for debugging
			log.Printf("Received %s request from %s: %s %s", 
				r.Method,
				r.RemoteAddr, 
				r.Host,
				r.URL.Path)
			
			// Always apply CORS headers to all responses
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "3600")
			
			// Handle preflight OPTIONS requests immediately
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// Process the request
			next.ServeHTTP(w, r)
		})
	}

	// Create a new ServeMux to register handlers
	mux := http.NewServeMux()
	
	// Register individual API handlers 
	mux.HandleFunc("/api/generate", handleGenerateRequest)
	mux.HandleFunc("/api/health", handleHealthCheck)
	mux.HandleFunc("/api/download-zip", handleZipDownload)
	
	// Serve static files from the output directory
	fs := http.FileServer(http.Dir("."))
	mux.Handle("/outputs/", fs)
	
	// Simplified file access handler
	mux.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		// Extract filename from path
		filename := strings.TrimPrefix(r.URL.Path, "/files/")
		if filename == "" {
			http.Error(w, "No file specified", http.StatusBadRequest)
			return
		}
		
		// Construct file path
		filePath := filepath.Join(config.OutputDir, filename)
		
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		
		// Determine content type
		contentType := "application/octet-stream"
		if strings.HasSuffix(filePath, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(filePath, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
			contentType = "image/jpeg"
		}
		
		// Set content type and serve the file
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, filePath)
	})
	
	// Apply global CORS middleware to all routes
	handler := globalCorsMiddleware(mux)
	
	// Start the server
	log.Printf("Starting Open Graph API service on port %s...", config.Port)
	log.Printf("Base URL: %s", config.BaseURL)
	log.Printf("Files will be served from %s/files/{filename}", config.BaseURL)
	log.Printf("CORS Enabled: Applying to all requests")
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}

// handleHealthCheck responds to health check requests
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Message: "Open Graph Generator API is running",
	}
	
	sendJSONResponse(w, response)
}

// handleZipDownload creates a zip file with the specified files and serves it for download
func handleZipDownload(w http.ResponseWriter, r *http.Request) {
	// Get the filenames from the query parameters
	files := r.URL.Query()["file"]
	if len(files) == 0 {
		http.Error(w, "No files specified", http.StatusBadRequest)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=opengraph_assets.zip")

	// Create a zip writer
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Add each file to the zip
	for _, filename := range files {
		// Validate filename to prevent directory traversal
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
			continue // Skip suspicious filenames
		}

		// Full path to the file
		filePath := filepath.Join(config.OutputDir, filename)

		// Check if file exists
		fileInfo, err := os.Stat(filePath)
		if os.IsNotExist(err) || fileInfo.IsDir() {
			continue // Skip files that don't exist or are directories
		}

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			continue
		}
		defer file.Close()

		// Create a file header
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			log.Printf("Error creating file header for %s: %v", filePath, err)
			continue
		}

		// Set compression
		header.Method = zip.Deflate

		// Create file in zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Printf("Error creating file in zip for %s: %v", filePath, err)
			continue
		}

		// Copy file contents to zip
		_, err = io.Copy(writer, file)
		if err != nil {
			log.Printf("Error writing file contents to zip for %s: %v", filePath, err)
			continue
		}
	}
}

// handleGenerateRequest processes requests to generate Open Graph assets
func handleGenerateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log the incoming request
	log.Printf("Received generate request from %s", r.RemoteAddr)
	
	// Parse the multipart form with a reasonable max memory
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		// If not multipart, try to parse as regular form
		if err = r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			sendErrorResponse(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	}
	
	// Log the form data received
	log.Printf("Form data received:")
	for key, values := range r.Form {
		log.Printf("  %s: %v", key, values)
	}

	// Generate a unique ID for this request
	requestID := generateRequestID()
	
	// Prepare the output paths
	imgOutputPath := filepath.Join(config.OutputDir, requestID+"_og_image.png")
	htmlOutputPath := filepath.Join(config.OutputDir, requestID+"_og_meta.html")

	// Make sure the output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Printf("Error creating output directory: %v", err)
		sendErrorResponse(w, "Failed to create output directory", http.StatusInternalServerError)
		return
	}

	// Instead of executing a separate binary, set up a context for direct generation
	// Save original args and restore them after
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	
	// Build args array for flag parsing
	args := []string{"og-generator"}
	
	// Always set the API service flag
	args = append(args, "-api-service=true")
	
	// Set environment variable for API service mode
	os.Setenv("OG_API_SERVICE", "true")
	
	// Add all the form parameters as command-line args
	for key, values := range r.Form {
		if len(values) > 0 && values[0] != "" {
			// Map form fields to command line arguments
			cmdFlag := key
			cmdValue := values[0]
			
			// Special case mappings if needed
			if key == "targetUrl" {
				cmdFlag = "target-url"
			} else if key == "twitterCard" {
				cmdFlag = "twitter-card"
			} else if key == "url" {
				cmdFlag = "url"
				log.Printf("Found URL parameter: %s", cmdValue)
			} else if key == "title" {
				cmdFlag = "title"
				log.Printf("Found title parameter: %s", cmdValue)
			}
			
			// Handle boolean flags (checkboxes)
			if key == "debug" || key == "verbose" {
				if cmdValue == "true" || cmdValue == "on" || cmdValue == "1" {
					args = append(args, "-"+cmdFlag)
					continue
				}
			}
			
			args = append(args, "-"+cmdFlag+"="+cmdValue)
		}
	}
	
	// Add output paths
	args = append(args, "-output="+imgOutputPath)
	args = append(args, "-html="+htmlOutputPath)
	
	// Log what we're doing
	log.Printf("Generating with args: %s", strings.Join(args, " "))
	
	// Set args for flag parsing
	os.Args = args
	
	// Create a channel to capture errors from ServerMain
	errChan := make(chan error, 1)
	
	// Run ServerMain in a goroutine so we can capture if it exits
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("generator panicked: %v", r)
			} else {
				// Signal completion
				errChan <- nil
			}
		}()
		
		// Call ServerMain directly
		ServerMain()
	}()
	
	// Wait for generation to complete or time out
	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Error during generation: %v", err)
			sendErrorResponse(w, fmt.Sprintf("Failed to generate Open Graph assets: %v", err), http.StatusInternalServerError)
			return
		}
	case <-time.After(30 * time.Second):
		log.Printf("Generation timed out after 30 seconds")
		sendErrorResponse(w, "Generation timed out", http.StatusRequestTimeout)
		return
	}

	// Construct the response URLs
	imageURL := fmt.Sprintf("%s/files/%s", config.BaseURL, filepath.Base(imgOutputPath))
	metaTagsURL := fmt.Sprintf("%s/files/%s", config.BaseURL, filepath.Base(htmlOutputPath))
	zipURL := fmt.Sprintf("%s/api/download-zip?file=%s&file=%s", 
		config.BaseURL, 
		filepath.Base(imgOutputPath), 
		filepath.Base(htmlOutputPath))

	// Check if image was actually generated
	if _, err := os.Stat(imgOutputPath); os.IsNotExist(err) {
		log.Printf("Warning: Image file was not created at %s", imgOutputPath)
		imageURL = "" // Don't include image URL if file doesn't exist
	}

	// Check if HTML was actually generated
	htmlContent := ""
	if _, err := os.Stat(htmlOutputPath); os.IsNotExist(err) {
		log.Printf("Warning: HTML file was not created at %s", htmlOutputPath)
		metaTagsURL = "" // Don't include meta tags URL if file doesn't exist
	} else {
		// Read HTML content for direct display
		htmlBytes, err := os.ReadFile(htmlOutputPath)
		if err == nil {
			htmlContent = string(htmlBytes)
		}
	}

	// Send the successful response
	response := APIResponse{
		Success:     true,
		Message:     "Open Graph assets generated successfully",
		ImageURL:    imageURL,
		MetaTagsURL: metaTagsURL,
		ZipURL:      zipURL,
		HtmlContent: htmlContent,
	}

	// Add more information to the message if the files were generated
	if imageURL != "" && metaTagsURL != "" {
		response.Message = fmt.Sprintf("Open Graph assets generated successfully. You can view and download the files individually or as a zip archive.")
	} else if metaTagsURL != "" {
		response.Message = fmt.Sprintf("Meta tags HTML generated successfully. You can access it at %s", metaTagsURL)
	} else if imageURL != "" {
		response.Message = fmt.Sprintf("Open Graph image generated successfully. You can access it at %s", imageURL)
	}

	sendJSONResponse(w, response)
}

// sendJSONResponse sends a structured JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response with the specified status code
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := APIResponse{
		Success: false,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// generateRequestID creates a unique ID for each request
func generateRequestID() string {
	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 36)
	nano := strconv.FormatInt(int64(now.Nanosecond()), 36)
	pid := strconv.FormatInt(int64(os.Getpid()), 36)
	
	// Add a random component for additional uniqueness
	rand.Seed(now.UnixNano())
	random := strconv.FormatInt(rand.Int63(), 36)
	
	return fmt.Sprintf("%s%s_%s%s", timestamp, nano, pid, random)
} 