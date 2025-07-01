#!/bin/bash
set -e

# Set up signal handling for graceful shutdown
cleanup() {
    echo "Shutting down services..."
    if [ ! -z "$BACKEND_PID" ] && kill -0 $BACKEND_PID 2>/dev/null; then
        echo "Stopping backend (PID: $BACKEND_PID)"
        kill -TERM $BACKEND_PID
        wait $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ] && kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "Stopping frontend (PID: $FRONTEND_PID)"
        kill -TERM $FRONTEND_PID
        wait $FRONTEND_PID 2>/dev/null || true
    fi
    exit 0
}

# Trap signals for graceful shutdown
trap cleanup SIGTERM SIGINT

# Create necessary directories
mkdir -p outputs
mkdir -p data

# Set environment variables for production
export NODE_ENV=production
export GO111MODULE=on

# Start the backend server using compiled binary
echo "Starting backend server..."
if [ -f "backend/build/ogdrip-backend" ]; then
    cd backend && ./build/ogdrip-backend -service &
    BACKEND_PID=$!
    echo "Backend started with PID: $BACKEND_PID"
else
    echo "ERROR: Backend binary not found at backend/build/ogdrip-backend"
    echo "Expected binary should be built by nixpacks build phase"
    exit 1
fi

# Check if frontend exists and has package.json
if [ -f "frontend/package.json" ]; then
    echo "Starting frontend server..."
    cd frontend && pnpm start &
    FRONTEND_PID=$!
    echo "Frontend started with PID: $FRONTEND_PID"
    
    # Wait for frontend to exit
    wait $FRONTEND_PID
else
    echo "Frontend not found or not built. Running backend only."
    # Wait for backend to exit if no frontend
    wait $BACKEND_PID
fi

# Clean up
cleanup
