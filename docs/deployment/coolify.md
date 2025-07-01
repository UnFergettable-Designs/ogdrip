# Coolify Deployment Guide

This guide covers deploying OG Drip on Coolify using nixpacks for automatic builds and deployments.

## Overview

OG Drip is optimized for Coolify deployment with:
- **nixpacks** for automatic build detection
- **Automatic SSL** with Let's Encrypt
- **Health checks** and monitoring
- **Persistent volumes** for data and generated images
- **Environment management** for production configuration

## Prerequisites

- **Coolify instance** (self-hosted or managed)
- **Git repository** connected to Coolify
- **Domain name** (optional but recommended)
- **2GB RAM minimum**, 4GB recommended

## Quick Deployment

### 1. Create New Service

In your Coolify dashboard:

1. Click **"New Service"**
2. Choose **"Source: GitHub"** (or your Git provider)
3. Select your **ogdrip repository**
4. Choose **"Build Pack: Nixpacks"**

### 2. Configure Environment Variables

Set these required environment variables in Coolify:

#### Required Variables
```env
# Security - Generate a secure token
ADMIN_TOKEN=your_secure_admin_token_here

# Backend URLs - Update with your domain
PUBLIC_BACKEND_URL=https://your-domain.com
BACKEND_URL=https://your-domain.com

# CORS - Match your domain
CORS_ORIGINS=https://your-domain.com
```

#### Optional Variables
```env
# Performance
BROWSER_TIMEOUT=60
MAX_CONCURRENT_GENERATIONS=2

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600

# Monitoring
SENTRY_DSN=your_sentry_dsn_here
LOG_LEVEL=info
```

### 3. Add Domain

In Coolify:
1. Go to **"Domains"** tab
2. Add your domain (e.g., `og-drip.com`)
3. Enable **"Auto SSL"** for HTTPS

### 4. Deploy

Click **"Deploy"** and Coolify will:
- Build using nixpacks
- Set up SSL certificates
- Configure reverse proxy
- Start health monitoring

## Configuration Files

### nixpacks.toml

The project includes a `nixpacks.toml` file that configures the build:

```toml
[phases.setup]
nixPkgs = [
  "nodejs_22",
  "go_1_24",
  "chromium",
  "font-manager",
  "fontconfig",
  "freetype",
  "liberation_ttf",
  "dejavu_fonts",
  "git"
]

[phases.install]
cmds = [
  "pnpm install --frozen-lockfile"
]

[phases.build]
cmds = [
  "pnpm build",
  "cd backend && mkdir -p build && go build -o build/ogdrip-backend *.go",
  "chmod +x backend/build/ogdrip-backend"
]

[start]
cmd = "./start.sh"
```

### coolify.json

The project includes Coolify-specific configuration:

```json
{
  "coolify": {
    "version": "1.0.0",
    "name": "Open Graph Generator",
    "healthCheck": {
      "path": "/api/health",
      "port": 8888,
      "interval": 30,
      "timeout": 10,
      "retries": 3,
      "startPeriod": 60
    },
    "volumes": [
      {
        "name": "outputs",
        "path": "/app/outputs",
        "persistent": true
      },
      {
        "name": "database",
        "path": "/app/data",
        "persistent": true
      }
    ]
  }
}
```

## Environment Setup

### Creating Environment Files

Before deployment, create these environment files:

#### Backend Environment (`backend/.env.production`)
```env
# Server Configuration
PORT=8888
HOST=0.0.0.0

# Database
DATABASE_PATH=./data/ogdrip.db

# Security - CHANGE THIS
ADMIN_TOKEN=your_secure_admin_token_here

# Browser Configuration
BROWSER_TIMEOUT=30
MAX_CONCURRENT_GENERATIONS=3

# CORS Configuration
CORS_ORIGINS=https://yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600

# Chrome Configuration
CHROME_PATH=/nix/store/*-chromium-*/bin/chromium
DISPLAY=:99
```

#### Frontend Environment (`frontend/.env.production`)
```env
# Backend API URLs - UPDATE THESE
PUBLIC_BACKEND_URL=https://yourdomain.com
BACKEND_URL=https://yourdomain.com

# Build Configuration
NODE_ENV=production
```

### Setting Environment Variables in Coolify

1. Go to your service in Coolify
2. Click **"Environment Variables"** tab
3. Add each variable individually or bulk import

**Bulk Import Format:**
```
ADMIN_TOKEN=your_secure_admin_token_here
PUBLIC_BACKEND_URL=https://your-domain.com
BACKEND_URL=https://your-domain.com
CORS_ORIGINS=https://your-domain.com
BROWSER_TIMEOUT=60
MAX_CONCURRENT_GENERATIONS=2
```

## Deployment Process

### Build Process

Coolify with nixpacks will:

1. **Setup Phase**: Install Node.js 22, Go 1.24, Chromium, and fonts
2. **Install Phase**: Run `pnpm install --frozen-lockfile`
3. **Build Phase**:
   - Build frontend with `pnpm build`
   - Compile Go backend binary
   - Set executable permissions

### Start Process

The application starts using `start.sh` which:
- Creates necessary directories
- Starts backend service
- Starts frontend (if applicable)
- Handles graceful shutdown

### Health Checks

