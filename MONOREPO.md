# OGDrip Monorepo Structure

This project uses a monorepo structure managed with [pnpm](https://pnpm.io/) and [Turborepo](https://turbo.build/).

## Directory Structure

- `frontend/`: Astro + Svelte 5 frontend application
- `backend/`: Go backend for Open Graph image generation
- `shared/`: TypeScript types and utilities shared between packages

## Setup

```bash
# Install dependencies
pnpm install

# Setup directories and files for frontend
pnpm run setup

# Setup directories and files for backend
pnpm run setup:backend
```

## Available Commands

### Root Commands

These commands can be run from the root directory:

```bash
# Build all packages
pnpm build

# Start development mode (runs both frontend and backend in parallel)
pnpm dev

# Start frontend development server only
pnpm dev:frontend

# Start backend development server only
pnpm dev:backend

# Start both frontend and backend in development mode (explicitly)
pnpm dev:all

# Lint all packages
pnpm lint

# Clean all packages
pnpm clean

# Format code with Prettier
pnpm format

# Preview the frontend build
pnpm preview

# Start the production build
pnpm start

# Docker commands
pnpm docker:build    # Build Docker containers
pnpm docker:start    # Start Docker containers
pnpm docker:clean    # Rebuild Docker containers from scratch
```

### Package-specific Commands

Each package can also be worked with independently:

```bash
# Run commands in a specific package
pnpm --filter @ogdrip/frontend build
pnpm --filter @ogdrip/backend build
pnpm --filter @ogdrip/shared build
```

## Package Dependencies

- **frontend**: Depends on shared for types
- **backend**: Standalone Go service (no Node.js dependencies)
- **shared**: No dependencies

## Configuration Files

- `pnpm-workspace.yaml`: Defines the workspace packages
- `turbo.json`: Defines the task relationships and dependencies
- `package.json`: Root package with scripts for the monorepo

## Development Workflow

1. Run `pnpm install` to install all dependencies
2. Run `pnpm setup` and `pnpm setup:backend` to create necessary directories
3. Run `pnpm dev` to start all services in parallel, or use:
   - `pnpm dev:frontend` to start just the frontend
   - `pnpm dev:backend` to start just the backend
   - `pnpm dev:all` to explicitly start both
4. Make changes to code
5. Use `pnpm build` to build all packages for production

## Notes

- The backend is built with Go and runs independently
- The frontend is built with Astro and Svelte 5
- Shared types ensure consistency between packages
- Building the backend requires Go to be installed on your system
- Persistent tasks in Turborepo (like dev, dev:frontend, dev:backend) run servers that stay alive
