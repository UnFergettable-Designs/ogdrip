#!/bin/bash
set -e

# Ensure directories exist
mkdir -p outputs
mkdir -p data

# Start the backend server
cd backend && ./build/ogdrip-service -service &
BACKEND_PID=$!

# Start the frontend server in production mode
cd /app/frontend && pnpm start:prod

# If frontend exits, kill the backend
kill $BACKEND_PID
