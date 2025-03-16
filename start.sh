#!/bin/bash
set -e

# Create necessary directories if they don't exist
mkdir -p /app/outputs
mkdir -p /app/data

# Start the backend with --service flag
cd backend
./ogdrip-backend --service &

# Start the frontend in production mode
cd ../frontend
NODE_ENV=production HOST=0.0.0.0 PORT=${PORT:-3000} node ./dist/server/entry.mjs
