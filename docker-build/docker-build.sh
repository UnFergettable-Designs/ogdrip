#!/bin/bash
set -e

# This script builds production-ready Docker images for the Open Graph Generator

echo "Building Open Graph Generator Docker Images"
echo "============================================"

# Set script directory and move to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

# Build both components first
echo "Step 1: Building frontend..."
cd frontend
pnpm install
pnpm build
cd ..

echo "Step 2: Building backend..."
cd backend
go mod tidy
go mod download
CGO_ENABLED=1 go build -v -o ogdrip-backend .
cd ..

# Build Docker images using docker-compose
echo "Step 3: Building Docker images..."
docker compose -f docker-compose.production.yml build

echo "Step 4: Creating local test container..."
# Create volume directory if it doesn't exist
VOLUME_DIR="$HOME/ogdrip-data"
mkdir -p "$VOLUME_DIR"

# Stop any existing container
echo "Stopping any existing containers..."
docker compose -f docker-compose.production.yml down

# Run the containers with local configuration
echo "Starting containers..."
docker compose -f docker-compose.production.yml up -d

echo ""
echo "Container started!"
echo "Frontend: http://localhost:3000"
echo "Backend API: http://localhost:8888/api/health"
echo "API Documentation: http://localhost:8888/docs/"
echo ""
echo "To view container logs:"
echo "  docker compose -f docker-compose.production.yml logs -f"
echo ""
echo "To stop the container:"
echo "  docker compose -f docker-compose.production.yml down"
echo ""
echo "Docker images are now ready for deployment to any Docker registry."
