# API Examples

This document provides practical examples of using the OG Drip API for various use cases and
integration scenarios.

## Basic Usage

### Simple Image Generation

Generate a basic Open Graph image from a URL:

```bash
curl -X POST http://localhost:8888/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com"
  }'
```

**Response:**

```json
{
  "success": true,
  "image_url": "/outputs/abc123def456.png",
  "metadata": {
    "title": "Example Domain",
    "description": "This domain is for use in illustrative examples",
    "url": "https://example.com",
    "width": 1200,
    "height": 630,
    "generated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Custom Dimensions

Generate an image with specific dimensions:

```bash
curl -X POST http://localhost:8888/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "width": 1200,
    "height": 630,
    "template": "default"
  }'
```

## JavaScript/TypeScript Examples

### Using Fetch API

```typescript
interface GenerateRequest {
  url: string;
  width?: number;
  height?: number;
  template?: string;
}

interface GenerateResponse {
  success: boolean;
  image_url: string;
  metadata: {
    title: string;
    description: string;
    url: string;
    width: number;
    height: number;
    generated_at: string;
  };
}

async function generateOpenGraphImage(request: GenerateRequest): Promise<GenerateResponse> {
  const response = await fetch('http://localhost:8888/api/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
}

// Usage
try {
  const result = await generateOpenGraphImage({
    url: 'https://example.com',
    width: 1200,
    height: 630,
  });

  console.log('Generated image:', result.image_url);
  console.log('Page title:', result.metadata.title);
} catch (error) {
  console.error('Failed to generate image:', error);
}
```

### React Hook Example

```typescript
import { useState, useCallback } from 'react';

interface UseOpenGraphReturn {
  generateImage: (url: string) => Promise<void>;
  loading: boolean;
  result: GenerateResponse | null;
  error: string | null;
}

export function useOpenGraph(): UseOpenGraphReturn {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<GenerateResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const generateImage = useCallback(async (url: string) => {
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        throw new Error(`Failed to generate image: ${response.statusText}`);
      }

      const data = await response.json();
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, []);

  return { generateImage, loading, result, error };
}

// Component usage
function OpenGraphGenerator() {
  const [url, setUrl] = useState('');
  const { generateImage, loading, result, error } = useOpenGraph();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (url) {
      generateImage(url);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="url"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="Enter URL..."
        required
      />
      <button type="submit" disabled={loading}>
        {loading ? 'Generating...' : 'Generate Image'}
      </button>

      {error && <div className="error">{error}</div>}
      {result && (
        <div>
          <img src={result.image_url} alt={result.metadata.title} />
          <p>{result.metadata.title}</p>
        </div>
      )}
    </form>
  );
}
```

## Python Examples

### Basic Python Client

```python
import requests
import json
from typing import Dict, Any, Optional

class OGDripClient:
    def __init__(self, base_url: str = "http://localhost:8888"):
        self.base_url = base_url.rstrip('/')

    def generate_image(
        self,
        url: str,
        width: int = 1200,
        height: int = 630,
        template: str = "default"
    ) -> Dict[str, Any]:
        """Generate an Open Graph image from a URL."""

        payload = {
            "url": url,
            "width": width,
            "height": height,
            "template": template
        }

        response = requests.post(
            f"{self.base_url}/api/generate",
            json=payload,
            headers={"Content-Type": "application/json"}
        )

        response.raise_for_status()
        return response.json()

    def get_health(self) -> Dict[str, Any]:
        """Check API health status."""
        response = requests.get(f"{self.base_url}/api/health")
        response.raise_for_status()
        return response.json()

# Usage example
if __name__ == "__main__":
    client = OGDripClient()

    try:
        # Check if API is healthy
        health = client.get_health()
        print(f"API Status: {health}")

        # Generate image
        result = client.generate_image("https://example.com")
        print(f"Generated image: {result['image_url']}")
        print(f"Title: {result['metadata']['title']}")

    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
```

### Async Python Client

```python
import aiohttp
import asyncio
from typing import Dict, Any

class AsyncOGDripClient:
    def __init__(self, base_url: str = "http://localhost:8888"):
        self.base_url = base_url.rstrip('/')

    async def generate_image(
        self,
        url: str,
        width: int = 1200,
        height: int = 630
    ) -> Dict[str, Any]:
        """Generate an Open Graph image asynchronously."""

        payload = {
            "url": url,
            "width": width,
            "height": height
        }

        async with aiohttp.ClientSession() as session:
            async with session.post(
                f"{self.base_url}/api/generate",
                json=payload
            ) as response:
                response.raise_for_status()
                return await response.json()

    async def generate_multiple(self, urls: list[str]) -> list[Dict[str, Any]]:
        """Generate images for multiple URLs concurrently."""
        tasks = [self.generate_image(url) for url in urls]
        return await asyncio.gather(*tasks, return_exceptions=True)

# Usage example
async def main():
    client = AsyncOGDripClient()

    urls = [
        "https://example.com",
        "https://github.com",
        "https://stackoverflow.com"
    ]

    results = await client.generate_multiple(urls)

    for i, result in enumerate(results):
        if isinstance(result, Exception):
            print(f"Error for {urls[i]}: {result}")
        else:
            print(f"Generated: {result['image_url']}")

# Run async example
# asyncio.run(main())
```

## Go Examples

### Go Client Library

```go
package ogdrip

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    BaseURL    string
    HTTPClient *http.Client
}

