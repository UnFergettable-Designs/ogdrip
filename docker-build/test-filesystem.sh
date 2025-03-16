#!/bin/bash
set -e

# This script runs the components directly in the file system
# No containers or images needed - just pure native execution

echo "Setting up native test environment..."

# Set script directory and move to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Create output directory if it doesn't exist
OUTPUT_DIR="$HOME/ogdrip-output"
mkdir -p "$OUTPUT_DIR"
echo "Using output directory: $OUTPUT_DIR"

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

# Create a test directory to run from
TEST_DIR="$HOME/ogdrip-test"
mkdir -p "$TEST_DIR"
echo "Created test directory: $TEST_DIR"

# Set up directories
mkdir -p "$TEST_DIR/frontend"
mkdir -p "$TEST_DIR/backend"
mkdir -p "$TEST_DIR/outputs"
mkdir -p "$TEST_DIR/logs"

# Copy files
echo "Copying files to test directory..."
cp -r frontend/dist/* "$TEST_DIR/frontend/"
cp backend/ogdrip-backend "$TEST_DIR/backend/"
chmod +x "$TEST_DIR/backend/ogdrip-backend"

# Create minimal package.json for frontend
echo '{"name":"ogdrip-frontend","type":"module","dependencies":{"@sentry/astro":"9.5.0","@sentry/browser":"9.5.0","clsx":"2.1.1","cookie":"0.6.0","mime":"4.0.1"}}' > "$TEST_DIR/frontend/package.json"

# Install frontend dependencies
cd "$TEST_DIR/frontend"
npm install --production --no-fund --no-audit
cd -

# Create start scripts
cat > "$TEST_DIR/start-backend.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"

export OGDRIP_OUTPUT_PATH="./outputs"
# Use appropriate Chrome path for your system
if [ -f "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" ]; then
  export OGDRIP_CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
elif [ -f "/usr/bin/google-chrome" ]; then
  export OGDRIP_CHROME_PATH="/usr/bin/google-chrome"
elif [ -f "/usr/bin/chromium-browser" ]; then
  export OGDRIP_CHROME_PATH="/usr/bin/chromium-browser"
else
  echo "No Chrome/Chromium found - using a dummy path"
  export OGDRIP_CHROME_PATH="/bin/echo"
fi

export OGDRIP_CHROME_FLAGS="--no-sandbox,--disable-dev-shm-usage,--disable-gpu,--headless,--disable-software-rasterizer"
export PORT="8888"
export BASE_URL="http://localhost:8888"
export ENABLE_CORS="true"

echo "Starting backend with Chrome at: $OGDRIP_CHROME_PATH"
./backend/ogdrip-backend -service > ./logs/backend.log 2>&1 &
echo "Backend started with PID $!"
EOF

cat > "$TEST_DIR/start-frontend.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"

export HOST=0.0.0.0
export PORT=3000
export NODE_ENV=production
export BACKEND_URL=http://localhost:8888
export PUBLIC_BACKEND_URL=http://localhost:8888

cd frontend
node ./server/entry.mjs > ../logs/frontend.log 2>&1 &
echo "Frontend started with PID $!"
EOF

cat > "$TEST_DIR/start-all.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"
echo "Starting backend..."
./start-backend.sh
echo "Starting frontend..."
./start-frontend.sh
echo "All services started!"
echo "Frontend: http://localhost:3000"
echo "Backend API: http://localhost:8888/api/health"
echo "Log files are in ./logs/"
EOF

cat > "$TEST_DIR/stop-all.sh" << 'EOF'
#!/bin/bash
echo "Stopping OG Drip services..."
pkill -f "ogdrip-backend" || true
pkill -f "node ./server/entry.mjs" || true
echo "All services stopped."
EOF

chmod +x "$TEST_DIR/start-backend.sh" "$TEST_DIR/start-frontend.sh" "$TEST_DIR/start-all.sh" "$TEST_DIR/stop-all.sh"

echo ""
echo "Test environment created at: $TEST_DIR"
echo ""
echo "To start the services:"
echo "  cd $TEST_DIR"
echo "  ./start-all.sh"
echo ""
echo "To check logs:"
echo "  tail -f $TEST_DIR/logs/backend.log"
echo "  tail -f $TEST_DIR/logs/frontend.log"
echo ""
echo "To stop the services:"
echo "  $TEST_DIR/stop-all.sh"
echo ""
echo "Starting services now..."

# Start services
cd "$TEST_DIR"
./start-all.sh
