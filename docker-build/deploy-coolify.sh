#!/bin/bash
set -e

echo "Building OG Drip combined image..."
cd "$(dirname "$0")/.."
docker build -t ogdrip-combined:latest -f docker-build/Dockerfile.combined .

echo "Saving image to tar file..."
docker save -o ogdrip-combined.tar ogdrip-combined:latest

echo "Create a deployment package..."
mkdir -p deploy-package
cp docker-build/docker-compose.yml deploy-package/
cp docker-build/supervisord.conf deploy-package/
cp ogdrip-combined.tar deploy-package/

echo "Creating deployment archive..."
tar -czf ogdrip-deploy.tar.gz deploy-package

echo "Deployment package created: ogdrip-deploy.tar.gz"
echo ""
echo "To deploy to Coolify:"
echo "1. Upload ogdrip-deploy.tar.gz to your Coolify server"
echo "2. Extract: tar -xzf ogdrip-deploy.tar.gz"
echo "3. Load image: docker load -i deploy-package/ogdrip-combined.tar"
echo "4. Deploy: cd deploy-package && docker-compose up -d"
echo ""
echo "Done!"
