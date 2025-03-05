#!/bin/bash

# Set API service environment variable
export OG_API_SERVICE=true

# Set port (default 8888)
export PORT=8888

# Base URL for responses
export BASE_URL=http://localhost:8888

# Log level (debug, info, warn, error)
export LOG_LEVEL=debug

# Build the application with all .go files
go build -o og-generator *.go

# Run the service
./og-generator -service 