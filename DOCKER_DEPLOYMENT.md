# Docker Deployment Guide

This guide explains how to deploy the OG Drip application using Docker with the multi-stage Dockerfile approach.

## Overview

The multi-stage Dockerfile enables deployment of the entire application stack (Go backend and Astro frontend) as a single unit. This approach simplifies the deployment process by building both services from the main `ogdrip` repository and running them in a single container.

## Dockerfile Structure

The `Dockerfile` is structured as follows:

1. **Frontend Builder Stage**: Builds the Astro/Svelte frontend and shared TypeScript types
2. **Backend Builder Stage**: Builds the Go backend binary
3. **Final Runtime Stage**: Combines the build artifacts into a minimal final image

### Available Dockerfiles

- `Dockerfile`: Basic version with minimal system dependencies
- `Dockerfile.production`: Full production version with Chromium and font packages

## Build Instructions

### Basic Build

```bash
# Build the basic version (minimal dependencies)
docker build -t ogdrip:latest .

# Build the production version (includes Chromium and fonts)
docker build -f Dockerfile.production -t ogdrip:production .
```

### Build Arguments

The Dockerfile supports building in environments with limited network access. The basic `Dockerfile` avoids external package installations that might fail in restricted environments.

## Running the Container

### Basic Run

```bash
# Run the application
docker run -p 8888:8888 -p 3000:3000 ogdrip:latest

# Run with environment variables
docker run -p 8888:8888 -p 3000:3000 \
  -e BASE_URL=http://localhost:8888 \
  -e CHROME_PATH=/usr/bin/chromium \
  ogdrip:latest
```

### Production Run

```bash
# Run the production version with all features
docker run -p 8888:8888 -p 3000:3000 \
  -e BASE_URL=http://your-domain.com \
  -e CHROME_PATH=/usr/bin/chromium \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/outputs:/app/outputs \
  ogdrip:production
```

### Docker Compose (Recommended)

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  ogdrip:
    build:
      context: .
      dockerfile: Dockerfile.production
    ports:
      - "8888:8888"
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - BASE_URL=http://localhost:8888
      - CHROME_PATH=/usr/bin/chromium
      - GO111MODULE=on
    volumes:
      - ./data:/app/data
      - ./outputs:/app/outputs
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NODE_ENV` | Node.js environment | `production` |
| `GO111MODULE` | Go modules setting | `on` |
| `BASE_URL` | Base URL for the application | `http://localhost:8888` |
| `CHROME_PATH` | Path to Chromium binary | `/usr/bin/chromium` |
| `CHROME_FLAGS` | Chromium launch flags | Pre-configured for headless mode |
| `PORT` | Frontend port | `3000` |

## Volumes

- `/app/data`: Database and persistent data storage
- `/app/outputs`: Generated files and outputs

## Health Checks

The application provides a health check endpoint at `/api/health` on port 8888.

Add to your Docker Compose:

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8888/api/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

## Troubleshooting

### Common Issues

1. **Chromium not found**: Use `Dockerfile.production` which includes Chromium installation
2. **Network timeouts during build**: Use the basic `Dockerfile` for environments with limited network access
3. **Permission issues**: Ensure the container has proper permissions for the mounted volumes

### Debugging

```bash
# Run with shell access
docker run -it --entrypoint /bin/sh ogdrip:latest

# View logs
docker logs <container-id>

# Check running processes
docker exec <container-id> ps aux
```

## Performance Considerations

- The frontend serves static files through Astro's SSR mode
- The backend is a compiled Go binary for optimal performance
- Use multi-stage builds to minimize final image size
- Consider using `.dockerignore` to exclude unnecessary files

## Security Notes

- The container runs as non-root where possible
- Chromium runs in sandboxed mode with security flags
- Use environment variables for sensitive configuration
- Consider using Docker secrets for production deployments

## Building for Different Architectures

```bash
# Build for multiple architectures
docker buildx build --platform linux/amd64,linux/arm64 -t ogdrip:multi-arch .
```

## Integration with Existing Deployment

This Docker approach supports the "Option 1" deployment strategy mentioned in the codebase and is compatible with:

- Nixpacks (alternative build system)
- Coolify deployment platform
- Standard container orchestration systems
- Cloud container services (AWS ECS, Google Cloud Run, etc.)