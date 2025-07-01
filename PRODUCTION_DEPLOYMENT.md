# Production Deployment Guide

This guide covers deploying OG Drip to production using Coolify with nixpacks.

## Overview

The application has been optimized for reliable production deployment with the following improvements:

- **Compiled binaries**: Uses pre-built Go binaries instead of `go run` for better performance and reliability
- **Chromium support**: Includes Chromium browser for headless operations
- **Process management**: Robust process management with graceful shutdown
- **Health monitoring**: Comprehensive health checks and monitoring
- **Production configuration**: Optimized environment configurations

## Deployment Architecture

```
Coolify (Nixpacks)
├── Build Phase
│   ├── Install dependencies (Node.js, Go, Chromium, fonts)
│   ├── Build frontend (if present)
│   └── Compile Go backend binary
├── Runtime
│   ├── start.sh (main entry point)
│   ├── Backend service (compiled binary)
│   ├── Frontend service (if present)
│   └── Health monitoring
└── Management
    ├── process-manager.sh (process control)
    ├── healthcheck.sh (health monitoring)
    └── Logging and cleanup
```

## Files and Configuration

### Core Files

- `start.sh` - Main startup script with improved process management
- `nixpacks.toml` - Build configuration with Chromium and dependencies
- `coolify.json` - Coolify deployment configuration
- `process-manager.sh` - Process management utility
- `healthcheck.sh` - Health monitoring script

### Environment Configuration

- `backend/.env.production` - Backend production environment
- `frontend/.env.production` - Frontend production environment

## Key Improvements

### 1. Start Script (`start.sh`)

**Before**: Used `go run .` (inefficient, unreliable)
**After**: Uses compiled binary with proper process management

- Uses pre-compiled binary from nixpacks build
- Graceful shutdown with signal handling
- Better error handling and logging
- Frontend detection and conditional startup

### 2. Nixpacks Configuration (`nixpacks.toml`)

**New Dependencies Added**:
- `chromium` - For headless browser operations
- `font-manager`, `fontconfig`, `freetype` - Font support
- `liberation_ttf`, `dejavu_fonts` - Font libraries

**Environment Variables**:
- `CHROME_PATH` - Path to Chromium binary
- `FONTCONFIG_PATH` - Font configuration
- Production-optimized settings

### 3. Process Management (`process-manager.sh`)

**Features**:
- Start/stop/restart commands
- Status monitoring
- Health checking
- Log management
- Cleanup utilities

**Usage**:
```bash
./process-manager.sh start    # Start service
./process-manager.sh stop     # Stop service
./process-manager.sh status   # Check status
./process-manager.sh health   # Run health check
./process-manager.sh logs     # View logs
```

### 4. Health Monitoring (`healthcheck.sh`)

**Features**:
- API endpoint health checking
- Retry logic with exponential backoff
- Detailed error reporting
- JSON response validation

### 5. Coolify Configuration (`coolify.json`)

**Improvements**:
- Switched from docker-compose to nixpacks
- Enhanced health check configuration
- Resource limits and reservations
- Better restart policies
- Monitoring integration

## Deployment Steps

### 1. Coolify Setup

1. Create new service in Coolify
2. Choose "Build Pack: Nixpacks"
3. Connect your Git repository
4. Configure environment variables

### 2. Environment Variables

**Required Variables**:
```bash
# Backend
BASE_URL=https://your-domain.com
PORT=8888
LOG_LEVEL=info
CHROME_DISABLE_SANDBOX=true

# Frontend (if present)
PUBLIC_BACKEND_URL=https://your-domain.com
NODE_ENV=production
```

**Optional Variables**:
```bash
# Performance
MAX_QUEUE_SIZE=50
MAX_CONCURRENT_REQUESTS=10
REQUEST_TIMEOUT=30s

# Security
ENABLE_CORS=true
CORS_ORIGINS=https://your-domain.com

# Monitoring
ENABLE_METRICS=true
HEALTH_CHECK_INTERVAL=30s
```

### 3. Domain Configuration

1. Add your domain in Coolify
2. Configure SSL certificates (automatic with Let's Encrypt)
3. Update `BASE_URL` environment variable

### 4. Storage Configuration

Persistent volumes are automatically configured:
- `/app/outputs` - Generated images and files
- `/app/data` - Database and application data
- `/app/logs` - Application logs

## Monitoring and Maintenance

### Health Checks

The service includes comprehensive health monitoring:

- **Coolify Integration**: Automatic health checks via `/api/health`
- **Manual Checks**: Use `./healthcheck.sh` for detailed status
- **Process Monitoring**: Use `./process-manager.sh status`

### Logs

Access logs through:
- Coolify dashboard (real-time logs)
- `./process-manager.sh logs` (local logs)
- Log files in `/app/logs` volume

### Performance Monitoring

Monitor key metrics:
- Response times via health checks
- Memory usage (limited to 2GB)
- CPU usage (limited to 2 cores)
- Queue size and processing times

## Troubleshooting

### Common Issues

1. **Chromium not found**
   - Check `CHROME_PATH` environment variable
   - Verify nixpacks included Chromium in build

2. **Frontend not starting**
   - Check if frontend directory exists and has package.json
   - Service will run backend-only if frontend is missing

3. **Database issues**
   - Check persistent volume mounts
   - Verify permissions on `/app/data` directory

4. **Performance issues**
   - Monitor resource usage in Coolify
   - Adjust `MAX_QUEUE_SIZE` and `MAX_CONCURRENT_REQUESTS`

### Debug Commands

```bash
# Check service status
./process-manager.sh status

# Run health check
./healthcheck.sh

# View recent logs
./process-manager.sh logs

# Restart service
./process-manager.sh restart

# Clean up and restart
./process-manager.sh cleanup
./process-manager.sh start
```

## Security Considerations

- CORS is properly configured for production domains
- Chromium runs with appropriate sandboxing
- Resource limits prevent resource exhaustion
- Environment variables are properly isolated

## Backup and Recovery

- Database files are stored in persistent volume `/app/data`
- Generated files are stored in persistent volume `/app/outputs`
- Configuration is version-controlled in Git
- Use Coolify's backup features for complete system backup