#!/bin/bash
set -e

# Create necessary directories
mkdir -p outputs
mkdir -p data

# Start the backend server
cd backend && go run . &
BACKEND_PID=$!

# Start the frontend server
cd /app/frontend && pnpm start

# If frontend exits, kill the backend
kill $BACKEND_PID
