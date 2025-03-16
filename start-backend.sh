#!/bin/bash
cd "$(dirname "$0")"

# Set environment variables
export OGDRIP_OUTPUT_PATH="./outputs"
export OGDRIP_CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
export OGDRIP_CHROME_FLAGS="--no-sandbox,--disable-dev-shm-usage,--disable-gpu,--headless,--disable-software-rasterizer"
export OGDRIP_DEBUG="true"
export OGDRIP_LOG_LEVEL="debug"
export PORT="8888"
export ENABLE_CORS="true"
export BASE_URL="http://localhost:8888"

# Start backend in service mode (API server)
./backend/ogdrip-backend -service > ./logs/backend.log 2>&1 &
echo "Backend started with PID $!"
