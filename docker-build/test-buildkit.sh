#!/bin/bash
set -e

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

# Create a temporary Dockerfile that doesn't need to pull from a registry
TMPDIR=$(mktemp -d)
cat > "$TMPDIR/Dockerfile" << EOF
# Use buildkit's scratch image
FROM scratch

# Create app structure
WORKDIR /app

# Copy pre-built files
COPY ./frontend/dist /app/frontend
COPY ./backend/ogdrip-backend /app/backend
COPY ./docker-build/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Create directories
WORKDIR /app
RUN mkdir -p /app/outputs

# Make the backend executable
RUN chmod +x /app/backend/ogdrip-backend

# Start script
COPY ./docker-build/start.sh /app/start.sh
RUN chmod +x /app/start.sh

# Command
CMD ["/app/start.sh"]
EOF

# Create a start script for the container
cat > "$SCRIPT_DIR/start.sh" << EOF
#!/bin/sh
echo "Starting OG Drip..."
echo "Backend API: http://localhost:8888/api/health"
echo "Frontend: http://localhost:3000"

# Start the backend
cd /app/backend
./ogdrip-backend -service &

# Start the frontend
cd /app/frontend
node ./server/entry.mjs
EOF
chmod +x "$SCRIPT_DIR/start.sh"

echo "Building container using tar archive method..."
# Create a tarball of the filesystem
echo "Creating tarball of the filesystem..."
cd "$TMPDIR" && tar -cf ../context.tar .

# Build using tar context
echo "Building container using tar context..."
docker buildx build --progress=plain -t ogdrip-local:test -f "$TMPDIR/Dockerfile" "$TMPDIR"

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
