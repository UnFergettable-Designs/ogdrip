#!/bin/bash
set -e

# Script to build and run the ogdrip application with Node.js and pnpm

MODE=${1:-dev}

# Make the script directory for env.ts if it doesn't exist
mkdir -p frontend/src/utils
mkdir -p frontend/public/images

# Always make sure env.ts exists
if [ ! -f frontend/src/utils/env.ts ]; then
  cat > frontend/src/utils/env.ts << 'EOL'
/**
 * Environment detection utility for Svelte and Astro
 * Compatible with Node.js environments
 * Exports required constants for Svelte 5 compatibility
 */

// =========== For Svelte 5 compatibility ===========
// Constants needed by Svelte 5 internals
export const DEV: boolean = process.env.NODE_ENV !== 'production';
export const BROWSER: boolean = typeof window !== 'undefined';
export const DEBUG: boolean = process.env.DEBUG === 'true';

// Additional constants that might be needed by Svelte 5
export const BUILD: {
  SSR: boolean;
} = {
  SSR: typeof window === 'undefined'
};

// Needed for server-side rendering
export const NODE_ENV: string = process.env.NODE_ENV || 'development';
// ==================================================

// Standard environment exports
export const browser: boolean = typeof window !== 'undefined';
export const dev: boolean = process.env.NODE_ENV !== 'production';
export const server: boolean = !browser;
EOL
  echo "Created env.ts file"
fi

# Always make sure placeholder.svg exists
if [ ! -f frontend/public/images/placeholder.svg ]; then
  cat > frontend/public/images/placeholder.svg << 'EOL'
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg width="50" height="50" viewBox="0 0 50 50" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="50" height="50" rx="10" fill="#4F46E5"/>
  <path d="M15 25C15 19.4772 19.4772 15 25 15C30.5228 15 35 19.4772 35 25C35 30.5228 30.5228 35 25 35C19.4772 35 15 30.5228 15 25Z" fill="white"/>
  <path d="M20 20L30 30M30 20L20 30" stroke="#4F46E5" stroke-width="3" stroke-linecap="round"/>
</svg>
EOL
  echo "Created placeholder logo"
fi

# Install pnpm if not available
if ! command -v pnpm &> /dev/null
then
    echo "pnpm not found, attempting to install..."
    if command -v npm &> /dev/null
    then
        npm install -g pnpm
    else
        echo "npm not found, please install npm and pnpm manually."
        exit 1
    fi
fi

case $MODE in
  install)
    echo "Installing dependencies..."
    cd frontend && pnpm install
    ;;
  dev)
    echo "Starting in development mode..."
    cd frontend && pnpm run dev
    ;;
  build)
    echo "Building the application..."
    cd frontend && pnpm run build
    ;;
  preview)
    echo "Building and previewing the application..."
    cd frontend && pnpm run build && pnpm run preview
    ;;
  docker)
    echo "Building Docker containers..."
    docker-compose build
    echo "Starting Docker containers..."
    docker-compose up
    ;;
  docker:build)
    echo "Building Docker containers..."
    docker-compose build
    ;;
  docker:clean)
    echo "Rebuilding Docker containers from scratch..."
    docker-compose down -v
    docker-compose build --no-cache
    echo "Starting Docker containers..."
    docker-compose up
    ;;
  clean)
    echo "Cleaning up node_modules directories..."
    rm -rf frontend/node_modules
    rm -rf node_modules
    echo "Removing dist directories..."
    rm -rf frontend/dist
    ;;
  *)
    echo "Unknown mode: $MODE"
    echo "Usage: ./run.sh [install|dev|build|preview|docker|docker:build|docker:clean|clean]"
    exit 1
    ;;
esac 