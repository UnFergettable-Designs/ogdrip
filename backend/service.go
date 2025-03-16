package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/cors"
)

// APIResponse represents the structure of the API response
type APIResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ImageURL    string `json:"image_url,omitempty"`
	MetaTagsURL string `json:"meta_tags_url,omitempty"`
	PreviewURL  string `json:"preview_url,omitempty"`
	ZipURL      string `json:"zip_url,omitempty"`      // URL to download files as zip
	HtmlContent string `json:"html_content,omitempty"` // HTML content for direct display
	ID          string `json:"id,omitempty"`
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

// Global variable to track if Sentry is initialized
var sentryInitialized bool

// Add package-level db variable
var db *Database

// GenerateRequest represents a request to generate Open Graph assets
type GenerateRequest struct {
	URL          string            `json:"url"`
	Title        string            `json:"title,omitempty"`
	Description  string            `json:"description,omitempty"`
	ImageWidth   int               `json:"image_width,omitempty"`
	ImageHeight  int               `json:"image_height,omitempty"`
	CustomParams map[string]string `json:"custom_params,omitempty"`
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

	// Initialize Sentry for error tracking
	if err := InitSentry(); err != nil {
		log.Printf("Warning: Failed to initialize Sentry: %v", err)
		log.Printf("Error reporting will be limited to logs only")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Initialize the database - global variable that will be reused
	var err error
	db, err = InitDB()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Printf("Generation history tracking will be disabled")
		// Report to Sentry
		CaptureException(err)
	} else {
		log.Printf("Database initialized successfully")
		// Start the cleanup task but don't close the database connection afterward
		// as it will be needed by other operations
		StartCleanupTask(db)
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

	// Add new API endpoints for history
	mux.HandleFunc("/api/history", handleHistoryRequest)
	mux.HandleFunc("/api/generation/", handleGetGenerationRequest)
	mux.HandleFunc("/api/download-complete", handleDownloadCompleteRequest)

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

	// Create middleware chain: CORS -> Sentry -> application handlers
	handler := globalCorsMiddleware(SentryMiddleware(mux))

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = config.Port // Use the config port
	} else {
		config.Port = port // Update config if environment variable is set
	}

	log.Printf("Starting Open Graph API service on port %s...", config.Port)
	log.Printf("Base URL: %s", config.BaseURL)
	log.Printf("Files will be served from %s/files/{filename}", config.BaseURL)
	log.Printf("CORS Enabled: Applying to all requests")
	log.Printf("Sentry Error Tracking: Enabled")
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

	// Create generation record in database with initial pending status
	generation := &Generation{
		ID:          requestID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		TargetURL:   r.FormValue("target_url"),
		ImagePath:   imgOutputPath,
		HTMLPath:    htmlOutputPath,
		CreatedAt:   time.Now(),
		ClientIP:    r.RemoteAddr,
		UserAgent:   r.UserAgent(),
		Parameters:  parametersToJSON(r.Form),
		Status:      "pending",
	}

	// Save initial generation record
	if err := db.SaveGeneration(generation); err != nil {
		log.Printf("Error saving generation record: %v", err)
		captureError(err, map[string]interface{}{
			"operation": "save_generation_record",
			"requestID": requestID,
		})
		// Continue anyway, as this is just for tracking
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

			// Report error to Sentry with context
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("request_type", "generation")
				scope.SetExtra("request_ip", r.RemoteAddr)
				scope.SetExtra("user_agent", r.UserAgent())

				// Add form data as context (filtering sensitive information)
				formData := make(map[string]string)
				for key, values := range r.Form {
					if len(values) > 0 {
						// Don't include large or sensitive values
						if key != "image" && key != "password" && len(values[0]) < 1000 {
							formData[key] = values[0]
						}
					}
				}
				// Convert to map[string]interface{} for Sentry context
				formDataContext := make(map[string]interface{})
				for k, v := range formData {
					formDataContext[k] = v
				}
				scope.SetContext("form_data", formDataContext)

				// Capture the exception
				CaptureException(err)
			})

			// Save error to database if database is available
			db, dbErr := InitDB()
			if dbErr == nil {
				// Try to extract request ID from args
				requestID := ""
				for _, arg := range args {
					if strings.HasPrefix(arg, "-output=") {
						// Extract ID from output path
						outputPath := strings.TrimPrefix(arg, "-output=")
						requestID = strings.TrimSuffix(filepath.Base(outputPath), "_og_image.png")
						break
					}
				}

				if requestID != "" {
					// Set error message and status
					db.SetErrorMessage(requestID, err.Error())
				}
			}

			sendErrorResponse(w, fmt.Sprintf("Failed to generate Open Graph assets: %v", err), http.StatusInternalServerError)
			return
		}
	case <-time.After(30 * time.Second):
		log.Printf("Generation timed out after 30 seconds")

		// Report timeout to Sentry
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("error_type", "timeout")
			scope.SetTag("request_type", "generation")
			scope.SetExtra("timeout_duration", "30s")
			CaptureMessage("Generation request timed out")
		})

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

	// Determine if generation was successful
	generationSuccess := imageURL != "" || metaTagsURL != ""

	// Update the generation status in the database
	if generationSuccess {
		// Update generation status to completed
		if err := db.UpdateGenerationStatus(requestID, "completed", ""); err != nil {
			log.Printf("Error updating generation status: %v", err)
			captureError(err, map[string]interface{}{
				"operation": "update_generation_status",
				"requestID": requestID,
				"status":    "completed",
			})
		}
	} else {
		// Update generation status to failed
		errorMsg := "Failed to generate Open Graph assets"
		if err := db.UpdateGenerationStatus(requestID, "failed", errorMsg); err != nil {
			log.Printf("Error updating generation status: %v", err)
			captureError(err, map[string]interface{}{
				"operation": "update_generation_status",
				"requestID": requestID,
				"status":    "failed",
			})
		}
	}

	// Send the successful response
	response := APIResponse{
		Success:     generationSuccess,
		Message:     "Open Graph assets generated successfully",
		ImageURL:    imageURL,
		MetaTagsURL: metaTagsURL,
		ZipURL:      zipURL,
		HtmlContent: htmlContent,
		ID:          requestID, // Include the ID in the response
	}

	// Add more information to the message if the files were generated
	if imageURL != "" && metaTagsURL != "" {
		response.Message = fmt.Sprintf("Open Graph assets generated successfully. You can view and download the files individually or as a zip archive.")
	} else if metaTagsURL != "" {
		response.Message = fmt.Sprintf("Meta tags HTML generated successfully. You can access it at %s", metaTagsURL)
	} else if imageURL != "" {
		response.Message = fmt.Sprintf("Open Graph image generated successfully. You can access it at %s", imageURL)
	} else {
		response.Message = "Failed to generate Open Graph assets. Please check your input parameters and try again."
		response.Success = false
	}

	sendJSONResponse(w, response)
}

