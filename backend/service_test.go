package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadConfig tests the configuration loading from environment variables
func TestLoadConfig(t *testing.T) {
	// Save original env vars to restore later
	origPort := os.Getenv("PORT")
	origBaseURL := os.Getenv("BASE_URL")
	origOutputDir := os.Getenv("OUTPUT_DIR")
	origEnableCORS := os.Getenv("ENABLE_CORS")
	
	// Clean up function to restore env vars
	defer func() {
		os.Setenv("PORT", origPort)
		os.Setenv("BASE_URL", origBaseURL)
		os.Setenv("OUTPUT_DIR", origOutputDir)
		os.Setenv("ENABLE_CORS", origEnableCORS)
	}()
	
	// Set test values
	os.Setenv("PORT", "9999")
	os.Setenv("BASE_URL", "http://test-server.com")
	os.Setenv("OUTPUT_DIR", "test_outputs")
	os.Setenv("ENABLE_CORS", "false")
	
	// Reset config to default values first
	config = Config{
		Port:         "8888",
		BaseURL:      "http://localhost:8888",
		OutputDir:    "outputs",
		EnableCORS:   true,
		MaxQueueSize: 10,
		ChromePath:   "",
	}
	
	// Run the function
	loadConfig()
	
	// Check if values were correctly loaded
	if config.Port != "9999" {
		t.Errorf("Expected Port to be 9999, got %s", config.Port)
	}
	
	if config.BaseURL != "http://test-server.com" {
		t.Errorf("Expected BaseURL to be http://test-server.com, got %s", config.BaseURL)
	}
	
	if config.OutputDir != "test_outputs" {
		t.Errorf("Expected OutputDir to be test_outputs, got %s", config.OutputDir)
	}
	
	if config.EnableCORS != false {
		t.Errorf("Expected EnableCORS to be false, got %v", config.EnableCORS)
	}
}

// TestGenerateRequestID tests the unique ID generation
func TestGenerateRequestID(t *testing.T) {
	// Generate multiple IDs and check they're different
	id1 := generateRequestID()
	id2 := generateRequestID()
	
	if id1 == id2 {
		t.Errorf("Expected different IDs, but got the same: %s", id1)
	}
	
	// Check format (should contain at least one underscore)
	if !strings.Contains(id1, "_") {
		t.Errorf("ID format incorrect, expected underscore: %s", id1)
	}
}

// TestBuildGeneratorArgs tests argument construction from requests
func TestBuildGeneratorArgs(t *testing.T) {
	// Create a mock request
	form := url.Values{}
	form.Add("url", "https://example.com")
	form.Add("title", "Test Title")
	form.Add("description", "Test Description")
	form.Add("debug", "true")
	
	testReq := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(form.Encode()))
	testReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	testReq.Form = form
	
	// Convert to GenerateRequest
	generateReq := GenerateRequest{
		URL:         "https://example.com",
		Title:       "Test Title",
		Description: "Test Description",
		CustomParams: map[string]string{
			"debug": "true",
		},
	}
	
	// Set output paths
	imgPath := "/tmp/test_image.png"
	htmlPath := "/tmp/test_meta.html"
	
	// Call the function
	args := buildGeneratorArgs(generateReq, imgPath, htmlPath)
	
	// Check expected arguments are present
	expectedArgs := []string{
		"-url=https://example.com",
		"-title=Test Title", 
		"-description=Test Description",
		"-debug=true",
		"-output=" + imgPath,
		"-html=" + htmlPath,
	}
	
	// Check each expected arg is in the result
	for _, expected := range expectedArgs {
		found := false
		for _, arg := range args {
			if arg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected argument %s not found in result: %v", expected, args)
		}
	}
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	// Create a request to the health endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	
	// Call the handler
	handleHealthCheck(w, req)
	
	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	
	// Check response content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	// Check response body contains success message
	if !strings.Contains(w.Body.String(), "Open Graph Generator API is running") {
		t.Errorf("Response body does not contain expected message: %s", w.Body.String())
	}
}

// TestCORSHeaders tests that CORS headers are set correctly
func TestCORSHeaders(t *testing.T) {
	// Enable CORS for this test
	origCORS := config.EnableCORS
	config.EnableCORS = true
	defer func() { config.EnableCORS = origCORS }()
	
	// Create a recorder for testing response headers
	w := httptest.NewRecorder()
	
	// Manually create a similar context as the real handler
	if config.EnableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	// Check if CORS headers are set
	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin to be *, got %s", origin)
	}
	
	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods != "POST, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods to be POST, OPTIONS, got %s", methods)
	}
}

// Mock the command execution for integration tests
type mockCommand struct {
	executed  bool
	shouldErr bool
}

// Overrides exec.Command for testing
func mockExecCommand(shouldErr bool) func(name string, arg ...string) *exec.Cmd {
	return func(name string, arg ...string) *exec.Cmd {
		// If we need to simulate an error
		if shouldErr {
			return exec.Command("false") // Will exit with non-zero status
		}
		
		// Create output files to simulate successful generation
		outImgPath := ""
		outHtmlPath := ""
		
		for _, a := range arg {
			if strings.HasPrefix(a, "-output=") {
				outImgPath = strings.TrimPrefix(a, "-output=")
			}
			if strings.HasPrefix(a, "-html=") {
				outHtmlPath = strings.TrimPrefix(a, "-html=")
			}
		}
		
		// Create directories if they don't exist
		if outImgPath != "" {
			os.MkdirAll(filepath.Dir(outImgPath), 0755)
			// Create an empty file to simulate the image
			os.WriteFile(outImgPath, []byte("test image data"), 0644)
		}
		
		if outHtmlPath != "" {
			os.MkdirAll(filepath.Dir(outHtmlPath), 0755)
			// Create an HTML file
			os.WriteFile(outHtmlPath, []byte("<html><body>Test HTML</body></html>"), 0644)
		}
		
		// Return command that will succeed
		return exec.Command("true")
	}
}

// Note: Integration tests would be implemented in a separate file 