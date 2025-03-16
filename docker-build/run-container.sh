#!/bin/bash
set -e

# Check if nerdctl is installed
if ! command -v nerdctl &> /dev/null; then
    echo "Error: nerdctl is not installed or not in your PATH"
    echo "Please install nerdctl and try again"
    exit 1
fi

# Create volume directory if it doesn't exist
VOLUME_DIR="$HOME/ogdrip-data"
mkdir -p "$VOLUME_DIR"
echo "Using volume directory: $VOLUME_DIR"

# Pull the latest image
echo "Pulling the latest image from GitHub Container Registry..."
nerdctl pull ghcr.io/unfergettable-designs/ogdrip:latest

# Stop any existing container
echo "Stopping any existing container..."
nerdctl stop ogdrip 2>/dev/null || true
nerdctl rm ogdrip 2>/dev/null || true

# Run the container
echo "Starting the container..."
nerdctl run --name ogdrip \
  -p 3000:3000 \
  -p 8888:8888 \
  -v "$VOLUME_DIR:/app/outputs" \
  -e PUBLIC_BACKEND_URL=http://localhost:8888 \
  -e BACKEND_URL=http://localhost:8888 \
  --health-cmd "curl -f http://localhost:8888/api/health || exit 1" \
  --health-interval 30s \
  --health-timeout 10s \
  --health-retries 3 \
  --health-start-period 40s \
  --restart unless-stopped \
  -d \
  ghcr.io/unfergettable-designs/ogdrip:latest

echo ""
echo "Container started!"
echo "Frontend: http://localhost:3000"
echo "Backend API: http://localhost:8888/api/health"
echo ""
echo "To view container logs:"
echo "  nerdctl logs -f ogdrip"
echo ""
echo "To stop the container:"
echo "  nerdctl stop ogdrip"
