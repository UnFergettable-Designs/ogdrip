# Docker Deployment Guide

This guide covers deploying OG Drip using Docker for both development and production environments.

## Overview

OG Drip supports multiple Docker deployment strategies:

1. **Single Container** - Frontend and backend in one container (recommended for production)
2. **Multi-Container** - Separate containers with Docker Compose (recommended for development)
3. **Production Container** - Optimized build for production deployment

## Prerequisites

- **Docker** >= 20.10.0
- **Docker Compose** >= 2.0.0 (for multi-container setup)
- **4GB RAM** minimum, 8GB recommended
- **10GB disk space** for images and data

## Quick Start

### Single Container Deployment

```bash
# Clone the repository
git clone https://github.com/yourusername/ogdrip.git
cd ogdrip

# Build and run with Docker Compose
docker-compose up -d

# Access the application
open http://localhost:3000
```

## Docker Images

### Development Image

The development image includes both frontend and backend with hot reloading:

```dockerfile
# Dockerfile
FROM node:22-alpine AS frontend-base
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install

FROM golang:1.24-alpine AS backend-base
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

FROM frontend-base AS frontend-build
COPY frontend/ .
RUN pnpm build

FROM backend-base AS backend-build
COPY backend/ .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates chromium
WORKDIR /app

COPY --from=frontend-build /app/frontend/dist ./frontend/dist
COPY --from=backend-build /app/backend/main ./backend/

EXPOSE 3000 8888
CMD ["./backend/main"]
```

### Production Image

The production image is optimized for size and performance:

```dockerfile
# Dockerfile.production
FROM node:22-alpine AS frontend-deps
WORKDIR /app
COPY frontend/package.json frontend/pnpm-lock.yaml ./frontend/
COPY shared/package.json ./shared/
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

FROM frontend-deps AS frontend-build
COPY frontend/ ./frontend/
COPY shared/ ./shared/
RUN cd frontend && pnpm build

FROM golang:1.24-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest AS production
RUN apk --no-cache add \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S ogdrip && \
    adduser -S ogdrip -u 1001

# Copy application files
COPY --from=frontend-build --chown=ogdrip:ogdrip /app/frontend/dist ./frontend/dist
COPY --from=backend-build --chown=ogdrip:ogdrip /app/backend/main ./backend/
COPY --chown=ogdrip:ogdrip nginx.conf ./

# Create necessary directories
RUN mkdir -p ./backend/data ./backend/outputs && \
    chown -R ogdrip:ogdrip ./backend/

USER ogdrip

EXPOSE 8888
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8888/api/health || exit 1

CMD ["./backend/main"]
```

## Docker Compose Configurations

### Development Configuration

```yaml
# docker-compose.yml
version: '3.8'

services:
  frontend:
    build:
      context: .
      dockerfile: Dockerfile.dev
      target: frontend
    ports:
      - '3000:3000'
    volumes:
      - ./frontend:/app/frontend
      - /app/frontend/node_modules
    environment:
      - PUBLIC_BACKEND_URL=http://localhost:8888
      - BACKEND_URL=http://backend:8888
    depends_on:
      - backend
    networks:
      - ogdrip-network

  backend:
    build:
      context: .
      dockerfile: Dockerfile.dev
      target: backend
    ports:
      - '8888:8888'
    volumes:
      - ./backend:/app/backend
      - ./backend/data:/app/backend/data
      - ./backend/outputs:/app/backend/outputs
    environment:
      - PORT=8888
      - HOST=0.0.0.0
      - DATABASE_PATH=./data/ogdrip.db
      - ADMIN_TOKEN=${ADMIN_TOKEN:-dev_admin_token}
    networks:
      - ogdrip-network

networks:
  ogdrip-network:
    driver: bridge

volumes:
  backend-data:
  backend-outputs:
```

