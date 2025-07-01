# Common Issues & Troubleshooting

This guide covers the most frequently encountered issues when setting up, developing, and deploying
OG Drip, along with their solutions.

## Installation & Setup Issues

### Node.js Version Issues

**Problem**: Getting errors about unsupported Node.js version

```
error This project requires Node.js >= 22.13.0
```

**Solutions**:

```bash
# Check current Node.js version
node --version

# Install Node.js 22.13.0 or later
# Using nvm (recommended)
nvm install 22.13.0
nvm use 22.13.0

# Using direct download
# Visit https://nodejs.org/ and download the latest LTS version
```

### pnpm Installation Issues

**Problem**: `pnpm` command not found

```
bash: pnpm: command not found
```

**Solutions**:

```bash
# Method 1: Use corepack (recommended)
corepack enable
corepack use pnpm@10.5.2

# Method 2: Install globally
npm install -g pnpm@10.5.2

# Method 3: Use npm dlx
npx pnpm install
```

### Go Module Issues

**Problem**: Go modules not downloading or building

```
go: module example.com/ogdrip: reading go.mod: open go.mod: no such file or directory
```

**Solutions**:

```bash
# Navigate to backend directory
cd backend

# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy modules
go mod tidy

# Verify Go version
go version  # Should be 1.24 or later
```

## Development Server Issues

### Port Already in Use

**Problem**: Development servers can't start due to port conflicts

```
Error: listen EADDRINUSE: address already in use :::3000
Error: listen EADDRINUSE: address already in use :::8888
```

**Solutions**:

```bash
# Find what's using the ports
lsof -i :3000
lsof -i :8888

# Kill the processes
kill -9 <PID>

# Or use different ports
# Frontend: Edit astro.config.mjs
export default defineConfig({
  server: { port: 3001 }
});

# Backend: Set PORT environment variable
PORT=8889 go run main.go
```

### Frontend Build Errors

**Problem**: Astro or Svelte build failures

```
[ERROR] Could not resolve "@astrojs/svelte"
```

**Solutions**:

```bash
# Clear node_modules and reinstall
rm -rf node_modules pnpm-lock.yaml
pnpm install

# Clear Astro cache
rm -rf frontend/.astro

# Verify Astro configuration
cd frontend
pnpm astro check
```

### Backend Compilation Errors

**Problem**: Go compilation failures

```
package chromedp is not in GOROOT
```

**Solutions**:

```bash
cd backend

# Ensure Go modules are initialized
go mod init ogdrip-backend

# Download missing dependencies
go mod download

# Update dependencies
go get -u github.com/chromedp/chromedp

# Build to test
go build -o ogdrip main.go
```

## Database Issues

### SQLite Permission Errors

**Problem**: Cannot create or write to database

```
unable to open database file: permission denied
```

**Solutions**:

```bash
# Create data directory with proper permissions
mkdir -p backend/data
chmod 755 backend/data

# Check file permissions
ls -la backend/data/

# Fix ownership if needed (Linux/macOS)
sudo chown -R $USER:$USER backend/data/
```

### Database Corruption

**Problem**: Database file is corrupted

```
database disk image is malformed
```

**Solutions**:

```bash
# Backup existing database (if possible)
cp backend/data/ogdrip.db backend/data/ogdrip.db.backup

# Remove corrupted database (development only)
rm backend/data/ogdrip.db

# Restart backend to recreate database
cd backend
go run main.go
```

### Migration Issues

**Problem**: Database schema is outdated

```
no such column: new_column_name
```

**Solutions**:

```bash
# Run migrations manually (if implemented)
cd backend
go run main.go -migrate

# Or recreate database (development only)
rm backend/data/ogdrip.db
go run main.go -init-db
```

## Browser & ChromeDP Issues

### Chrome Dependencies Missing

**Problem**: ChromeDP fails to launch browser

```
chrome failed to start: exit status 1
```

**Solutions**:

**Linux (Ubuntu/Debian)**:

```bash
# Install Chrome dependencies
sudo apt-get update
sudo apt-get install -y \
    chromium-browser \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libgbm1 \
    libxss1 \
    libnss3
```

**macOS**:

```bash
# Install Chromium via Homebrew
brew install chromium

# Or install Google Chrome manually
# Download from https://www.google.com/chrome/
```

**Docker**:

```dockerfile
# Add to Dockerfile
RUN apt-get update && apt-get install -y \
    chromium \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*
```

### Browser Launch Timeout

**Problem**: Browser takes too long to start

```
context deadline exceeded
```

**Solutions**:

```bash
# Increase timeout in backend configuration
export BROWSER_TIMEOUT=60  # seconds

# Or modify Go code
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
```

### Headless Browser Issues

**Problem**: Browser crashes or hangs

```
chrome crashed
```

**Solutions**:

```bash
# Add Chrome flags for stability
--no-sandbox
--disable-gpu
--disable-dev-shm-usage
--disable-extensions
--remote-debugging-port=9222
```

## Network & API Issues

### CORS Errors

**Problem**: Frontend cannot connect to backend

```
Access to fetch at 'http://localhost:8888' from origin 'http://localhost:3000' has been blocked by CORS policy
```

**Solutions**:

```go
// In backend Go code, ensure CORS is properly configured
func setupCORS() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })
}
```

### API Connection Refused

**Problem**: Frontend cannot reach backend API

```
net::ERR_CONNECTION_REFUSED
```

**Solutions**:

