package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
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
	}
	
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}
	
	if outputDir := os.Getenv("OUTPUT_DIR"); outputDir != "" {
		config.OutputDir = outputDir
	}
	
	if enableCORS := os.Getenv("ENABLE_CORS"); enableCORS != "" {
		config.EnableCORS = enableCORS == "true" || enableCORS == "1"
	}
	
	if maxQueue := os.Getenv("MAX_QUEUE_SIZE"); maxQueue != "" {
		if val, err := strconv.Atoi(maxQueue); err == nil && val > 0 {
			config.MaxQueueSize = val
		}
	}
	
	if chromePath := os.Getenv("CHROME_PATH"); chromePath != "" {
		config.ChromePath = chromePath
	}
	
	log.Printf("Configuration loaded: Port=%s, BaseURL=%s, OutputDir=%s, EnableCORS=%v, MaxQueueSize=%d",
		config.Port, config.BaseURL, config.OutputDir, config.EnableCORS, config.MaxQueueSize)
	
	if config.ChromePath != "" {
		log.Printf("Using custom Chrome path: %s", config.ChromePath)
	}
}

// ServiceMain is the entry point for the API service
func ServiceMain() {
	// Load configuration from environment variables
	loadConfig()
	
	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Define API routes
	http.HandleFunc("/api/generate", handleGenerateRequest)
	http.HandleFunc("/api/health", handleHealthCheck)
	
	// Serve static files from the output directory
	fs := http.FileServer(http.Dir(config.OutputDir))
	http.Handle("/outputs/", http.StripPrefix("/outputs/", fs))

	// Start the server
	log.Printf("Starting Open Graph API service on port %s...", config.Port)
	log.Printf("Base URL: %s", config.BaseURL)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

// handleHealthCheck responds to health check requests
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Message: "Open Graph Generator API is running",
	}
	
	sendJSONResponse(w, response)
}

// handleGenerateRequest processes requests to generate Open Graph assets
func handleGenerateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Enable CORS if configured
	if config.EnableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse request parameters
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for this request
	requestID := generateRequestID()
	
	// Prepare the output paths
	imgOutputPath := filepath.Join(config.OutputDir, requestID+"_og_image.png")
	htmlOutputPath := filepath.Join(config.OutputDir, requestID+"_og_meta.html")

	// Build command to run the original generator
	args := buildGeneratorArgs(r, imgOutputPath, htmlOutputPath)
	
	// Log the command for debugging
	log.Printf("Executing: ./og-generator %s", strings.Join(args, " "))
	
	// Execute the command using the pre-compiled binary
	cmd := exec.Command("./og-generator", args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		log.Printf("Error executing generator: %v\nOutput: %s", err, string(output))
		sendErrorResponse(w, fmt.Sprintf("Failed to generate Open Graph assets: %v", err), http.StatusInternalServerError)
		return
	}

	// Construct the response URLs
	imageURL := fmt.Sprintf("%s/outputs/%s", config.BaseURL, filepath.Base(imgOutputPath))
	metaTagsURL := fmt.Sprintf("%s/outputs/%s", config.BaseURL, filepath.Base(htmlOutputPath))

	// Check if image was actually generated
	if _, err := os.Stat(imgOutputPath); os.IsNotExist(err) {
		log.Printf("Warning: Image file was not created at %s", imgOutputPath)
		imageURL = "" // Don't include image URL if file doesn't exist
	}

	// Check if HTML was actually generated
	if _, err := os.Stat(htmlOutputPath); os.IsNotExist(err) {
		log.Printf("Warning: HTML file was not created at %s", htmlOutputPath)
		metaTagsURL = "" // Don't include meta tags URL if file doesn't exist
	}

	// Send the successful response
	response := APIResponse{
		Success:     true,
		Message:     "Open Graph assets generated successfully",
		ImageURL:    imageURL,
		MetaTagsURL: metaTagsURL,
	}
	
	sendJSONResponse(w, response)
}

// buildGeneratorArgs constructs the command line arguments for the generator
func buildGeneratorArgs(r *http.Request, imgOutputPath, htmlOutputPath string) []string {
	args := []string{}

	// Map form fields to command line arguments
	paramMap := map[string]string{
		"url":          "-url",
		"title":        "-title",
		"description":  "-description",
		"type":         "-type",
		"site":         "-site",
		"target-url":   "-target-url",
		"width":        "-width",
		"height":       "-height",
		"twitter-card": "-twitter-card",
		"wait":         "-wait",
		"selector":     "-selector",
		"quality":      "-quality",
	}

	// Add each parameter if it exists in the request
	for formParam, cmdArg := range paramMap {
		if value := r.FormValue(formParam); value != "" {
			args = append(args, cmdArg+"="+value)
		}
	}

	// Add boolean flags
	if r.FormValue("debug") == "true" {
		args = append(args, "-debug=true")
	}
	
	if r.FormValue("verbose") == "true" {
		args = append(args, "-verbose=true")
	}

	// Always set the output paths
	args = append(args, "-output="+imgOutputPath)
	args = append(args, "-html="+htmlOutputPath)

	return args
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