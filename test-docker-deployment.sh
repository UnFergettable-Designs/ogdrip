#!/bin/bash
# Test script for Docker deployment

set -e

echo "🐳 Testing OG Drip Docker Deployment"
echo "===================================="

# Test 1: Check if required files exist
echo "✅ Checking required files..."
required_files=(
    "Dockerfile"
    "Dockerfile.production" 
    ".dockerignore"
    "docker-compose.yml"
    "DOCKER_DEPLOYMENT.md"
    "start.sh"
    "package.json"
    "pnpm-workspace.yaml"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✓ $file exists"
    else
        echo "  ✗ $file missing"
        exit 1
    fi
done

# Test 2: Validate Dockerfile structure
echo ""
echo "✅ Validating Dockerfile structure..."

if grep -q "FROM node:22-alpine AS frontend-builder" Dockerfile; then
    echo "  ✓ Frontend builder stage defined"
else
    echo "  ✗ Frontend builder stage missing"
    exit 1
fi

if grep -q "FROM golang:1.22-alpine AS backend-builder" Dockerfile; then
    echo "  ✓ Backend builder stage defined"
else
    echo "  ✗ Backend builder stage missing"
    exit 1
fi

if grep -q "FROM node:22-alpine AS runtime" Dockerfile; then
    echo "  ✓ Runtime stage defined"
else
    echo "  ✗ Runtime stage missing"
    exit 1
fi

# Test 3: Check if build artifacts are copied correctly
echo ""
echo "✅ Checking build artifact copying..."

if grep -q "COPY --from=frontend-builder /app/frontend/dist ./frontend/dist" Dockerfile; then
    echo "  ✓ Frontend build artifacts copied"
else
    echo "  ✗ Frontend build artifacts not copied correctly"
    exit 1
fi

if grep -q "COPY --from=backend-builder /app/backend/build/ogdrip-backend ./backend/build/ogdrip-backend" Dockerfile; then
    echo "  ✓ Backend binary copied"
else
    echo "  ✗ Backend binary not copied correctly"
    exit 1
fi

# Test 4: Verify startup script is copied and made executable
echo ""
echo "✅ Checking startup script setup..."

if grep -q "COPY start.sh ./" Dockerfile && grep -q "RUN chmod +x start.sh" Dockerfile; then
    echo "  ✓ Startup script copied and made executable"
else
    echo "  ✗ Startup script not properly configured"
    exit 1
fi

# Test 5: Check port exposure
echo ""
echo "✅ Checking port configuration..."

if grep -q "EXPOSE 8888 3000" Dockerfile; then
    echo "  ✓ Ports 8888 and 3000 exposed"
else
    echo "  ✗ Required ports not exposed"
    exit 1
fi

# Test 6: Validate docker-compose configuration  
echo ""
echo "✅ Validating docker-compose.yml..."

if grep -q "dockerfile: Dockerfile.production" docker-compose.yml; then
    echo "  ✓ Production Dockerfile referenced"
else
    echo "  ✗ Production Dockerfile not referenced"
    exit 1
fi

if grep -q "8888:8888" docker-compose.yml && grep -q "3000:3000" docker-compose.yml; then
    echo "  ✓ Port mappings configured"
else
    echo "  ✗ Port mappings missing"
    exit 1
fi

# Test 7: Check environment variables
echo ""
echo "✅ Checking environment variables..."

required_env_vars=("NODE_ENV" "BASE_URL" "CHROME_PATH" "GO111MODULE")
for env_var in "${required_env_vars[@]}"; do
    if grep -q "$env_var" docker-compose.yml; then
        echo "  ✓ $env_var configured"
    else
        echo "  ✗ $env_var missing"
        exit 1
    fi
done

echo ""
echo "🎉 All Docker deployment tests passed!"
echo ""
echo "Next steps:"
echo "  1. Build the Docker image: docker build -t ogdrip ."
echo "  2. Or use production version: docker build -f Dockerfile.production -t ogdrip:prod ."
echo "  3. Run with docker-compose: docker-compose up"
echo "  4. Access application at http://localhost:8888 (backend) and http://localhost:3000 (frontend)"