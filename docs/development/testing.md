# Testing Guide

This document outlines the testing strategy, tools, and practices for OG Drip. We follow a
comprehensive testing approach to ensure code quality, reliability, and maintainability.

## Testing Philosophy

Our testing strategy follows these principles:

1. **Test-Driven Development (TDD)** - Write tests before implementing features when practical
2. **Comprehensive Coverage** - Aim for 80-90% code coverage on critical paths
3. **Fast Feedback** - Tests should run quickly and provide clear feedback
4. **Realistic Testing** - Test real scenarios, not just happy paths
5. **Maintainable Tests** - Tests should be easy to understand and maintain

## Testing Stack

### Frontend Testing

- **Vitest** - Unit testing framework
- **Testing Library** - Component testing utilities
- **Playwright** - End-to-end testing
- **axe-core** - Accessibility testing
- **MSW** - API mocking

### Backend Testing

- **Go testing** - Built-in testing package
- **testify** - Testing assertions and mocking
- **SQLite** - In-memory database for tests
- **httptest** - HTTP testing utilities

## Test Categories

### 1. Unit Tests

#### Frontend Unit Tests

Located in: `frontend/src/test/`

```bash
# Run frontend unit tests
cd frontend
pnpm test

# Run in watch mode
pnpm test:watch

# Run with coverage
pnpm test:coverage
```

**Example Svelte Component Test:**

```typescript
import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import OpenGraphForm from '../components/OpenGraphForm.svelte';

describe('OpenGraphForm', () => {
  it('renders form with URL input', () => {
    render(OpenGraphForm);

    const urlInput = screen.getByLabelText('URL');
    expect(urlInput).toBeInTheDocument();
    expect(urlInput).toHaveAttribute('type', 'url');
  });

  it('validates URL input', async () => {
    const { component } = render(OpenGraphForm);

    const urlInput = screen.getByLabelText('URL');
    await fireEvent.input(urlInput, { target: { value: 'invalid-url' } });

    expect(screen.getByText('Please enter a valid URL')).toBeInTheDocument();
  });
});
```

#### Backend Unit Tests

Located in: `backend/*_test.go`

```bash
# Run backend unit tests
cd backend
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestGenerateImage ./...
```

**Example Go Unit Test:**

```go
func TestValidateURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        wantErr bool
    }{
        {"valid HTTP URL", "http://example.com", false},
        {"valid HTTPS URL", "https://example.com", false},
        {"invalid scheme", "ftp://example.com", true},
        {"empty URL", "", true},
        {"malformed URL", "not-a-url", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 2. Integration Tests

#### API Integration Tests

```bash
# Run integration tests
cd backend
go test -tags=integration ./...
```

**Example Integration Test:**

```go
//go:build integration
// +build integration

func TestGenerateImageAPI(t *testing.T) {
    // Setup test server
    server := setupTestServer()
    defer server.Close()

    // Test data
    payload := map[string]interface{}{
        "url":    "https://example.com",
        "width":  1200,
        "height": 630,
    }

    // Make request
    resp, err := makeAPIRequest("POST", "/api/generate", payload)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)

    // Verify response
    var result GenerateResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err)
    assert.True(t, result.Success)
    assert.NotEmpty(t, result.ImageURL)
}
```

### 3. End-to-End Tests

#### Playwright E2E Tests

Located in: `frontend/tests/e2e/`

```bash
# Install Playwright
cd frontend
pnpm exec playwright install

# Run E2E tests
pnpm test:e2e

# Run in headed mode
pnpm test:e2e --headed

# Run specific test
pnpm test:e2e tests/image-generation.spec.ts
```

**Example E2E Test:**

```typescript
import { test, expect } from '@playwright/test';

test.describe('Image Generation', () => {
  test('should generate image from URL', async ({ page }) => {
    // Navigate to app
    await page.goto('/');

    // Fill form
    await page.fill('[data-testid="url-input"]', 'https://example.com');
    await page.click('[data-testid="generate-button"]');

    // Wait for generation
    await expect(page.locator('[data-testid="loading"]')).toBeVisible();
    await expect(page.locator('[data-testid="loading"]')).toBeHidden();

    // Verify result
    await expect(page.locator('[data-testid="generated-image"]')).toBeVisible();
    await expect(page.locator('[data-testid="download-link"]')).toBeVisible();
  });

  test('should handle invalid URLs', async ({ page }) => {
    await page.goto('/');

    await page.fill('[data-testid="url-input"]', 'invalid-url');
    await page.click('[data-testid="generate-button"]');

    await expect(page.locator('[data-testid="error-message"]')).toContainText(
      'Please enter a valid URL'
    );
  });
});
```

### 4. Accessibility Tests

#### Automated Accessibility Testing

```typescript
import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Accessibility', () => {
  test('should not have accessibility violations', async ({ page }) => {
    await page.goto('/');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('should be keyboard navigable', async ({ page }) => {
    await page.goto('/');

    // Tab through interactive elements
    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="url-input"]')).toBeFocused();

    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="generate-button"]')).toBeFocused();
  });
});
```

### 5. Performance Tests

#### Load Testing

```bash
# Install k6 for load testing
brew install k6  # macOS
# or
sudo apt-get install k6  # Ubuntu

# Run load test
k6 run tests/load/api-load-test.js
```

**Example Load Test:**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10, // 10 virtual users
  duration: '30s',
};

export default function () {
  const payload = JSON.stringify({
    url: 'https://example.com',
    width: 1200,
    height: 630,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const response = http.post('http://localhost:8888/api/generate', payload, params);

  check(response, {
    'status is 200': r => r.status === 200,
    'response time < 5000ms': r => r.timings.duration < 5000,
  });

  sleep(1);
}
```

