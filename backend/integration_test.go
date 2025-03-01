// +build integration

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// These tests require a working instance of Chrome/Chromium
// Run with: go test -tags=integration

// TestGenerateImageFromWebpage tests the full flow of generating an image from a webpage
func TestGenerateImageFromWebpage(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Create temp directory for test outputs
	testDir, err := os.MkdirTemp("", "opengraph-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)
	
	// Set test config
	origOutputDir := config.OutputDir
	config.OutputDir = testDir
	defer func() { config.OutputDir = origOutputDir }()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			handleGenerateRequest(w, r)
		} else if r.URL.Path == "/api/health" {
			handleHealthCheck(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	
	// Set the base URL for the API
	config.BaseURL = server.URL
	
	// Create a multipart form with test data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	
	// Add form fields
	fields := map[string]string{
		"url":         "https://example.com",
		"title":       "Integration Test",
		"description": "Testing the Open Graph generator",
		"debug":       "true",
	}
	
	for field, value := range fields {
		fw, err := writer.CreateFormField(field)
		if err != nil {
			t.Fatalf("Error creating form field: %v", err)
		}
		fw.Write([]byte(value))
	}
	
	// Close the writer
	writer.Close()
	
	// Create the HTTP request
	req, err := http.NewRequest("POST", server.URL+"/api/generate", &b)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	
	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Create an HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error executing request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
		
		// Read error body
		body, _ := ioutil.ReadAll(resp.Body)
		t.Logf("Response body: %s", string(body))
		return
	}
	
	// Parse the response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Error parsing response JSON: %v", err)
	}
	
	// Check response fields
	if !apiResp.Success {
		t.Errorf("Expected success=true, got false. Message: %s", apiResp.Message)
	}
	
	// Note: In a real integration test, we'd check for actual files
	// but this may be difficult depending on how the server is configured
	t.Logf("Response: %+v", apiResp)
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheckIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(handleHealthCheck))
	defer server.Close()
	
	// Make request to the health check endpoint
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Error making health check request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	// Parse the response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Error parsing response JSON: %v", err)
	}
	
	// Check response fields
	if !apiResp.Success {
		t.Errorf("Expected success=true, got false")
	}
	
	if apiResp.Message != "Open Graph Generator API is running" {
		t.Errorf("Unexpected message: %s", apiResp.Message)
	}
}

// TestRunGeneratorDirectly tests running the generator executable directly
func TestRunGeneratorDirectly(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Create temp directory for test outputs
	testDir, err := os.MkdirTemp("", "opengraph-direct-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)
	
	// Set up output paths
	imgPath := filepath.Join(testDir, "test-image.png")
	htmlPath := filepath.Join(testDir, "test-html.html")
	
	// Create a simple static HTML test server
	testHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
    <meta name="description" content="A test page for Open Graph generator">
</head>
<body>
    <h1>Test Page</h1>
    <p>This is a test page for the Open Graph generator integration test.</p>
</body>
</html>`
	
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTML))
	}))
	defer testServer.Close()
	
	// Build the command
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "go", "run", "server.go",
		"-url="+testServer.URL,
		"-output="+imgPath,
		"-html="+htmlPath,
		"-title=Test Title",
		"-description=Test Description",
		"-wait=2000", // Shorter wait time for tests
		"-debug=true",
	)
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Run the command
	err = cmd.Run()
	t.Logf("Command stdout: %s", stdout.String())
	t.Logf("Command stderr: %s", stderr.String())
	
	if err != nil {
		t.Fatalf("Error running generator: %v", err)
	}
	
	// Check if the image was created
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		t.Errorf("Expected image file wasn't created at: %s", imgPath)
	} else {
		// Get image file size to verify it's not empty
		fi, err := os.Stat(imgPath)
		if err == nil {
			t.Logf("Image file size: %d bytes", fi.Size())
			if fi.Size() < 100 { // Arbitrary small size check
				t.Errorf("Image file is too small: %d bytes", fi.Size())
			}
		}
	}
	
	// Check if the HTML was created
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Errorf("Expected HTML file wasn't created at: %s", htmlPath)
	} else {
		// Verify HTML content
		content, err := ioutil.ReadFile(htmlPath)
		if err == nil {
			htmlStr := string(content)
			
			// Check for expected elements in the HTML
			expectations := []string{
				"Test Title", // Our custom title
				"Test Description", // Our custom description
				"Open Graph / Facebook", // Template marker
				"Twitter", // Template marker
				"og:image", // OG tag
			}
			
			for _, expected := range expectations {
				if !bytes.Contains(content, []byte(expected)) {
					t.Errorf("HTML doesn't contain expected content: %s", expected)
				}
			}
		} else {
			t.Errorf("Error reading HTML file: %v", err)
		}
	}
}

// TestHealthEndpoint tests the health endpoint of the API
func TestHealthEndpoint(t *testing.T) {
	// Create a new server instance
	service := NewService()
	
	// Create a request to the health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.HealthHandler(w, r)
	})
	
	// Serve the request
	handler.ServeHTTP(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	
	// Check the response body
	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestAPIIntegration tests the full API flow with a mock request
func TestAPIIntegration(t *testing.T) {
	// Skip this test in CI environments without a browser
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Create a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Create a new service
	service := NewService()
	
	// Initialize the service
	err := service.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}
	defer service.Shutdown()
	
	// Test the service with a sample URL
	// This is a simplified test - adjust based on your actual API
	testURL := "https://example.com"
	result, err := service.GenerateOpenGraph(ctx, testURL)
	
	if err != nil {
		t.Fatalf("Failed to generate Open Graph: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected result but got nil")
	}
	
	// Check that the result contains expected fields
	if result.URL != testURL {
		t.Errorf("Expected URL %s but got %s", testURL, result.URL)
	}
	
	if result.Title == "" {
		t.Error("Expected title but got empty string")
	}
} 