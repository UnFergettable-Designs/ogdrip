#!/bin/bash
set -e

echo "=== Building frontend locally with SSR ==="
# Change to the repo root directory
cd "$(dirname "$0")"

# Install dependencies at the monorepo level
echo "Installing dependencies..."
pnpm install

# Build the frontend
echo "Building frontend..."
cd frontend
NODE_ENV=production pnpm build
cd ..

echo "=== Building the ultra-lightweight frontend container ==="
# Build the Docker image using the ultra-lightweight approach
nerdctl build -f Dockerfile.frontend.production -t ogdrip-frontend:production .

echo "=== Stopping any existing containers ==="
nerdctl rm -f ogdrip-frontend-production-test 2>/dev/null || true

echo "=== Running frontend container for testing ==="
nerdctl run --name ogdrip-frontend-production-test -p 3000:3000 -d ogdrip-frontend:production

echo "=== Testing frontend availability ==="
echo "Waiting for container to start..."
sleep 3
echo -n "Testing connection to http://localhost:3000... "
if curl -s --head http://localhost:3000 | grep "200 OK" > /dev/null; then
  echo "SUCCESS! Frontend is running properly."
else
  echo "WARNING: Could not connect to frontend. Check container logs:"
  nerdctl logs ogdrip-frontend-production-test
fi

echo ""
echo "=== Container started! ==="
echo "You can access the frontend at: http://localhost:3000"
echo ""
echo "To view logs:"
echo "nerdctl logs -f ogdrip-frontend-production-test"
echo ""
echo "To stop the container:"
echo "nerdctl stop ogdrip-frontend-production-test"
echo ""
echo "To remove the container:"
echo "nerdctl rm ogdrip-frontend-production-test"
echo ""
echo "For Coolify deployment, use the deploy-coolify.sh script."
