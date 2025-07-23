# OGDrip - Project Context for Claude

## Project Overview

**OGDrip** is an Open Graph image generator service that creates beautiful OG images and metadata for web pages. It's built as a modern, efficient service with both a Go backend and an Astro + Svelte frontend.

### Tech Stack
- **Frontend**: Astro + Svelte 5
- **Backend**: Go with ChromeDP for headless browser automation
- **Database**: SQLite
- **Build System**: Turborepo + pnpm workspaces
- **Deployment**: Docker, Coolify with nixpacks

### Repository Structure
```
ogdrip/
├── frontend/           # Astro + Svelte frontend (git submodule)
├── backend/           # Go backend service
├── shared/            # Shared TypeScript types
├── docs/              # Documentation
├── scripts/           # Deployment and setup scripts
└── [deployment files] # Docker, nixpacks, and deployment configs
```

## Key Commands

### Development
```bash
# Setup project
pnpm install
pnpm setup

# Start all services in development
pnpm dev

# Start individual services
pnpm dev:frontend    # Astro dev server on :5000
pnpm dev:backend     # Go backend on :8888

# Build all packages
pnpm build

# Lint and format
pnpm lint
pnpm lint:fix
pnpm format

# Test
pnpm test
pnpm test:coverage
```

### Docker Deployment
```bash
# Production with Docker Compose
cp .env.example .env  # Edit with your settings
docker-compose up -d

# Development with hot reload
docker-compose -f docker-compose.dev.yml up

# Build individual Docker image
docker build -f Dockerfile.production -t ogdrip:prod .
```

### Manual Backend Commands
```bash
cd backend

# Run backend directly
go run *.go

# Build backend binary
go build -o build/ogdrip-backend *.go

# Test backend
go test ./...
./test-backend.sh
```

### Process Management
```bash
# Start/stop services using process manager
./process-manager.sh start
./process-manager.sh stop
./process-manager.sh status
./process-manager.sh health

# Health check
./healthcheck.sh
```

## Important Files & Configuration

### Environment Configuration
- `.env.example` - Environment template
- `frontend/.env.example` - Frontend environment template
- `backend/.env.example` - Backend environment template

### Key Configuration Files
- `package.json` - Root package configuration with scripts
- `turbo.json` - Turborepo build configuration
- `pnpm-workspace.yaml` - pnpm workspace configuration
- `nixpacks.toml` - Nixpacks build configuration for Coolify
- `coolify.json` - Coolify deployment configuration

### Docker Files
- `Dockerfile` - Basic Docker build
- `Dockerfile.production` - Production Docker build with Chromium
- `docker-compose.yml` - Production Docker Compose setup
- `docker-compose.dev.yml` - Development Docker Compose setup

### Scripts & Management
- `start.sh` - Production startup script
- `process-manager.sh` - Process management utility
- `healthcheck.sh` - Health check script
- `validate-deployment.sh` - Deployment validation

## API Endpoints

### Backend (Port 8888)
- `GET /api/health` - Health check endpoint
- `GET /docs/` - Swagger UI documentation
- `GET /api/openapi.yaml` - OpenAPI specification
- Image generation endpoints (see backend/openapi.yaml for full API)

### Frontend (Port 5000)
- Admin interface for managing OG image generation
- Built with Astro + Svelte 5

## Development Workflow

1. **Setup**: Run `pnpm install && pnpm setup`
2. **Development**: Use `pnpm dev` for hot reload development
3. **Testing**: Run `pnpm test` before committing
4. **Linting**: Code is auto-formatted with Prettier via husky hooks
5. **Building**: Use `pnpm build` to build all packages

## Deployment Options

### 1. Coolify (Recommended)
- Uses nixpacks for automatic builds
- Configured via `nixpacks.toml` and `coolify.json`
- Supports automatic SSL and domain management
- See `docs/deployment/coolify.md` for detailed guide

### 2. Docker Compose
- Production: `docker-compose up -d`
- Development: `docker-compose -f docker-compose.dev.yml up`
- Includes persistent volumes, health checks, and resource limits

### 3. Manual Deployment
- Build with `pnpm build`
- Run backend binary and serve frontend
- See `DEPLOYMENT.md` for detailed instructions

## Troubleshooting

### Common Issues
- **Frontend submodule not initialized**: Run `git submodule update --init --recursive`
- **Build failures**: Check Node.js version (requires >=22.13.0)
- **Chromium issues**: Ensure proper Chrome/Chromium installation for headless browser
- **CORS errors**: Check CORS_ORIGINS environment variable
- **Database issues**: Ensure proper permissions on data directory

### Logs & Monitoring
- Application logs: Check Docker logs or `./logs/` directory
- Health checks: Use `./healthcheck.sh` or hit `/api/health`
- Process status: Use `./process-manager.sh status`

## Security Notes

- **ADMIN_TOKEN**: Always set a secure admin token in production
- **CORS**: Restrict CORS_ORIGINS to your actual domains
- **Rate Limiting**: Configure appropriate rate limits for your use case
- **Database**: SQLite database contains application state - ensure proper backups

## Performance Considerations

- **Concurrent Generations**: Limit via MAX_CONCURRENT_GENERATIONS
- **Browser Timeout**: Adjust BROWSER_TIMEOUT for complex pages
- **Memory**: Chromium requires significant memory - allocate 2GB+ in production
- **Storage**: Generated images are stored in outputs/ - monitor disk usage

## Environment Variables

### Required
- `ADMIN_TOKEN` - Secure token for admin access
- `PUBLIC_BACKEND_URL` - Public URL where backend is accessible
- `CORS_ORIGINS` - Allowed CORS origins

### Optional Performance
- `BROWSER_TIMEOUT` - Browser timeout in seconds (default: 60)
- `MAX_CONCURRENT_GENERATIONS` - Max concurrent generations (default: 3)
- `RATE_LIMIT_REQUESTS` - Rate limit requests per window (default: 100)
- `RATE_LIMIT_WINDOW` - Rate limit window in seconds (default: 3600)

### Optional Monitoring
- `SENTRY_DSN` - Sentry DSN for error tracking
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

## Git Submodules

The frontend is managed as a git submodule. When working with the repository:

```bash
# Initialize submodules
git submodule update --init --recursive

# Update submodules
git submodule update --remote

# Check submodule status
git submodule status
```

## Monorepo Structure

This project uses Turborepo for efficient monorepo management:
- **Root**: Orchestration and shared configuration
- **frontend/**: Astro + Svelte frontend package
- **backend/**: Go backend service
- **shared/**: Shared TypeScript types and utilities

Each package can be built and run independently, but Turborepo handles dependencies and caching for optimal build performance.