// parametersToJSON converts form values to a JSON string
func parametersToJSON(form url.Values) string {
	params := make(map[string]interface{})

	for key, values := range form {
		if len(values) == 1 {
			params[key] = values[0]
		} else if len(values) > 1 {
			params[key] = values
		}
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Printf("Error marshaling parameters to JSON: %v", err)
		return "{}"
	}

	return string(jsonData)
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

// handleHistoryRequest handles requests to get the history of generation requests
func handleHistoryRequest(w http.ResponseWriter, r *http.Request) {
	// Set headers for CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limit := 50 // Default limit
	offset := 0 // Default offset

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get generations from database
	generations, err := db.GetRecentGenerations(limit, offset)
	if err != nil {
		log.Printf("Error getting generations: %v", err)
		http.Error(w, "Failed to retrieve generations", http.StatusInternalServerError)
		return
	}

	// Get the total count
	count, err := db.GetGenerationCount()
	if err != nil {
		log.Printf("Error getting generation count: %v", err)
		// Continue anyway, just set count to 0
		count = 0
	}

	// Create response object
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"generations": generations,
			"total":       count,
			"limit":       limit,
			"offset":      offset,
		},
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleGetGenerationRequest returns details for a specific generation
func handleGetGenerationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the generation ID from the URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		sendErrorResponse(w, "Invalid generation ID", http.StatusBadRequest)
		return
	}
	generationID := parts[len(parts)-1]

	// Get the database instance
	db, err := InitDB()
	if err != nil {
		sendErrorResponse(w, "Database not available", http.StatusInternalServerError)
		return
	}

	// Get the generation details
	generation, err := db.GetGeneration(generationID)
	if err != nil {
		log.Printf("Error retrieving generation %s: %v", generationID, err)
		sendErrorResponse(w, "Generation not found", http.StatusNotFound)
		return
	}

	// Construct the response URLs
	imageURL := fmt.Sprintf("%s/files/%s", config.BaseURL, filepath.Base(generation.ImagePath))
	metaTagsURL := fmt.Sprintf("%s/files/%s", config.BaseURL, filepath.Base(generation.HTMLPath))
	zipURL := fmt.Sprintf("%s/api/download-zip?file=%s&file=%s",
		config.BaseURL,
		filepath.Base(generation.ImagePath),
		filepath.Base(generation.HTMLPath))

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message":    "Generation retrieved successfully",
		"generation": generation,
		"image_url":  imageURL,
		"meta_url":   metaTagsURL,
		"zip_url":    zipURL,
	})
}

