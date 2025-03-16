#!/bin/bash
set -e

echo "Cleaning up Open Graph Generator project..."

# Stop any running test services
if [ -d "$HOME/ogdrip-test" ]; then
  echo "Stopping test services..."
  cd "$HOME/ogdrip-test" && ./stop-all.sh
fi

# Remove temporary test directories
echo "Removing temporary test directories..."
rm -rf "$HOME/ogdrip-test"
rm -rf "$HOME/ogdrip-data"
rm -rf "$HOME/ogdrip-output"

# Clean up docker build artifacts
echo "Cleaning Docker build artifacts..."
rm -rf docker-build/alpine-root
rm -f docker-build/alpine.tar
rm -f docker-build/alpine-minirootfs-*.tar.gz

# Clean up build artifacts
echo "Cleaning build artifacts..."
cd backend && rm -f ogdrip-backend og-generator
cd ../frontend && rm -rf dist dev-dist .astro

# Clean up node modules (optional, comment out if you want to keep them)
# echo "Cleaning node_modules (this might take a while)..."
# rm -rf node_modules
# rm -rf frontend/node_modules
# rm -rf shared/node_modules

echo "Cleanup complete! Your project is now ready for GitHub."
echo "Use ./docker-build/docker-build.sh to test Docker builds."
