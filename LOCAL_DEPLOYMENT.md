# Local Development Guide

This guide explains how to set up and run the Open Graph Generator service locally for development.

## Prerequisites

- Node.js 20.x or later
- pnpm 8.x or later
- Go 1.24 or later
- SQLite 3.x
- Chromium/Chrome (for backend image generation)

## Setup Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/ogdrip.git
   cd ogdrip
   ```

2. Install dependencies:

   ```bash
   pnpm install
   ```

3. Set up environment variables:

   ```bash
   # Frontend (.env)
   cp frontend/.env.example frontend/.env

   # Backend (.env)
   cp backend/.env.example backend/.env
   ```

4. Start the development servers:

   ```bash
   # Start all services
   pnpm dev:all

   # Or start services individually:
   pnpm dev:frontend  # Frontend only
   pnpm dev:backend   # Backend only
   ```

## Development Environment

### Frontend (Astro + Svelte)

- Development server runs on port 3000
- Hot module reloading enabled
- TypeScript checking
- ESLint + Prettier for code formatting

### Backend (Go)

- Development server runs on port 8888
- Live reloading with air (if installed)
- SQLite database for storage
- Swagger UI for API documentation

## Available Scripts

```bash
# Development
pnpm dev:all         # Start all services
pnpm dev:frontend    # Start frontend only
pnpm dev:backend     # Start backend only

# Building
pnpm build           # Build all services
pnpm build:frontend  # Build frontend only
pnpm build:backend   # Build backend only

# Testing
pnpm test           # Run all tests
pnpm test:frontend  # Run frontend tests
pnpm test:backend   # Run backend tests

# Linting
pnpm lint          # Lint all code
pnpm lint:fix      # Fix linting issues
```

## Development Features

### Hot Reloading

- Frontend changes are automatically reflected
- Backend requires restart for Go file changes
- Static assets are served immediately

### API Documentation

- Swagger UI available at `/swagger`
- OpenAPI spec at `/api/openapi.yaml`
- Postman collection in `/docs`

### Debugging

1. Frontend:

   - Browser DevTools
   - Vite debugging tools
   - Svelte DevTools extension

2. Backend:
   - Go debugger (delve)
   - API logs in console
   - SQLite database browser

## Common Development Tasks

### Adding Dependencies

```bash
# Frontend dependencies
cd frontend
pnpm add package-name

# Backend dependencies
cd backend
go get package-name
```

### Database Management

- SQLite database location: `backend/data/ogdrip.db`
- Use SQLite browser for direct DB access
- Backup: `cp backend/data/ogdrip.db backup.db`

### Testing

1. Unit Tests:

   ```bash
   pnpm test
   ```

2. Integration Tests:

   ```bash
   cd backend
   go test -v ./...
   ```

3. End-to-End Tests:
   ```bash
   pnpm test:e2e
   ```

## Troubleshooting

### Frontend Issues

1. Module not found:

   ```bash
   pnpm install  # Reinstall dependencies
   ```

2. TypeScript errors:
   ```bash
   pnpm clean    # Clean build cache
   pnpm build    # Rebuild
   ```

### Backend Issues

1. Database errors:

   ```bash
   rm backend/data/ogdrip.db  # Remove corrupt DB
   go run main.go             # New DB will be created
   ```

2. Chrome/Chromium not found:
   - Set CHROME_PATH in backend/.env
   - Install Chrome/Chromium if missing

## Best Practices

1. Code Style

   - Follow existing patterns
   - Use TypeScript for frontend
   - Format with prettier/gofmt
   - Run linters before committing

2. Testing

   - Write tests for new features
   - Update existing tests
   - Run full test suite before PR

3. Git Workflow

   - Branch from main
   - Use descriptive commit messages
   - Keep PRs focused and small

4. Performance
   - Optimize image sizes
   - Use proper caching
   - Monitor memory usage