type GenerateRequest struct {
    URL      string `json:"url"`
    Width    int    `json:"width,omitempty"`
    Height   int    `json:"height,omitempty"`
    Template string `json:"template,omitempty"`
}

type GenerateResponse struct {
    Success  bool `json:"success"`
    ImageURL string `json:"image_url"`
    Metadata struct {
        Title       string    `json:"title"`
        Description string    `json:"description"`
        URL         string    `json:"url"`
        Width       int       `json:"width"`
        Height      int       `json:"height"`
        GeneratedAt time.Time `json:"generated_at"`
    } `json:"metadata"`
}

func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) GenerateImage(req GenerateRequest) (*GenerateResponse, error) {
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    resp, err := c.HTTPClient.Post(
        c.BaseURL+"/api/generate",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    var result GenerateResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &result, nil
}

// Usage example
func main() {
    client := NewClient("http://localhost:8888")

    result, err := client.GenerateImage(GenerateRequest{
        URL:    "https://example.com",
        Width:  1200,
        Height: 630,
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated image: %s\n", result.ImageURL)
    fmt.Printf("Title: %s\n", result.Metadata.Title)
}
```

## Integration Examples

### WordPress Plugin Integration

```php
<?php
/**
 * OG Drip WordPress Plugin
 */

class OGDripPlugin {
    private $api_url;

    public function __construct($api_url = 'http://localhost:8888') {
        $this->api_url = rtrim($api_url, '/');
        add_action('save_post', array($this, 'generate_og_image'));
    }

    public function generate_og_image($post_id) {
        // Only generate for published posts
        if (get_post_status($post_id) !== 'publish') {
            return;
        }

        $post_url = get_permalink($post_id);
        $result = $this->call_api($post_url);

        if ($result && $result['success']) {
            // Save image URL as post meta
            update_post_meta($post_id, '_og_image_url', $result['image_url']);
            update_post_meta($post_id, '_og_metadata', $result['metadata']);
        }
    }

    private function call_api($url) {
        $data = json_encode(array(
            'url' => $url,
            'width' => 1200,
            'height' => 630
        ));

        $options = array(
            'http' => array(
                'header' => "Content-type: application/json\r\n",
                'method' => 'POST',
                'content' => $data
            )
        );

        $context = stream_context_create($options);
        $result = file_get_contents($this->api_url . '/api/generate', false, $context);

        return json_decode($result, true);
    }

    // Add OG tags to head
    public function add_og_tags() {
        if (is_single()) {
            global $post;
            $og_image = get_post_meta($post->ID, '_og_image_url', true);

            if ($og_image) {
                echo '<meta property="og:image" content="' . esc_url($og_image) . '" />';
                echo '<meta property="og:image:width" content="1200" />';
                echo '<meta property="og:image:height" content="630" />';
            }
        }
    }
}

// Initialize plugin
add_action('init', function() {
    $ogdrip = new OGDripPlugin();
    add_action('wp_head', array($ogdrip, 'add_og_tags'));
});
?>
```

### Next.js API Route

```typescript
// pages/api/generate-og.ts
import { NextApiRequest, NextApiResponse } from 'next';

interface OGGenerateRequest {
  url: string;
  width?: number;
  height?: number;
}

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { url, width = 1200, height = 630 }: OGGenerateRequest = req.body;

  if (!url) {
    return res.status(400).json({ error: 'URL is required' });
  }

  try {
    const response = await fetch(`${process.env.OGDRIP_API_URL}/api/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ url, width, height }),
    });

    if (!response.ok) {
      throw new Error(`API responded with status ${response.status}`);
    }

    const result = await response.json();
    res.status(200).json(result);
  } catch (error) {
    console.error('OG generation failed:', error);
    res.status(500).json({ error: 'Failed to generate Open Graph image' });
  }
}
```

### Express.js Middleware

```javascript
const express = require('express');
const axios = require('axios');

// Middleware for automatic OG image generation
function ogImageMiddleware(options = {}) {
  const {
    apiUrl = 'http://localhost:8888',
    cacheDuration = 3600000, // 1 hour
    skipRoutes = ['/api', '/admin'],
  } = options;

  const cache = new Map();

  return async (req, res, next) => {
    // Skip API routes and admin pages
    if (skipRoutes.some(route => req.path.startsWith(route))) {
      return next();
    }

    const fullUrl = `${req.protocol}://${req.get('host')}${req.originalUrl}`;
    const cacheKey = `og:${fullUrl}`;

    // Check cache first
    const cached = cache.get(cacheKey);
    if (cached && Date.now() - cached.timestamp < cacheDuration) {
      req.ogImage = cached.data;
      return next();
    }

    try {
      // Generate OG image
      const response = await axios.post(`${apiUrl}/api/generate`, {
        url: fullUrl,
        width: 1200,
        height: 630,
      });

      const ogData = response.data;

      // Cache the result
      cache.set(cacheKey, {
        data: ogData,
        timestamp: Date.now(),
      });

      req.ogImage = ogData;
    } catch (error) {
      console.error('Failed to generate OG image:', error.message);
      req.ogImage = null;
    }

    next();
  };
}

// Usage
const app = express();
app.use(
  ogImageMiddleware({
    apiUrl: process.env.OGDRIP_API_URL,
    cacheDuration: 3600000, // 1 hour
  })
);

app.get('*', (req, res) => {
  const ogImage = req.ogImage;

  res.send(`
    <!DOCTYPE html>
    <html>
      <head>
        <title>My Site</title>
        ${
          ogImage
            ? `
          <meta property="og:image" content="${ogImage.image_url}" />
          <meta property="og:title" content="${ogImage.metadata.title}" />
          <meta property="og:description" content="${ogImage.metadata.description}" />
        `
            : ''
        }
      </head>
      <body>
        <h1>Welcome to My Site</h1>
      </body>
    </html>
  `);
});
```

## Error Handling Examples

### Comprehensive Error Handling

```typescript
interface APIError {
  error: boolean;
  message: string;
  code: string;
  details?: any;
}

class OGDripError extends Error {
  constructor(
    message: string,
    public code: string,
    public status: number,
    public details?: any
  ) {
    super(message);
    this.name = 'OGDripError';
  }
}

async function generateWithRetry(
  url: string,
  maxRetries: number = 3,
  delay: number = 1000
): Promise<GenerateResponse> {
  let lastError: Error;

  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      const response = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        const errorData: APIError = await response.json();
        throw new OGDripError(
          errorData.message,
          errorData.code,
          response.status,
          errorData.details
        );
      }

      return await response.json();
    } catch (error) {
      lastError = error as Error;

      // Don't retry on client errors (4xx)
      if (error instanceof OGDripError && error.status < 500) {
        throw error;
      }

      if (attempt < maxRetries) {
        console.warn(`Attempt ${attempt} failed, retrying in ${delay}ms...`);
        await new Promise(resolve => setTimeout(resolve, delay));
        delay *= 2; // Exponential backoff
      }
    }
  }

  throw new Error(`Failed after ${maxRetries} attempts: ${lastError.message}`);
}

