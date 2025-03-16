#!/bin/bash
set -e

# This script builds a container image from a local Alpine Linux tarball
# No pulling from Docker Hub or other registries needed

echo "Building local Alpine-based container for testing..."

# Set script directory and move to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Create volume directory if it doesn't exist
VOLUME_DIR="$HOME/ogdrip-data"
mkdir -p "$VOLUME_DIR"
echo "Using volume directory: $VOLUME_DIR"

# Download Alpine Linux if we don't have it already
ALPINE_VERSION="3.18.0"
ALPINE_TARBALL="$SCRIPT_DIR/alpine-minirootfs-$ALPINE_VERSION-x86_64.tar.gz"
if [ ! -f "$ALPINE_TARBALL" ]; then
    echo "Downloading Alpine Linux $ALPINE_VERSION..."
    curl -L -o "$ALPINE_TARBALL" "https://dl-cdn.alpinelinux.org/alpine/v${ALPINE_VERSION%.*}/releases/x86_64/alpine-minirootfs-$ALPINE_VERSION-x86_64.tar.gz"
fi

# Import Alpine as a local image
echo "Importing Alpine Linux as a local image..."
mkdir -p "$SCRIPT_DIR/alpine-root"
tar -xzf "$ALPINE_TARBALL" -C "$SCRIPT_DIR/alpine-root"

# Package it into a local image
echo "Creating local Alpine image..."
tar -C "$SCRIPT_DIR/alpine-root" -c . | nerdctl load
nerdctl tag imported-image alpine:local

# Build components
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

# Create a Dockerfile that uses our local Alpine image
cat > docker-build/Dockerfile.alpine << 'EOF'
FROM alpine:local

# Install necessary packages
RUN apk add --no-cache nodejs npm curl chromium sqlite supervisor

# Set working directory
WORKDIR /app

# Create directories
RUN mkdir -p /app/frontend /app/backend /app/outputs /etc/supervisor.d

# Copy pre-built files into the image
COPY frontend/dist/ /app/frontend/
COPY backend/ogdrip-backend /app/backend/
COPY docker-build/run.sh /app/run.sh

# Make scripts executable
RUN chmod +x /app/backend/ogdrip-backend /app/run.sh

# Create a minimal package.json for frontend
RUN echo '{"name":"ogdrip-frontend","type":"module","dependencies":{"@sentry/astro":"9.5.0","@sentry/browser":"9.5.0","clsx":"2.1.1","cookie":"0.6.0","mime":"4.0.1"}}' > /app/frontend/package.json

# Install frontend runtime dependencies
RUN cd /app/frontend && npm install --production --no-fund --no-audit

# Expose ports
EXPOSE 3000 8888

# Set entrypoint
CMD ["/app/run.sh"]
EOF

# Create run script
cat > docker-build/run.sh << 'EOF'
#!/bin/sh
echo "Starting OG Drip..."

# Create runtime directories
mkdir -p /app/outputs

# Start the backend
cd /app/backend
export OGDRIP_OUTPUT_PATH=/app/outputs
export OGDRIP_CHROME_PATH=/usr/bin/chromium-browser
export OGDRIP_CHROME_FLAGS="--no-sandbox,--disable-dev-shm-usage,--disable-gpu,--headless,--disable-software-rasterizer"
./ogdrip-backend -service &
backend_pid=$!
echo "Backend started with PID $backend_pid"

# Start the frontend
cd /app/frontend
export HOST=0.0.0.0
export PORT=3000
export NODE_ENV=production
export BACKEND_URL=http://localhost:8888
export PUBLIC_BACKEND_URL=http://localhost:8888
node ./server/entry.mjs &
frontend_pid=$!
echo "Frontend started with PID $frontend_pid"

echo "All services started!"
echo "API: http://localhost:8888/api/health"
echo "Frontend: http://localhost:3000"

# Keep the container running
wait $backend_pid $frontend_pid
EOF
chmod +x docker-build/run.sh

# Build the container
echo "Building container locally..."
nerdctl build --no-cache -t ogdrip-local:alpine -f docker-build/Dockerfile.alpine .

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed. Please check the error messages above."
    exit 1
fi

# Stop any existing container
echo "Stopping any existing test container..."
nerdctl stop ogdrip-test 2>/dev/null || true
nerdctl rm ogdrip-test 2>/dev/null || true

# Run the container
echo "Starting the container..."
nerdctl run --name ogdrip-test \
  -p 3000:3000 \
  -p 8888:8888 \
  -v "$VOLUME_DIR:/app/outputs" \
  -e PUBLIC_BACKEND_URL=http://localhost:8888 \
  -e BACKEND_URL=http://localhost:8888 \
  --restart unless-stopped \
  -d \
  ogdrip-local:alpine

echo ""
echo "Container started!"
echo "Frontend: http://localhost:3000"
echo "Backend API: http://localhost:8888/api/health"
echo ""
echo "To view container logs:"
echo "  nerdctl logs -f ogdrip-test"
echo ""
echo "To stop the container:"
echo "  nerdctl stop ogdrip-test"

# Clean up
rm -rf "$SCRIPT_DIR/alpine-root"