### Production Configuration

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  ogdrip:
    build:
      context: .
      dockerfile: Dockerfile.production
    ports:
      - '8888:8888'
    volumes:
      - backend-data:/app/backend/data
      - backend-outputs:/app/backend/outputs
    environment:
      - PORT=8888
      - HOST=0.0.0.0
      - DATABASE_PATH=./data/ogdrip.db
      - ADMIN_TOKEN=${ADMIN_TOKEN}
      - BROWSER_TIMEOUT=${BROWSER_TIMEOUT:-30}
    restart: unless-stopped
    healthcheck:
      test:
        ['CMD', 'wget', '--no-verbose', '--tries=1', '--spider', 'http://localhost:8888/api/health']
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  nginx:
    image: nginx:alpine
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
      - backend-outputs:/usr/share/nginx/html/outputs:ro
    depends_on:
      - ogdrip
    restart: unless-stopped

volumes:
  backend-data:
    driver: local
  backend-outputs:
    driver: local
```

## Environment Configuration

### Environment Variables

Create a `.env` file for production deployment:

```env
# .env
# Required
ADMIN_TOKEN=your_secure_admin_token_here

# Optional
PORT=8888
HOST=0.0.0.0
DATABASE_PATH=./data/ogdrip.db
BROWSER_TIMEOUT=30
MAX_CONCURRENT_GENERATIONS=3

# Security
CORS_ORIGINS=https://yourdomain.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600

# Monitoring
SENTRY_DSN=your_sentry_dsn_here
LOG_LEVEL=info
```

### Docker Environment File

```env
# docker.env
COMPOSE_PROJECT_NAME=ogdrip
COMPOSE_FILE=docker-compose.prod.yml

# Application
ADMIN_TOKEN=your_secure_admin_token_here
BROWSER_TIMEOUT=60
MAX_CONCURRENT_GENERATIONS=2

# Database
DATABASE_PATH=./data/ogdrip.db

# Monitoring
LOG_LEVEL=info
```

## Deployment Commands

### Development Deployment

```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f

# Rebuild after changes
docker-compose up -d --build

# Stop services
docker-compose down

# Clean up (removes volumes)
docker-compose down -v
```

### Production Deployment

```bash
# Build production image
docker build -f Dockerfile.production -t ogdrip:latest .

# Run with Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Or run single container
docker run -d \
  --name ogdrip \
  -p 8888:8888 \
  -v ogdrip-data:/app/backend/data \
  -v ogdrip-outputs:/app/backend/outputs \
  -e ADMIN_TOKEN=your_secure_token \
  ogdrip:latest
```

### Health Checks

```bash
# Check container health
docker ps
docker inspect ogdrip --format='{{.State.Health.Status}}'

# Check application health
curl http://localhost:8888/api/health

# View application logs
docker logs ogdrip -f

