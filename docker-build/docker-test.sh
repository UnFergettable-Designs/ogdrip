#!/bin/bash
set -e

# Use Docker instead of nerdctl to avoid credential helper issues
# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: docker is not installed or not in your PATH"
    echo "Please install Docker Desktop and try again"
    exit 1
fi

# Set script directory and move to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Create volume directory if it doesn't exist
VOLUME_DIR="$HOME/ogdrip-data"
mkdir -p "$VOLUME_DIR"
echo "Using volume directory: $VOLUME_DIR"

# Build components first
echo "Building frontend..."
cd frontend
pnpm install
pnpm build
cd ..

echo "Building backend..."
cd backend
go mod download
CGO_ENABLED=1 go build -v -o ogdrip-backend .
cd ..

echo "Making sure supervisord.conf exists..."
if [ ! -f "docker-build/supervisord.conf" ]; then
    echo "Error: supervisord.conf not found in docker-build directory"
    exit 1
fi

# Build the container
echo "Building container locally with Docker..."
docker build -t ogdrip-local:test -f docker-build/Dockerfile.test .

# Stop any existing container
echo "Stopping any existing test container..."
docker stop ogdrip-test 2>/dev/null || true
docker rm ogdrip-test 2>/dev/null || true

# Run the container
echo "Starting the container..."
docker run --name ogdrip-test \
  -p 3000:3000 \
  -p 8888:8888 \
  -v "$VOLUME_DIR:/app/outputs" \
  -e PUBLIC_BACKEND_URL=http://localhost:8888 \
  -e BACKEND_URL=http://localhost:8888 \
  -e OGDRIP_CHROME_PATH=/usr/bin/chromium-browser \
  -e OGDRIP_OUTPUT_PATH=/app/outputs \
  --restart unless-stopped \
  -d \
  ogdrip-local:test

echo ""
echo "Container started!"
echo "Frontend: http://localhost:3000"
echo "Backend API: http://localhost:8888/api/health"
echo ""
echo "To view container logs:"
echo "  docker logs -f ogdrip-test"
echo ""
echo "To stop the container:"
echo "  docker stop ogdrip-test"
