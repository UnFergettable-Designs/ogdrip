# API Reference

The OG Drip API provides programmatic access to Open Graph image generation functionality. This
RESTful API allows you to integrate Open Graph image generation into your applications, websites,
and workflows.

## Base URL

```
http://localhost:8888/api  # Development
https://your-domain.com/api  # Production
```

## Authentication

Most endpoints are public, but admin endpoints require authentication via the `Authorization`
header:

```http
Authorization: Bearer YOUR_ADMIN_TOKEN
```

## Quick Start

Generate your first Open Graph image:

```bash
curl -X POST http://localhost:8888/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "width": 1200,
    "height": 630
  }'
```

## API Endpoints

### Image Generation

#### `POST /api/generate`

Generate an Open Graph image from a URL.

**Request Body:**

```json
{
  "url": "https://example.com",
  "width": 1200,
  "height": 630,
  "template": "default"
}
```

**Response:**

```json
{
  "success": true,
  "image_url": "/outputs/abc123.png",
  "metadata": {
    "title": "Example Site",
    "description": "Example description",
    "width": 1200,
    "height": 630
  }
}
```

### History and Management

#### `GET /api/history`

Retrieve generation history (admin only).

#### `DELETE /api/history/{id}`

Delete a specific generation record (admin only).

### Health and Status

#### `GET /api/health`

Check API health status.

#### `GET /api/version`

Get API version information.

## Rate Limiting

- **Public endpoints**: 100 requests per hour per IP
- **Admin endpoints**: 1000 requests per hour per token

Rate limit headers are included in all responses:

- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the rate limit resets

## Error Handling

The API uses standard HTTP status codes and returns errors in JSON format:

```json
{
  "error": true,
  "message": "Invalid URL provided",
  "code": "INVALID_URL"
}
```

Common error codes:

- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## SDKs and Client Libraries

Official client libraries:

- JavaScript/TypeScript (coming soon)
- Python (coming soon)
- Go (coming soon)

## Interactive Documentation

For interactive API exploration, visit:

- Swagger UI: `/docs/`
- OpenAPI Spec: `/api/openapi.yaml`

## Further Reading

- [Authentication Guide](authentication.md)
- [Rate Limiting Details](rate-limiting.md)
- [API Examples](examples.md)
- [Error Codes Reference](error-codes.md)
