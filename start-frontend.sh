#!/bin/bash
cd "$(dirname "$0")"
export HOST=0.0.0.0
export PORT=3000
export NODE_ENV=production
export BACKEND_URL=http://localhost:8888
export PUBLIC_BACKEND_URL=http://localhost:8888
export DEBUG=astro:*
export VERBOSE=true
cd frontend
node ./server/entry.mjs > ../logs/frontend.log 2>&1 &
echo "Frontend started with PID $!"
