#!/bin/bash
set -e

# Make script executable
chmod +x start.sh

# Create necessary directories
mkdir -p outputs
mkdir -p data

# Start the backend service
cd backend
./ogdrip-backend --service

# Start the frontend in production mode
cd ../frontend
NODE_ENV=production HOST=0.0.0.0 PORT=${PORT:-3000} node ./dist/server/entry.mjs
