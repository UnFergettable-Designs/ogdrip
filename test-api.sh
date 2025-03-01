#!/bin/bash

# Test the API health endpoint
echo "Testing API health..."
curl -v http://localhost:8888/api/health

# Test the API generate endpoint with minimal parameters
echo -e "\n\nTesting API generate endpoint..."
curl -v -X POST \
  -F "title=Test Title" \
  -F "description=Test Description" \
  -F "type=website" \
  http://localhost:8888/api/generate 