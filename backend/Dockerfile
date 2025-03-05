FROM golang:1.23-alpine

# Note: This base image will be pulled from Docker Hub.
# If you encounter authentication issues, you may need to:
# 1. Login to Docker Hub: docker login
# 2. Or use a local development approach without Docker (see LOCAL_DEPLOYMENT.md)

# Install required dependencies
RUN apk add --no-cache chromium ca-certificates git gcc musl-dev

# Set environment variables
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_FLAGS=--headless,--disable-gpu,--no-sandbox
ENV CGO_ENABLED=1

# Set the working directory
WORKDIR /app

# Create necessary files for minimal build
COPY <<EOF /app/go.mod
module og-generator

go 1.23
EOF

COPY <<EOF /app/main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8888"
    }
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Open Graph Generator API Running")
    })
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "{\"success\":true,\"message\":\"Open Graph Generator API is running\"}")
    })
    
    log.Printf("Starting server on port %s", port)
    if err := http.ListenAndServe(":" + port, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
EOF

# Build the executable
RUN go build -o og-generator .

# Create output and data directories
RUN mkdir -p /app/outputs
RUN mkdir -p /app/data

# Expose port
EXPOSE 8888

# Command that will be overridden at runtime
CMD ["./og-generator"] 