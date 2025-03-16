#!/bin/bash
set -e

# Set current directory to project root
cd "$(dirname "$0")/.."

echo "Building OG Drip components locally for testing..."

# Build frontend
echo "Building frontend..."
cd frontend
pnpm install
pnpm build
cd ..

# Build backend
echo "Building backend..."
cd backend
go mod download
CGO_ENABLED=1 go build -v -o ogdrip-backend .
cd ..

# Create a test directory
echo "Creating test directory..."
TEST_DIR="./local-test"
mkdir -p $TEST_DIR/{frontend,backend,outputs,logs}

# Copy files to test directory
echo "Copying files to test directory..."
cp -r frontend/dist/* $TEST_DIR/frontend/
cp backend/ogdrip-backend $TEST_DIR/backend/
chmod +x $TEST_DIR/backend/ogdrip-backend

# Create minimal package.json for frontend
echo "Creating minimal package.json..."
echo '{"name":"ogdrip-frontend","type":"module","dependencies":{"@sentry/astro":"9.5.0","@sentry/browser":"9.5.0","clsx":"2.1.1","cookie":"0.6.0","mime":"4.0.1"}}' > $TEST_DIR/frontend/package.json

# Install frontend runtime dependencies
echo "Installing frontend runtime dependencies..."
cd $TEST_DIR/frontend
npm install --production --no-fund --no-audit
cd ../..

# Create startup scripts
echo "Creating startup scripts..."
cat > $TEST_DIR/start-backend.sh << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"

# Set environment variables
export OGDRIP_OUTPUT_PATH="./outputs"
export OGDRIP_CHROME_PATH="/usr/bin/google-chrome"
export OGDRIP_CHROME_FLAGS="--no-sandbox,--disable-dev-shm-usage,--disable-gpu,--headless,--disable-software-rasterizer"

# Start backend in service mode (API server)
./backend/ogdrip-backend -service > ./logs/backend.log 2>&1 &
echo "Backend started with PID $!"
EOF

cat > $TEST_DIR/start-frontend.sh << 'EOF'
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

chmod +x $TEST_DIR/start-backend.sh $TEST_DIR/start-frontend.sh

# Create a combined startup script
cat > $TEST_DIR/start-all.sh << 'EOF'
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

chmod +x $TEST_DIR/start-all.sh

# Create a stop script
cat > $TEST_DIR/stop-all.sh << 'EOF'
#!/bin/bash
echo "Stopping OG Drip services..."
pkill -f "ogdrip-backend" || true
pkill -f "node ./server/entry.mjs" || true
echo "All services stopped."
EOF

chmod +x $TEST_DIR/stop-all.sh

echo "Local test environment created in $TEST_DIR"
echo "To start the services:"
echo "  cd $TEST_DIR"
echo "  ./start-all.sh"
echo ""
echo "To stop the services:"
echo "  ./stop-all.sh"
echo ""
echo "Note: You'll need Chrome/Chromium installed for the backend to work properly."
echo "If you have issues, check the log files in $TEST_DIR/logs/"
