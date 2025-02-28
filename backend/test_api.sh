#!/bin/bash

# Test script for Open Graph Generator API
# Usage: ./test_api.sh [URL]

# Default URL if none provided
URL=${1:-"https://example.com"}

echo "Testing Open Graph Generator API..."
echo "Using URL: $URL"
echo

# Health check
echo "1. Testing health endpoint..."
curl -s http://localhost:8888/api/health | jq .
echo

# Generate Open Graph assets
echo "2. Testing generate endpoint..."
RESPONSE=$(curl -s -X POST \
  http://localhost:8888/api/generate \
  -F "url=$URL" \
  -F "title=API Test" \
  -F "description=Testing the Open Graph Generator API")

echo $RESPONSE | jq .
echo

# Extract image URL if available
IMAGE_URL=$(echo $RESPONSE | jq -r .image_url)
if [ "$IMAGE_URL" != "null" ] && [ "$IMAGE_URL" != "" ]; then
  echo "3. Opening generated image in browser..."
  open $IMAGE_URL 2>/dev/null || xdg-open $IMAGE_URL 2>/dev/null || echo "Could not open browser automatically. Image URL: $IMAGE_URL"
fi

# Extract meta tags URL if available
META_URL=$(echo $RESPONSE | jq -r .meta_tags_url)
if [ "$META_URL" != "null" ] && [ "$META_URL" != "" ]; then
  echo "4. Opening meta tags HTML in browser..."
  open $META_URL 2>/dev/null || xdg-open $META_URL 2>/dev/null || echo "Could not open browser automatically. Meta tags URL: $META_URL"
fi

echo
echo "Test completed!" 