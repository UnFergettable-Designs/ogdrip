#!/bin/bash
set -e

# Check if nerdctl is installed
if ! command -v nerdctl &> /dev/null; then
    echo "Error: nerdctl is not installed or not in your PATH"
    echo "Please install nerdctl and try again"
    exit 1
fi

# Set script directory and move to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Create volume directory if it doesn't exist
VOLUME_DIR="$HOME/ogdrip-data"
mkdir -p "$VOLUME_DIR"
echo "Using volume directory: $VOLUME_DIR"

# Fix nerdctl credential helper issue
echo "Setting up dummy credential helper..."
# Create a directory for the dummy credential helper if it doesn't exist
mkdir -p "$HOME/bin"

# Create a dummy credential helper script
cat > "$HOME/bin/docker-credential-undefined" << 'EOF'
#!/bin/bash
# This is a dummy credential helper that always returns an empty credential
echo "{}"
EOF

# Make it executable
chmod +x "$HOME/bin/docker-credential-undefined"

# Add it to PATH temporarily if it's not there
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    export PATH="$HOME/bin:$PATH"
    echo "Added $HOME/bin to PATH temporarily"
fi

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

# Show the docker build context
echo "Build context: $(pwd)"
echo "Key files:"
ls -la frontend/dist backend/ogdrip-backend docker-build/supervisord.conf

# Build the container using the simpler test Dockerfile
echo "Building container locally..."
nerdctl build --insecure-registry --progress=plain -t ogdrip-local:test -f docker-build/Dockerfile.test .

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
echo "  nerdctl logs -f ogdrip-test"
echo ""
echo "To stop the container:"
echo "  nerdctl stop ogdrip-test"
