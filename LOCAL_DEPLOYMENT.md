# Local Development Guide

This guide covers setting up OG Drip for local development.

## Prerequisites

- Node.js 20+ (22+ recommended)
- Go 1.23+
- pnpm
- Chrome or Chromium browser
- Git

## Initial Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/ogdrip.git
   cd ogdrip
   ```

2. Install dependencies:

   ```bash
   pnpm install
   ```

3. Set up environment files:

   ```bash
   # Frontend environment
   cp frontend/.env.example frontend/.env

   # Backend environment
   cp backend/.env.example backend/.env
   ```

4. Configure environment variables:

   **frontend/.env**:

   ```bash
   PUBLIC_BACKEND_URL=http://localhost:8888
   BACKEND_URL=http://localhost:8888
   ```

   **backend/.env**:

   ```bash
   PORT=8888
   ADMIN_TOKEN=local-dev-token
   CHROME_PATH=/path/to/your/chrome
   DATABASE_PATH=./data/ogdrip.db
   OUTPUT_DIR=./outputs
   ```

## Development Workflow

### Starting the Development Servers

1. Start all services:

   ```bash
   pnpm dev
   ```

   This will start:

   - Frontend at http://localhost:3000
   - Backend at http://localhost:8888

2. Or start services individually:

   ```bash
   # Frontend only
   pnpm dev:frontend

   # Backend only
   pnpm dev:backend
   ```

### Development Commands

```bash
# Install dependencies
pnpm install

# Start development servers
pnpm dev

# Build all packages
pnpm build

# Run tests
pnpm test

# Lint code
pnpm lint

# Preview production build
pnpm preview
```

## Project Structure

```
ogdrip/
├── frontend/           # Astro + Svelte frontend
│   ├── src/
│   ├── public/
│   └── package.json
├── backend/           # Go backend service
│   ├── cmd/
│   ├── internal/
│   └── go.mod
├── shared/           # Shared TypeScript types
│   └── src/
└── package.json     # Root package.json
```

## Development Guidelines

### Code Style

- Frontend: Uses Prettier and ESLint
- Backend: Uses `gofmt` and `golangci-lint`
- Shared: Uses Prettier and ESLint

### Git Workflow

1. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature
   ```

2. Make your changes and commit:

   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

3. Push and create a pull request:
   ```bash
   git push origin feature/your-feature
   ```

### Testing

1. Frontend tests:

   ```bash
   cd frontend
   pnpm test
   ```

2. Backend tests:
   ```bash
   cd backend
   go test ./...
   ```

## Debugging

### Frontend

1. Use browser developer tools
2. Check Vite/Astro development logs
3. Enable Svelte debugging in browser devtools

### Backend

1. Use Go debugging tools:

   ```bash
   # Run with delve
   dlv debug ./cmd/ogdrip

   # Or use VS Code Go debugger
   ```

2. Check logs:
   ```bash
   tail -f backend/logs/ogdrip.log
   ```

## Common Issues

### Frontend

1. **Module not found errors**

   - Run `pnpm install`
   - Clear `.astro` cache
   - Check import paths

2. **Hot reload not working**
   - Check file watchers limit
   - Restart dev server

### Backend

1. **Chrome/Chromium issues**

   - Verify CHROME_PATH in .env
   - Check Chrome installation
   - Ensure proper permissions

2. **Database errors**
   - Check DATABASE_PATH permissions
   - Verify SQLite installation
   - Check disk space

## Performance Optimization

### Frontend

1. Use dynamic imports for large components
2. Optimize images and assets
3. Enable proper caching strategies

### Backend

1. Use connection pooling
2. Implement proper caching
3. Optimize Chrome instances

## Security Best Practices

1. Keep dependencies updated
2. Use environment variables for secrets
3. Implement proper input validation
4. Follow security headers best practices

## Additional Resources

- [Astro Documentation](https://docs.astro.build)
- [Svelte Documentation](https://svelte.dev/docs)
- [Go Documentation](https://golang.org/doc/)
- [ChromeDP Documentation](https://pkg.go.dev/github.com/chromedp/chromedp)
