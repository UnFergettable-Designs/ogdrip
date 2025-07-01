# Development Setup

This guide will help you set up a local development environment for contributing to OG Drip.

## Prerequisites

Before you begin, ensure you have the following tools installed:

### Required Software

- **Node.js** >= 22.13.0 ([Download](https://nodejs.org/))
- **pnpm** >= 10.5.2 (Package manager)
- **Go** >= 1.24 ([Download](https://golang.org/dl/))
- **Git** for version control
- **VS Code** (recommended) or your preferred editor

### Optional but Recommended

- **Docker** and **Docker Compose** for containerized development
- **Chrome/Chromium** for testing (ChromeDP dependency)

## Initial Setup

### 1. Clone the Repository

```bash
# Clone your fork of the repository
git clone https://github.com/YOUR_USERNAME/ogdrip.git
cd ogdrip

# Add the upstream remote
git remote add upstream https://github.com/original-owner/ogdrip.git
```

### 2. Install Package Manager

```bash
# Enable corepack (recommended)
corepack enable

# Or install pnpm globally
npm install -g pnpm@10.5.2
```

### 3. Install Dependencies

```bash
# Install all dependencies for the monorepo
pnpm install

# Verify Go modules
cd backend
go mod download
go mod tidy
cd ..
```

### 4. Environment Configuration

```bash
# Copy environment files
cp frontend/.env.example frontend/.env
cp backend/.env.example backend/.env

# Edit the files with your local configuration
nano frontend/.env
nano backend/.env
```

#### Frontend Environment (`.env`)

```env
# Backend API URL for development
PUBLIC_BACKEND_URL=http://localhost:8888
BACKEND_URL=http://localhost:8888

# Optional: Development tools
PUBLIC_SENTRY_DSN=your_development_sentry_dsn
```

#### Backend Environment (`.env`)

```env
# Server configuration
PORT=8888
HOST=localhost

# Database
DATABASE_PATH=./data/ogdrip.db

# Admin access (use a secure token)
ADMIN_TOKEN=dev_admin_token_change_in_production

# Optional: Development services
SENTRY_DSN=your_development_sentry_dsn
```

### 5. Create Required Directories

```bash
# Create necessary directories
mkdir -p backend/data backend/outputs
mkdir -p frontend/public/outputs
```

## Development Workflow

### Starting Development Servers

```bash
# Start all services (frontend + backend)
pnpm dev

# Or start services individually
pnpm dev:frontend  # Starts Astro dev server on port 3000
pnpm dev:backend   # Starts Go API server on port 8888
```

### Available Scripts

```bash
# Development
pnpm dev           # Start all development servers
pnpm dev:frontend  # Start only frontend
pnpm dev:backend   # Start only backend

# Building
pnpm build         # Build all packages
pnpm build:frontend # Build only frontend
pnpm build:backend  # Build only backend

# Testing
pnpm test          # Run all tests
pnpm test:frontend # Run frontend tests
pnpm test:backend  # Run backend tests
pnpm test:coverage # Run tests with coverage

# Linting and Formatting
pnpm lint          # Lint all packages
pnpm lint:fix      # Fix linting issues
pnpm format        # Format code with Prettier

# Validation
pnpm validate      # Run all validation checks
```

## IDE Setup

### VS Code Configuration

Install the following extensions:

```json
{
  "recommendations": [
    "astro-build.astro-vscode",
    "svelte.svelte-vscode",
    "golang.go",
    "esbenp.prettier-vscode",
    "dbaeumer.vscode-eslint",
    "bradlc.vscode-tailwindcss",
    "ms-vscode.vscode-typescript-next"
  ]
}
```

### Workspace Settings

Create `.vscode/settings.json`:

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go"
  },
  "[astro]": {
    "editor.defaultFormatter": "astro-build.astro-vscode"
  },
  "[svelte]": {
    "editor.defaultFormatter": "svelte.svelte-vscode"
  }
}
```

## Project Structure

```
ogdrip/
├── frontend/          # Astro + Svelte frontend
│   ├── src/
│   │   ├── components/    # Svelte components
│   │   ├── layouts/       # Astro layouts
│   │   ├── pages/         # Astro pages/routes
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets
│   └── astro.config.mjs   # Astro configuration
├── backend/           # Go API server
│   ├── *.go              # Go source files
│   ├── data/             # SQLite database
│   ├── outputs/          # Generated images
│   └── go.mod            # Go module file
├── shared/            # Shared TypeScript types
│   ├── types.ts          # Type definitions
│   └── index.ts          # Exports
└── docs/              # Documentation
```

## Development Guidelines

### Code Style

#### Frontend (TypeScript/Svelte/Astro)

- Use TypeScript for all new code
- Follow Svelte 5 runes syntax (`$state`, `$derived`)
- Use semantic HTML elements
- Follow accessibility guidelines (WCAG 2.2 AA)
- Use REM for sizing, HSLA for colors

#### Backend (Go)

- Follow Go conventions and idioms
- Use proper error handling with context
- Implement timeouts for ChromeDP operations
- Use prepared statements for database queries
- Include comprehensive logging

### Testing

#### Frontend Testing

```bash
# Run unit tests
cd frontend
pnpm test

# Run tests in watch mode
pnpm test:watch

# Run tests with coverage
pnpm test:coverage
```

#### Backend Testing

```bash
# Run Go tests
cd backend
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...
```

### Database Development

#### SQLite Management

```bash
# Access database directly
sqlite3 backend/data/ogdrip.db

# View schema
.schema

# Common queries
SELECT * FROM generations ORDER BY created_at DESC LIMIT 10;
```

#### Database Migrations

```bash
# Apply migrations (if implemented)
cd backend
go run main.go -migrate

# Reset database (development only)
rm backend/data/ogdrip.db
go run main.go -init-db
```

## Debugging

### Frontend Debugging

- Use browser developer tools
- Astro dev server provides detailed error pages
- Check console for JavaScript errors
- Use Svelte DevTools extension

### Backend Debugging

```bash
# Run with debug logging
cd backend
go run main.go -debug

# Use delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

### Common Issues

#### Port Conflicts

```bash
# Check what's using a port
lsof -i :3000
lsof -i :8888

# Kill process if needed
kill -9 <PID>
```

#### ChromeDP Issues

```bash
# Install Chrome dependencies (Linux)
sudo apt-get update
sudo apt-get install -y chromium-browser

# macOS (if using Homebrew)
brew install chromium
```

#### Go Module Issues

```bash
cd backend
go clean -modcache
go mod download
go mod tidy
```

## Performance Optimization

### Development Performance

- Use `pnpm dev` for hot reloading
- Enable Go build cache: `export GOCACHE=$(go env GOCACHE)`
- Use SSD storage for better I/O performance

### Profiling

```bash
# Profile Go application
cd backend
go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Contributing Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Follow coding standards
- Write tests for new functionality
- Update documentation as needed

### 3. Test Changes

```bash
# Run all tests
pnpm test

# Run linting
pnpm lint

# Build to ensure no errors
pnpm build
```

### 4. Commit Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add new feature description"
```

### 5. Push and Create PR

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create pull request on GitHub
```

## Troubleshooting

### Getting Help

1. Check the [FAQ](../troubleshooting/faq.md)
2. Search existing [GitHub Issues](https://github.com/yourusername/ogdrip/issues)
3. Join community discussions
4. Ask questions in pull requests

### Common Development Issues

**Build Failures:**

- Ensure all dependencies are installed
- Check Node.js and Go versions
- Clear build caches: `pnpm clean`

**Test Failures:**

- Run tests individually to isolate issues
- Check test database setup
- Ensure test data is properly cleaned up

**Linting Errors:**

- Run `pnpm lint:fix` to auto-fix issues
- Check ESLint and Prettier configurations
- Ensure code follows project conventions

## Next Steps

After setting up your development environment:

1. **Explore the codebase** - Start with the main entry points
2. **Run the test suite** - Ensure everything works correctly
3. **Make a small change** - Try fixing a small issue or adding a minor feature
4. **Read the architecture docs** - Understand the system design
5. **Join the community** - Participate in discussions and code reviews

---

_Need help with setup? Check our [troubleshooting guide](../troubleshooting/common-issues.md) or
[create an issue](https://github.com/yourusername/ogdrip/issues/new)._
