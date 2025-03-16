package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"path/filepath"

	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed openapi.yaml
var openAPISpec embed.FS

// setupSwagger configures the Swagger UI and documentation endpoints
func setupSwagger(mux *http.ServeMux) {
	// Create a handler for serving the OpenAPI YAML file
	mux.HandleFunc("/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		// Read the embedded OpenAPI spec
		data, err := openAPISpec.ReadFile("openapi.yaml")
		if err != nil {
			log.Printf("Error reading OpenAPI spec: %v", err)
			http.Error(w, "Failed to read OpenAPI specification", http.StatusInternalServerError)
			return
		}

		// Set content type and serve the file
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(data)
	})

	// Create a function to serve pre-modified Swagger UI
	// This allows us to point the default URL to our openapi.yaml file
	createSwaggerUIHandler := func() http.Handler {
		return httpSwagger.Handler(
			httpSwagger.URL("/api/openapi.yaml"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
			httpSwagger.DomID("swagger-ui"),
			httpSwagger.PersistAuthorization(true),
		)
	}

	// Register Swagger documentation endpoint
	mux.Handle("/docs/", http.StripPrefix("/docs", createSwaggerUIHandler()))
	mux.Handle("/swagger/", http.StripPrefix("/swagger", createSwaggerUIHandler()))

	// Copy the OpenAPI spec to the output directory for easier access
	// in case embedded files are not working in some environments
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Printf("Warning: Failed to create output directory for OpenAPI spec: %v", err)
		return
	}

	// Extract the openapi.yaml to the outputs directory
	data, err := openAPISpec.ReadFile("openapi.yaml")
	if err != nil {
		log.Printf("Warning: Failed to read embedded OpenAPI spec: %v", err)
		return
	}

	specPath := filepath.Join(config.OutputDir, "openapi.yaml")
	if err := os.WriteFile(specPath, data, 0644); err != nil {
		log.Printf("Warning: Failed to write OpenAPI spec to outputs directory: %v", err)
		return
	}

	log.Printf("OpenAPI documentation available at: %s/docs/", config.BaseURL)
	log.Printf("OpenAPI specification available at: %s/api/openapi.yaml", config.BaseURL)
}
