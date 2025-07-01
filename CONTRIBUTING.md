# Contributing to OG Drip

Thank you for considering contributing to OG Drip! This document provides guidelines and information
for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Security Guidelines](#security-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow:

- Be respectful and inclusive
- Focus on constructive feedback
- Help create a welcoming environment for all contributors
- Report any unacceptable behavior to the project maintainers

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a feature branch for your changes
5. Make your changes and test them
6. Submit a pull request

## Development Setup

### Prerequisites

- Node.js >= 22.13.0
- pnpm >= 10.5.2
- Go >= 1.24
- Docker (for containerized development)

### Initial Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/ogdrip.git
cd ogdrip

# Install dependencies
pnpm install

# Set up environment files
cp frontend/.env.example frontend/.env
cp backend/.env.example backend/.env

# Start development servers
pnpm dev
```

### Project Structure

This is a monorepo with the following components:

- `frontend/`: Astro + Svelte 5 frontend application
- `backend/`: Go backend service with ChromeDP
- `shared/`: Shared TypeScript types and utilities
- `docs/`: Project documentation

## Contributing Process

1. **Check existing issues** - Look for existing issues or create a new one
2. **Discuss major changes** - For significant changes, discuss with maintainers first
3. **Create a branch** - Use descriptive branch names (e.g., `feature/add-new-template`,
   `fix/cors-issue`)
4. **Make changes** - Follow coding standards and write tests
5. **Test thoroughly** - Ensure all tests pass and add new tests for your changes
6. **Update documentation** - Update relevant documentation
7. **Submit PR** - Create a pull request with a clear description

## Coding Standards

### General Guidelines

- Follow the existing code style and patterns
- Write self-documenting code with clear variable and function names
- Keep functions small and focused (10-30 lines when possible)
- Use TypeScript for all frontend code
- Follow Go conventions for backend code

### Frontend Standards (Astro + Svelte)

- Use Svelte 5 runes syntax (`$state`, `$derived`, etc.)
- Include `lang="ts"` in script tags
- Use semantic HTML elements
- Follow accessibility guidelines (WCAG 2.2 AA)
- Use REM for sizing and spacing
- Use HSLA for colors
- Import Svelte components with `.svelte` extension

### Backend Standards (Go)

- Use proper error handling with context
- Implement timeouts for ChromeDP operations
- Use prepared statements for database operations
- Follow Go naming conventions
- Include proper logging
- Implement graceful shutdown

### CSS Standards

- Use CSS custom properties for design tokens
- Follow BEM naming convention when applicable
- Use semantic class names
- Maintain consistent spacing scale
- Ensure minimum 4.5:1 contrast ratio

## Testing Requirements

### Test Coverage

- Maintain 80-90% code coverage for critical paths
- Test both happy and error paths
- Include edge cases and boundary conditions

### Frontend Testing

- Unit tests with Vitest
- Component testing with Testing Library
- E2E tests with Playwright
- Accessibility testing with axe-core

### Backend Testing

- Unit tests with Go's built-in testing package
- Integration tests for API endpoints
- Database testing with test fixtures
- Performance testing for critical paths

### Running Tests

```bash
# Run all tests
pnpm test

# Run tests with coverage
pnpm test:coverage

# Run specific test suites
pnpm test --filter=@ogdrip/frontend
pnpm test --filter=@ogdrip/backend
```

## Security Guidelines

### Security Requirements

- Never commit secrets, tokens, or credentials
- Use environment variables for sensitive data
- Validate and sanitize all user inputs
- Use parameterized queries for database operations
- Implement proper authentication and authorization
- Follow the principle of least privilege

### Security Testing

- Run security audits on dependencies
- Test for common vulnerabilities (XSS, CSRF, SQL injection)
- Validate input sanitization
- Test authentication and authorization flows

### Reporting Security Issues

Please report security vulnerabilities privately to the maintainers. Do not create public issues for
security problems.

## Documentation

### Documentation Requirements

- Update README.md for significant changes
- Document all public APIs
- Include code examples in documentation
- Update CHANGELOG.md for all changes
- Add JSDoc comments for complex functions

### Documentation Standards

- Use clear, concise language
- Include practical examples
- Keep documentation up to date with code changes
- Follow markdown best practices

## Pull Request Process

### Before Submitting

1. Ensure all tests pass
2. Run linting and formatting tools
3. Update documentation
4. Add entry to CHANGELOG.md
5. Rebase on the latest main branch

### PR Description Template

```markdown
## Description

Brief description of changes

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing

- [ ] Tests pass locally
- [ ] Added tests for new features
- [ ] Manual testing completed

## Checklist

- [ ] Code follows project standards
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
```

### Review Process

1. Automated checks must pass
2. At least one maintainer review required
3. Address all review feedback
4. Maintain clean commit history
5. Squash commits if requested

## Issue Reporting

### Bug Reports

Include the following information:

- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, browser, versions)
- Screenshots or logs if applicable

### Feature Requests

Include the following information:

- Clear description of the feature
- Use case and motivation
- Proposed implementation approach
- Alternatives considered

### Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Improvements or additions to docs
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention is needed
- `security`: Security-related issue

## Development Tools

### Required Tools

- **ESLint**: Code linting with accessibility plugins
- **Prettier**: Code formatting
- **EditorConfig**: Consistent coding styles
- **Husky**: Git hooks for pre-commit validation
- **Turbo**: Monorepo task runner

### Recommended IDE Setup

- VS Code with recommended extensions
- Svelte extension for Svelte support
- Go extension for Go development
- ESLint and Prettier extensions

## Getting Help

- Check existing documentation
- Search existing issues
- Join discussions in GitHub Discussions
- Reach out to maintainers for major questions

## License

By contributing to OG Drip, you agree that your contributions will be licensed under the MIT
License.
