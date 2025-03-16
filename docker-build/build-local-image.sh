#!/bin/bash
set -e

# This script builds a container image purely from local components
# No pulling from Docker Hub or other registries

echo "Building purely local container image..."

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

# Create a scratch image with our components
echo "Creating local container image from scratch..."

# Create a simple Dockerfile that doesn't require any base image from registries
cat > docker-build/Dockerfile.scratch << 'EOF'
# Start with nothing
FROM scratch

# Create filesystem structure
WORKDIR /

# Add a basic shell (using busybox)
COPY docker-build/busybox /bin/busybox
RUN ["/bin/busybox", "--install", "-s", "/bin"]

# Set up runtime environment
RUN mkdir -p /app/frontend /app/backend /app/outputs /var/log/supervisor

# Copy our application components
COPY frontend/dist/ /app/frontend/
COPY backend/ogdrip-backend /app/backend/
COPY docker-build/run.sh /app/run.sh

# Make scripts executable
RUN chmod +x /app/backend/ogdrip-backend /app/run.sh

# Expose ports
EXPOSE 3000 8888

# Set entrypoint
CMD ["/app/run.sh"]
EOF

# Download busybox as our minimal userland
echo "Downloading busybox for minimal userland..."
if [ ! -f "docker-build/busybox" ]; then
    curl -L -o docker-build/busybox https://busybox.net/downloads/binaries/1.31.0-defconfig-multiarch/busybox-x86_64
    chmod +x docker-build/busybox
fi

# Create a simple run script
cat > docker-build/run.sh << 'EOF'
#!/bin/sh
echo "Starting OG Drip..."

# Create runtime directories
mkdir -p /app/outputs

# Start the backend
cd /app/backend
./ogdrip-backend -service &
backend_pid=$!
echo "Backend started with PID $backend_pid"

# Start the frontend
cd /app/frontend
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

# Build the container using our local Dockerfile
echo "Building container locally..."
nerdctl build --no-cache -t ogdrip-local:scratch -f docker-build/Dockerfile.scratch .

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
  -e OGDRIP_CHROME_PATH=/bin/true \
  -e OGDRIP_OUTPUT_PATH=/app/outputs \
  --restart unless-stopped \
  -d \
  ogdrip-local:scratch

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