// handleDownloadCompleteRequest marks a generation as downloaded
func handleDownloadCompleteRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var requestData struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.ID == "" {
		sendErrorResponse(w, "Generation ID is required", http.StatusBadRequest)
		return
	}

	// Get the database instance
	db, err := InitDB()
	if err != nil {
		sendErrorResponse(w, "Database not available", http.StatusInternalServerError)
		return
	}

	// Mark the generation as downloaded
	if err := db.MarkAsDownloaded(requestData.ID); err != nil {
		log.Printf("Error marking generation %s as downloaded: %v", requestData.ID, err)
		sendErrorResponse(w, "Failed to update generation", http.StatusInternalServerError)
		return
	}

	// Set cleanup time to 1 hour from now for downloaded items
	cleanupTime := time.Now().Add(1 * time.Hour)
	if err := db.SetCleanupTime(requestData.ID, cleanupTime); err != nil {
		log.Printf("Warning: Failed to set cleanup time for generation %s: %v", requestData.ID, err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Generation marked as downloaded",
	})
}

// handleGenerationDetailsRequest gets details for a specific generation
func handleGenerationDetailsRequest(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract generation ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid request: missing generation ID", http.StatusBadRequest)
		return
	}

	generationID := pathParts[3]
	if generationID == "" {
		http.Error(w, "Invalid request: empty generation ID", http.StatusBadRequest)
		return
	}

	// Get the generation details from database
	generation, err := db.GetGenerationByID(generationID)
	if err != nil {
		log.Printf("Error getting generation %s: %v", generationID, err)
		http.Error(w, "Failed to retrieve generation details", http.StatusInternalServerError)
		return
	}

	if generation == nil {
		http.Error(w, "Generation not found", http.StatusNotFound)
		return
	}

	// Create response object
	response := map[string]interface{}{
		"success": true,
		"data":    generation,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleAdminVerify validates an admin token against the environment variable
func handleAdminVerify(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Extract the token from the Authorization header
	// Format should be "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		sendErrorResponse(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	token := parts[1]

	// Get the admin token from environment variable
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		// If ADMIN_TOKEN is not set, admin access is disabled
		sendErrorResponse(w, "Admin access is disabled", http.StatusUnauthorized)
		return
	}

	// Validate the token
	if token != adminToken {
		sendErrorResponse(w, "Invalid admin token", http.StatusUnauthorized)
		return
	}

	// If the token is valid, send a success response
	sendJSONResponse(w, APIResponse{
		Success: true,
		Message: "Admin authentication successful",
	})
}