# Check resource usage
docker stats ogdrip
```

## Nginx Configuration

### Basic Nginx Configuration

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server ogdrip:8888;
    }

    server {
        listen 80;
        server_name localhost;

        # Frontend static files
        location / {
            root /usr/share/nginx/html;
            try_files $uri $uri/ /index.html;
        }

        # API proxy
        location /api/ {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Generated images
        location /outputs/ {
            alias /usr/share/nginx/html/outputs/;
            expires 1d;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

### Production Nginx with SSL

```nginx
# nginx.prod.conf
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server ogdrip:8888;
    }

    # Redirect HTTP to HTTPS
    server {
        listen 80;
        server_name yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name yourdomain.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

        # Security headers
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

        # Frontend
        location / {
            root /usr/share/nginx/html;
            try_files $uri $uri/ /index.html;
        }

        # API
        location /api/ {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Generated images
        location /outputs/ {
            alias /usr/share/nginx/html/outputs/;
            expires 7d;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

## Volume Management

### Data Persistence

```bash
# Create named volumes
docker volume create ogdrip-data
docker volume create ogdrip-outputs

# Backup data volume
docker run --rm -v ogdrip-data:/data -v $(pwd):/backup alpine \
    tar czf /backup/ogdrip-data-$(date +%Y%m%d).tar.gz -C /data .

# Restore data volume
docker run --rm -v ogdrip-data:/data -v $(pwd):/backup alpine \
    tar xzf /backup/ogdrip-data-20240115.tar.gz -C /data

# Backup outputs volume
docker run --rm -v ogdrip-outputs:/outputs -v $(pwd):/backup alpine \
    tar czf /backup/ogdrip-outputs-$(date +%Y%m%d).tar.gz -C /outputs .
```

### Cleanup

```bash
# Clean up old images
find ./backend/outputs -name "*.png" -mtime +7 -delete

# Docker cleanup
docker system prune -f
docker volume prune -f
docker image prune -f

# Remove specific volumes
docker volume rm ogdrip-data ogdrip-outputs
```

## Monitoring and Logging

### Log Configuration

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  ogdrip:
    # ... existing configuration
    logging:
      driver: 'json-file'
      options:
        max-size: '10m'
        max-file: '3'

  nginx:
    # ... existing configuration
    logging:
      driver: 'json-file'
      options:
        max-size: '10m'
        max-file: '3'
```

### Health Monitoring

```bash
#!/bin/bash
# health-check.sh

# Check container status
if ! docker ps | grep -q ogdrip; then
    echo "ERROR: OG Drip container is not running"
    exit 1
fi

# Check application health
if ! curl -f http://localhost:8888/api/health >/dev/null 2>&1; then
    echo "ERROR: OG Drip API is not responding"
    exit 1
fi

# Check disk space
USAGE=$(df /var/lib/docker | tail -1 | awk '{print $5}' | sed 's/%//')
if [ $USAGE -gt 90 ]; then
    echo "WARNING: Disk usage is at ${USAGE}%"
fi

echo "OK: All health checks passed"
```

## Troubleshooting

### Common Docker Issues

**Container won't start:**

```bash
# Check logs
docker logs ogdrip

# Check container configuration
docker inspect ogdrip

# Run container interactively
docker run -it --rm ogdrip /bin/sh
```

**Permission issues:**

```bash
# Fix file permissions
sudo chown -R $USER:$USER ./backend/data ./backend/outputs

# Or run with correct user
docker run --user $(id -u):$(id -g) ogdrip
```

**Out of memory:**

```bash
# Limit container memory
docker run -m 2g ogdrip

# Check memory usage
docker stats ogdrip
```

**Browser launch failures:**

```bash
# Add required Chrome flags
docker run -e CHROME_FLAGS="--no-sandbox --disable-gpu" ogdrip

# Check Chrome dependencies
docker run -it ogdrip /bin/sh
chromium --version
```

### Performance Optimization

```bash
# Use multi-stage builds for smaller images
docker build --target production -t ogdrip:prod .

# Enable BuildKit for faster builds
export DOCKER_BUILDKIT=1
docker build -t ogdrip .

# Use build cache
docker build --cache-from ogdrip:latest -t ogdrip:new .

# Optimize for production
docker run --memory=4g --cpus=2 ogdrip
```

## Security Considerations

### Container Security

```dockerfile
# Use non-root user
RUN addgroup -g 1001 -S ogdrip && \
    adduser -S ogdrip -u 1001
USER ogdrip

# Remove unnecessary packages
RUN apk del build-dependencies && \
    rm -rf /var/cache/apk/*

# Set read-only filesystem
docker run --read-only --tmpfs /tmp ogdrip
```

### Network Security

```yaml
# docker-compose.security.yml
services:
  ogdrip:
    networks:
      - internal
    # Don't expose ports directly

  nginx:
    ports:
      - '80:80'
      - '443:443'
    networks:
      - internal
      - external

networks:
  internal:
    driver: bridge
    internal: true
  external:
    driver: bridge
```

## Production Checklist

- [ ] Use production Dockerfile
- [ ] Set secure admin token
- [ ] Configure SSL certificates
- [ ] Set up log rotation
- [ ] Configure health checks
- [ ] Set resource limits
- [ ] Enable security headers
- [ ] Set up backup strategy
- [ ] Configure monitoring
- [ ] Test disaster recovery

---

_For more deployment options, see our [Production Setup Guide](production.md) or
[Coolify Deployment](coolify.md)._
