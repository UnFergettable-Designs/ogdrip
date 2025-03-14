#!/bin/bash
set -e

# Make sure we're running in the backend directory
cd "$(dirname "$0")"
echo "Running in $(pwd)"

echo "Cleaning up any existing test containers..."
nerdctl rm -f ogdrip-backend-test 2>/dev/null || true

echo "Building backend image..."
nerdctl build -f Dockerfile.production -t ogdrip-backend:test .

echo "Running backend container..."
nerdctl run --name ogdrip-backend-test -p 8888:8888 \
  -e PORT=8888 \
  -e BASE_URL=http://localhost:8888 \
  -e CHROME_PATH=/usr/bin/chromium \
  -e OUTPUT_DIR=/app/outputs \
  -e ENABLE_CORS=true \
  -d ogdrip-backend:test

echo "Waiting for backend to start..."
sleep 3

echo "Testing health endpoint..."
curl -v http://localhost:8888/api/health || echo "Failed to connect to health endpoint"

echo "Container logs:"
nerdctl logs ogdrip-backend-test

# Keep the container running for further tests
echo "Backend container is running. Press Ctrl+C to stop it."
echo "Run 'nerdctl logs -f ogdrip-backend-test' to view logs"
echo "Run 'nerdctl rm -f ogdrip-backend-test' to stop and remove the container"