## Test Configuration

### Frontend Test Configuration

#### Vitest Config (`vitest.config.ts`)

```typescript
import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte({ hot: !process.env.VITEST })],
  test: {
    environment: 'jsdom',
    setupFiles: ['src/test/setup.ts'],
    coverage: {
      reporter: ['text', 'json', 'html'],
      exclude: ['node_modules/', 'src/test/', '**/*.d.ts', '**/*.config.*'],
    },
  },
});
```

#### Test Setup (`src/test/setup.ts`)

```typescript
import '@testing-library/jest-dom';
import { beforeAll, afterAll, afterEach } from 'vitest';
import { server } from './mocks/server';

// Setup MSW
beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
```

### Backend Test Configuration

#### Test Database Setup

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }

    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }

    return db
}
```

## Testing Best Practices

### General Practices

1. **AAA Pattern** - Arrange, Act, Assert
2. **Descriptive Names** - Test names should describe what they test
3. **Single Responsibility** - Each test should test one thing
4. **Independent Tests** - Tests should not depend on each other
5. **Clean Up** - Always clean up resources after tests

### Frontend Testing Best Practices

1. **Test User Behavior** - Focus on what users do, not implementation details
2. **Use Test IDs** - Add `data-testid` attributes for reliable element selection
3. **Mock External Dependencies** - Use MSW for API mocking
4. **Test Accessibility** - Include accessibility tests in your test suite

### Backend Testing Best Practices

1. **Use Table-Driven Tests** - Test multiple scenarios efficiently
2. **Test Error Cases** - Don't just test happy paths
3. **Use Test Doubles** - Mock external services and dependencies
4. **Test Concurrency** - Test concurrent operations where applicable

## Continuous Integration

### GitHub Actions Workflow

```yaml
name: Tests

on: [push, pull_request]

jobs:
  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '22'
      - run: corepack enable
      - run: pnpm install
      - run: pnpm test:frontend
      - run: pnpm test:e2e

  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.24'
      - run: cd backend && go test -cover ./...
      - run: cd backend && go test -tags=integration ./...
```

## Test Coverage

### Coverage Goals

- **Critical paths**: 90%+ coverage
- **Business logic**: 85%+ coverage
- **UI components**: 80%+ coverage
- **Utilities**: 95%+ coverage

### Checking Coverage

```bash
# Frontend coverage
cd frontend
pnpm test:coverage

# Backend coverage
cd backend
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Testing Tools and Utilities

### Custom Test Utilities

#### Frontend Test Helpers

```typescript
// src/test/helpers.ts
export function renderWithProviders(component: any, options = {}) {
  return render(component, {
    ...options,
    // Add any global providers here
  });
}

export function createMockApiResponse(data: any) {
  return {
    success: true,
    data,
    ...data,
  };
}
```

#### Backend Test Helpers

```go
// test_helpers.go
func CreateTestServer() *httptest.Server {
    handler := setupRoutes()
    return httptest.NewServer(handler)
}

func MakeAPIRequest(method, path string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader
    if body != nil {
        jsonBody, _ := json.Marshal(body)
        reqBody = bytes.NewBuffer(jsonBody)
    }

    req, err := http.NewRequest(method, testServer.URL+path, reqBody)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    return http.DefaultClient.Do(req)
}
```

## Debugging Tests

### Frontend Test Debugging

```bash
# Run tests in debug mode
cd frontend
pnpm test --inspect-brk

# Run specific test file
pnpm test OpenGraphForm.test.ts

# Run with verbose output
pnpm test --reporter=verbose
```

### Backend Test Debugging

```bash
# Run tests with verbose output
cd backend
go test -v ./...

# Run specific test with debugging
go test -v -run TestSpecificFunction

# Use delve for debugging
dlv test -- -test.run TestSpecificFunction
```

## Common Testing Patterns

### Testing Async Operations

```typescript
test('should handle async image generation', async () => {
  const mockGenerate = vi.fn().mockResolvedValue({
    success: true,
    imageUrl: '/outputs/test.png',
  });

  render(OpenGraphForm, {
    props: { onGenerate: mockGenerate },
  });

  await fireEvent.click(screen.getByRole('button', { name: /generate/i }));

  await waitFor(() => {
    expect(screen.getByText('Image generated successfully')).toBeInTheDocument();
  });
});
```

### Testing Error Handling

```go
func TestHandleError(t *testing.T) {
    tests := []struct {
        name           string
        input          string
        expectedStatus int
        expectedError  string
    }{
        {
            name:           "invalid URL",
            input:          "not-a-url",
            expectedStatus: http.StatusBadRequest,
            expectedError:  "invalid URL format",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := makeRequest(tt.input)
            require.NoError(t, err)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
        })
    }
}
```

## Troubleshooting Tests

### Common Issues

**Tests timing out:**

- Increase timeout values
- Check for infinite loops or blocking operations
- Ensure proper cleanup of resources

**Flaky tests:**

- Add proper waits for async operations
- Use deterministic test data
- Avoid relying on timing

**Memory leaks in tests:**

- Clean up event listeners
- Close database connections
- Clear timers and intervals

---

_For more testing help, see our [troubleshooting guide](../troubleshooting/common-issues.md) or
[create an issue](https://github.com/yourusername/ogdrip/issues/new)._