// Usage with error handling
try {
  const result = await generateWithRetry('https://example.com');
  console.log('Success:', result.image_url);
} catch (error) {
  if (error instanceof OGDripError) {
    switch (error.code) {
      case 'INVALID_URL':
        console.error('Please provide a valid URL');
        break;
      case 'RATE_LIMIT_EXCEEDED':
        console.error('Too many requests, please try again later');
        break;
      case 'GENERATION_TIMEOUT':
        console.error('Image generation timed out');
        break;
      default:
        console.error('API Error:', error.message);
    }
  } else {
    console.error('Unexpected error:', error.message);
  }
}
```

## Testing Examples

### API Testing with Jest

```typescript
import { describe, it, expect, beforeAll, afterAll } from '@jest/globals';

describe('OG Drip API', () => {
  const API_URL = 'http://localhost:8888';

  beforeAll(async () => {
    // Wait for API to be ready
    await waitForAPI(API_URL);
  });

  it('should generate image for valid URL', async () => {
    const response = await fetch(`${API_URL}/api/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        url: 'https://example.com',
        width: 1200,
        height: 630,
      }),
    });

    expect(response.status).toBe(200);

    const result = await response.json();
    expect(result.success).toBe(true);
    expect(result.image_url).toMatch(/^\/outputs\/.*\.png$/);
    expect(result.metadata.title).toBeTruthy();
  });

  it('should reject invalid URLs', async () => {
    const response = await fetch(`${API_URL}/api/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        url: 'not-a-valid-url',
      }),
    });

    expect(response.status).toBe(400);

    const result = await response.json();
    expect(result.error).toBe(true);
    expect(result.code).toBe('INVALID_URL');
  });

  it('should handle rate limiting', async () => {
    // Make multiple rapid requests
    const promises = Array(10)
      .fill(null)
      .map(() =>
        fetch(`${API_URL}/api/generate`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ url: 'https://example.com' }),
        })
      );

    const responses = await Promise.all(promises);
    const rateLimited = responses.some(r => r.status === 429);

    // Should eventually hit rate limit
    expect(rateLimited).toBe(true);
  });
});

async function waitForAPI(url: string, timeout: number = 30000): Promise<void> {
  const start = Date.now();

  while (Date.now() - start < timeout) {
    try {
      const response = await fetch(`${url}/api/health`);
      if (response.ok) return;
    } catch (error) {
      // API not ready yet
    }

    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  throw new Error('API did not become ready within timeout');
}
```

---

_For more examples and use cases, see our [API Reference](README.md) or
[create an issue](https://github.com/yourusername/ogdrip/issues/new) with your specific integration
question._