Coolify monitors the application using:
- **Health endpoint**: `/api/health` on port 8888
- **Interval**: Every 30 seconds
- **Timeout**: 10 seconds
- **Retries**: 3 attempts
- **Start period**: 60 seconds grace period

## Persistent Storage

### Volumes Configuration

The deployment uses persistent volumes for:

1. **Database Storage** (`/app/data`)
   - SQLite database file
   - Application logs
   - Configuration files

2. **Generated Images** (`/app/outputs`)
   - Open Graph images
   - Cached results
   - Temporary files

### Volume Management

```bash
# Access volumes via Coolify terminal
cd /app/data      # Database and logs
cd /app/outputs   # Generated images

# Check disk usage
df -h /app/data /app/outputs

# Backup data (run from Coolify terminal)
tar -czf backup-$(date +%Y%m%d).tar.gz data/ outputs/
```

## SSL Configuration

### Automatic SSL

Coolify automatically configures SSL when you:
1. Add a domain to your service
2. Enable "Auto SSL"
3. Ensure DNS points to your Coolify instance

### Manual SSL

For custom certificates:
1. Upload certificate files in Coolify
2. Configure SSL settings
3. Update environment variables if needed

## Monitoring & Logging

### Built-in Monitoring

Coolify provides:
- **Resource usage** graphs
- **Application logs** in real-time
- **Health status** monitoring
- **Deployment history**

### Custom Monitoring

The application includes:
- **Health check script** (`healthcheck.sh`)
- **Process manager** (`process-manager.sh`)
- **Structured logging** in JSON format

### Accessing Logs

```bash
# View logs in Coolify UI
# Or access via terminal:

# Application logs
tail -f /app/logs/application.log

# System logs
journalctl -u your-service-name -f

# Container logs
docker logs container_name -f
```

## Scaling & Performance

### Resource Allocation

Configure in Coolify:
- **Memory**: 2GB minimum, 4GB recommended
- **CPU**: 1 core minimum, 2 cores recommended
- **Storage**: 10GB minimum for images

### Performance Optimization

```env
# Environment variables for performance
BROWSER_TIMEOUT=60
MAX_CONCURRENT_GENERATIONS=2
NODE_OPTIONS=--max-old-space-size=4096
GOGC=100
```

### Horizontal Scaling

For high traffic:
1. Deploy multiple instances
2. Use load balancer
3. Share persistent volumes
4. Configure session affinity if needed

## Backup & Recovery

### Automated Backups

Set up automated backups in Coolify:
1. Configure backup schedule
2. Choose backup destination
3. Set retention policy

### Manual Backup

```bash
# Create backup
tar -czf ogdrip-backup-$(date +%Y%m%d).tar.gz \
  /app/data /app/outputs

# Upload to external storage
# (configure based on your backup solution)
```

### Disaster Recovery

1. **Data Recovery**: Restore from volume backups
2. **Configuration Recovery**: Re-import environment variables
3. **Service Recovery**: Redeploy from Git repository

## Troubleshooting

### Common Issues

**Build Failures:**
```bash
# Check build logs in Coolify
# Common causes:
# - Missing environment variables
# - Go/Node version issues
# - Dependency conflicts
```

**Runtime Issues:**
```bash
# Check application logs
# Common causes:
# - Chrome/Chromium not starting
# - Database permissions
# - Network connectivity
```

**Performance Issues:**
```bash
# Monitor resource usage
# Check for:
# - Memory leaks
# - CPU bottlenecks
# - Disk space issues
```

### Debug Mode

Enable debug logging:
```env
LOG_LEVEL=debug
```

### Health Check Debugging

```bash
# Test health endpoint manually
curl -f http://localhost:8888/api/health

# Run health check script
./healthcheck.sh

# Check process status
./process-manager.sh status
```

## Security Considerations

### Environment Security

- **Never commit** `.env.production` files
- **Use strong tokens** for ADMIN_TOKEN
- **Restrict CORS** origins to your domain
- **Enable rate limiting** in production

### Network Security

- **Use HTTPS** only in production
- **Configure firewall** rules
- **Monitor access** logs
- **Regular security** updates

### Data Security

- **Encrypt sensitive** data at rest
- **Regular backups** with encryption
- **Access control** for admin functions
- **Audit logging** for security events

## Production Checklist

Before going live:

- [ ] Set secure ADMIN_TOKEN
- [ ] Configure proper domain and SSL
- [ ] Set up persistent volumes
- [ ] Configure monitoring and alerts
- [ ] Test backup and recovery
- [ ] Set appropriate resource limits
- [ ] Configure rate limiting
- [ ] Test health checks
- [ ] Verify CORS settings
- [ ] Set up log rotation

## Advanced Configuration

### Custom Build

For custom build requirements, modify `nixpacks.toml`:

```toml
[phases.setup]
nixPkgs = ["nodejs_22", "go_1_24", "chromium", "your-package"]

[phases.custom]
cmds = ["your-custom-command"]
```

### Multiple Environments

Deploy separate instances for:
- **Staging**: Test new features
- **Production**: Live application
- **Development**: Development testing

### Integration

Connect with external services:
- **CDN**: For static asset delivery
- **Database**: External PostgreSQL if needed
- **Monitoring**: External monitoring services
- **Analytics**: Usage tracking services

---

*For more deployment options, see our [Docker Deployment Guide](docker.md) or [Production Setup Guide](production.md).*
