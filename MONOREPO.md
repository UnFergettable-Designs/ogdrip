# Monorepo Structure

This repository uses pnpm workspaces to manage multiple packages in a monorepo structure.

## Project Structure

```
ogdrip/
├── frontend/           # Astro + Svelte frontend
│   ├── src/           # Frontend source code
│   ├── public/        # Static assets
│   └── package.json   # Frontend package config
├── backend/           # Go backend service
│   ├── cmd/          # Command line tools
│   ├── internal/     # Internal packages
│   └── go.mod        # Go module file
├── shared/           # Shared TypeScript types
│   └── src/         # Shared code
├── package.json     # Root package.json
└── pnpm-workspace.yaml # Workspace configuration
```

## Workspace Commands

```bash
# Install all dependencies
pnpm install

# Development
pnpm dev:all         # Start all services
pnpm dev:frontend    # Start frontend only
pnpm dev:backend     # Start backend only

# Building
pnpm build           # Build all packages
pnpm build:frontend  # Build frontend only
pnpm build:backend   # Build backend only

# Testing
pnpm test           # Run all tests
pnpm test:frontend  # Run frontend tests
pnpm test:backend   # Run backend tests

# Linting
pnpm lint          # Lint all code
pnpm lint:fix      # Fix linting issues

# Clean
pnpm clean        # Clean all build artifacts
```

## Package Management

### Adding Dependencies

1. Workspace-specific dependencies:

   ```bash
   # Frontend dependencies
   cd frontend
   pnpm add package-name

   # Shared dependencies
   cd shared
   pnpm add package-name
   ```

2. Root dependencies:
   ```bash
   pnpm add -w package-name
   ```

### Updating Dependencies

```bash
# Update all dependencies
pnpm update

# Update specific package
pnpm update package-name
```

## Development Workflow

1. Start development servers:

   ```bash
   pnpm dev:all
   ```

2. Make changes in respective packages:

   - Frontend: `frontend/`
   - Backend: `backend/`
   - Shared types: `shared/`

3. Build for production:
   ```bash
   pnpm build
   ```

## Shared Code

The `shared` package contains:

- TypeScript types
- Utility functions
- Constants
- API interfaces

Import shared code in frontend:

```typescript
import { type OpenGraphRequest } from '@ogdrip/shared';
```

## Testing

1. Run all tests:

   ```bash
   pnpm test
   ```

2. Test specific package:
   ```bash
   pnpm test:frontend
   pnpm test:backend
   ```

## Best Practices

1. Package Organization:

   - Keep packages focused and minimal
   - Share code through `shared` package
   - Maintain clear package boundaries

2. Dependencies:

   - Use workspace references when possible
   - Keep dependencies up to date
   - Avoid duplicate dependencies

3. TypeScript:

   - Share types through `shared` package
   - Use strict TypeScript settings
   - Maintain type consistency

4. Testing:
   - Write tests for shared code
   - Test package integration
   - Maintain test coverage

## Troubleshooting

1. Dependency Issues:

   ```bash
   pnpm clean
   rm -rf node_modules
   pnpm install
   ```

2. Build Issues:

   ```bash
   pnpm clean
   pnpm build
   ```

3. TypeScript Errors:
   - Check `shared` types
   - Update TypeScript version
   - Clear TypeScript cache