// verifyAdminToken middleware checks if the request has a valid admin token
func verifyAdminToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the admin token from environment variable
		adminToken := os.Getenv("ADMIN_TOKEN")

		// If ADMIN_TOKEN is not set, skip admin auth (for backward compatibility)
		if adminToken == "" {
			next(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendErrorResponse(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendErrorResponse(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate the token
		if token != adminToken {
			sendErrorResponse(w, "Invalid admin token", http.StatusUnauthorized)
			return
		}

		// If the token is valid, call the next handler
		next(w, r)
	}
}

// StartAPIService starts the API service on the specified port
func StartAPIService(port string) {
	log.Printf("Starting API service on port %s", port)

	// Initialize service components
	initService()

	// Default to port 8888 if not provided
	if port == "" {
		port = "8888"
	}

	// Get port from environment if available
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
		log.Printf("Using port from environment: %s", port)
	}

	mux := http.NewServeMux()

	// Set up static file serving
	setupStaticFileServing(mux)

	// API endpoints
	mux.HandleFunc("/api/generate", handleGenerateRequest)
	mux.HandleFunc("/api/get/", handleGetGenerationRequest)
	mux.HandleFunc("/api/download/", handleDownloadRequest)
	mux.HandleFunc("/api/health", handleHealthCheck)
	mux.HandleFunc("/api/admin/verify", handleAdminVerify)

	// History endpoints with admin auth
	mux.HandleFunc("/api/history", verifyAdminToken(handleHistoryRequest))
	mux.HandleFunc("/api/history/", verifyAdminToken(handleGenerationDetailsRequest))

	// Set up Swagger UI for API documentation
	setupSwagger(mux)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allowing all origins for now
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	// Start the server
	handler := c.Handler(mux)

	// Add Sentry middleware if available
	if sentryInitialized {
		handler = sentryHandler(handler)
	}

	log.Printf("API Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// captureError sends an error to Sentry with additional context
func captureError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		if context != nil {
			for k, v := range context {
				scope.SetExtra(k, v)
			}
		}
		sentry.CaptureException(err)
	})
}

// handleDownloadRequest handles requests to download generated files
func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	// Extract the filename from the URL path
	filename := strings.TrimPrefix(r.URL.Path, "/api/download/")
	if filename == "" {
		http.Error(w, "No filename specified", http.StatusBadRequest)
		return
	}

	// Build the full path to the file
	outputDir := getOutputDir()
	filePath := filepath.Join(outputDir, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// sentryHandler wraps an http handler with Sentry error tracking
func sentryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a placeholder for actual Sentry implementation
		// In a real implementation, this would use the Sentry SDK to track errors
		h.ServeHTTP(w, r)
	})
}

// setupStaticFileServing configures paths for serving static files
func setupStaticFileServing(mux *http.ServeMux) {
	// Get the output directory
	outputDir := getOutputDir()

	// Create a file server for the output directory
	outputFileServer := http.FileServer(http.Dir(outputDir))

	// Set up a handler for the outputs path
	mux.Handle("/outputs/", http.StripPrefix("/outputs/", outputFileServer))
}

// buildGeneratorArgs builds the arguments for the generator command
func buildGeneratorArgs(req GenerateRequest, imgPath, htmlPath string) []string {
	args := []string{}

	// Add required parameters
	if req.URL != "" {
		args = append(args, fmt.Sprintf("-url=%s", req.URL))
	}

	// Add output path
	args = append(args, fmt.Sprintf("-output=%s", imgPath))

	// Add HTML output path
	args = append(args, fmt.Sprintf("-html=%s", htmlPath))

	// Add optional parameters if provided
	if req.Title != "" {
		args = append(args, fmt.Sprintf("-title=%s", req.Title))
	}

	if req.Description != "" {
		args = append(args, fmt.Sprintf("-description=%s", req.Description))
	}

	if req.ImageWidth > 0 {
		args = append(args, fmt.Sprintf("-width=%d", req.ImageWidth))
	}

	if req.ImageHeight > 0 {
		args = append(args, fmt.Sprintf("-height=%d", req.ImageHeight))
	}

	// Add any custom parameters from the request
	for key, value := range req.CustomParams {
		args = append(args, fmt.Sprintf("-%s=%s", key, value))
	}

	return args
}

// getOutputDir returns the configured output directory
func getOutputDir() string {
	// Use the configured output directory or default to ./outputs
	if config.OutputDir != "" {
		return config.OutputDir
	}
	return "./outputs"
}

// StartCleanupTask starts the periodic cleanup task
func StartCleanupTask(db *Database) {
	// Run cleanup every hour
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		// Run cleanup on the ticker interval
		for range ticker.C {
			log.Printf("Running scheduled cleanup")
			if err := db.RunCleanup(); err != nil {
				log.Printf("Error running scheduled cleanup: %v", err)
			}
		}
	}()
}