```bash
# Verify backend is running
curl http://localhost:8888/api/health

# Check backend logs for errors
cd backend
go run main.go -debug

# Verify environment variables
echo $PUBLIC_BACKEND_URL  # Should be http://localhost:8888
```

### Rate Limiting Issues

**Problem**: API requests being rate limited

```
HTTP 429 Too Many Requests
```

**Solutions**:

```bash
# Increase rate limits (development only)
export RATE_LIMIT_REQUESTS=1000
export RATE_LIMIT_WINDOW=3600

# Or disable rate limiting for development
export DISABLE_RATE_LIMITING=true
```

## Docker Issues

### Docker Permission Errors

**Problem**: Permission denied when using Docker

```
permission denied while trying to connect to the Docker daemon socket
```

**Solutions**:

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in, or run:
newgrp docker

# Verify Docker access
docker run hello-world
```

### Docker Build Failures

**Problem**: Docker build fails

```
failed to solve with frontend dockerfile.v0
```

**Solutions**:

```bash
# Clear Docker cache
docker system prune -a

# Build with no cache
docker build --no-cache -t ogdrip .

# Check Dockerfile syntax
docker build --dry-run -t ogdrip .
```

### Container Startup Issues

**Problem**: Container exits immediately

```
container exited with code 1
```

**Solutions**:

```bash
# Check container logs
docker logs <container_id>

# Run container interactively
docker run -it ogdrip /bin/sh

# Check environment variables
docker run ogdrip env
```

## Performance Issues

### Slow Image Generation

**Problem**: Image generation takes too long

```
Image generation timeout after 30 seconds
```

**Solutions**:

```bash
# Increase timeouts
export BROWSER_TIMEOUT=60
export GENERATION_TIMEOUT=120

# Optimize browser flags
--disable-background-timer-throttling
--disable-backgrounding-occluded-windows
--disable-renderer-backgrounding

# Monitor system resources
top
htop
docker stats  # If using Docker
```

### High Memory Usage

**Problem**: Application consumes too much memory

```
out of memory: cannot allocate memory
```

**Solutions**:

```bash
# Monitor memory usage
free -h
docker stats

# Limit browser instances
export MAX_CONCURRENT_GENERATIONS=2

# Add memory limits (Docker)
docker run -m 2g ogdrip

# Optimize Go garbage collection
export GOGC=100
```

### Disk Space Issues

**Problem**: Running out of disk space

```
no space left on device
```

**Solutions**:

```bash
# Check disk usage
df -h
du -sh backend/outputs/

# Clean up old images
find backend/outputs -name "*.png" -mtime +7 -delete

# Set up log rotation
logrotate /etc/logrotate.d/ogdrip
```

## Testing Issues

### Test Failures

**Problem**: Tests failing unexpectedly

```
Test suite failed to run
```

**Solutions**:

```bash
# Clear test cache
cd frontend
rm -rf node_modules/.cache

# Run tests in isolation
pnpm test --no-coverage --reporter=verbose

# Check test database
cd backend
go test -v ./...
```

### E2E Test Issues

**Problem**: Playwright tests failing

```
browserType.launch: Executable doesn't exist
```

**Solutions**:

```bash
# Install Playwright browsers
cd frontend
pnpm exec playwright install

# Run tests in headed mode for debugging
pnpm test:e2e --headed

# Check Playwright configuration
npx playwright doctor
```

## Deployment Issues

### Environment Variables

**Problem**: Environment variables not loading

```
undefined is not a function (reading 'PUBLIC_BACKEND_URL')
```

**Solutions**:

```bash
# Verify environment files exist
ls -la frontend/.env backend/.env

# Check variable names (must start with PUBLIC_ for frontend)
grep PUBLIC_ frontend/.env

# Restart development servers after changes
pnpm dev
```

### Build Deployment Issues

**Problem**: Production build fails

```
Build failed with errors
```

**Solutions**:

```bash
# Test build locally
pnpm build

# Check build output
ls -la frontend/dist/

# Verify all dependencies are installed
pnpm install --frozen-lockfile

# Check for TypeScript errors
cd frontend
pnpm astro check
```

## Getting Help

### Diagnostic Commands

Run these commands to gather information for bug reports:

```bash
# System information
node --version
go version
pnpm --version
docker --version

# Project status
pnpm list
cd backend && go list -m all

# Check running processes
ps aux | grep -E "(node|go|chrome)"
lsof -i :3000,8888

# Disk usage
df -h
du -sh backend/outputs backend/data
```

### Log Collection

```bash
# Frontend logs (browser console)
# Open browser dev tools > Console tab

# Backend logs
cd backend
go run main.go 2>&1 | tee backend.log

# Docker logs
docker logs <container_name> > docker.log 2>&1

# System logs (Linux)
journalctl -u ogdrip --since "1 hour ago"
```

### Creating Bug Reports

When reporting issues, include:

1. **Environment details**: OS, Node.js version, Go version
2. **Steps to reproduce**: Exact commands and actions
3. **Expected vs actual behavior**: What should happen vs what happens
4. **Error messages**: Complete error output
5. **Logs**: Relevant log files
6. **Configuration**: Environment variables (redact secrets)

### Community Support

- **GitHub Issues**: [Create an issue](https://github.com/yourusername/ogdrip/issues/new)
- **Discussions**: [Join discussions](https://github.com/yourusername/ogdrip/discussions)
- **FAQ**: Check our [FAQ](faq.md) for common questions

---

_Still having issues?
[Create a detailed bug report](https://github.com/yourusername/ogdrip/issues/new?template=bug_report.md)
with the diagnostic information above._
