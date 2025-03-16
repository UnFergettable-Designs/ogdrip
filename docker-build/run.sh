#!/bin/sh
echo "Starting OG Drip..."

# Create runtime directories
mkdir -p /app/outputs

# Start the backend
cd /app/backend
./ogdrip-backend -service &
backend_pid=$!
echo "Backend started with PID $backend_pid"

# Start the frontend
cd /app/frontend
node ./server/entry.mjs &
frontend_pid=$!
echo "Frontend started with PID $frontend_pid"

echo "All services started!"
echo "API: http://localhost:8888/api/health"
echo "Frontend: http://localhost:3000"

# Keep the container running
wait $backend_pid $frontend_pid
