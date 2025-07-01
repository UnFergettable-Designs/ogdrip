# Frequently Asked Questions (FAQ)

This document answers common questions about OG Drip installation, usage, and troubleshooting.

## General Questions

### What is OG Drip?

OG Drip is an Open Graph image generator that creates beautiful, customizable images for social
media sharing. It automatically generates images from web page URLs with proper metadata extraction.

### What are the system requirements?

**For Development:**

- Node.js >= 22.13.0
- pnpm >= 10.5.2
- Go >= 1.24
- 4GB RAM minimum, 8GB recommended

**For Production:**

- Docker and Docker Compose, OR
- Coolify deployment platform
- 2GB RAM minimum, 4GB recommended
- 10GB disk space for images and database

### Is OG Drip free to use?

Yes, OG Drip is open source and free to use under the MIT License. You can use it for personal and
commercial projects.

## Installation & Setup

### Why am I getting "port already in use" errors?

This usually happens when another service is using ports 3000 or 8888. To fix this:

```bash
# Check what's using the port
lsof -i :8888
lsof -i :3000

# Kill the process or change ports in your configuration
```

You can also change the ports in your environment variables:

```env
PORT=8889  # Backend port
# Frontend port can be changed in astro.config.mjs
```

### How do I fix Docker permission errors?

Add your user to the docker group:

```bash
sudo usermod -aG docker $USER
# Log out and back in for changes to take effect
```

### Why isn't pnpm working?

Make sure you have the correct version installed:

```bash
# Install pnpm if not already installed
npm install -g pnpm@10.5.2

# Or use corepack (recommended)
corepack enable
corepack use pnpm@10.5.2
```

## Usage Questions

### How do I generate an Open Graph image?

1. **Via Web Interface:**

   - Go to http://localhost:3000
   - Enter a URL in the input field
   - Click "Generate Image"
   - Download or copy the result

2. **Via API:**
   ```bash
   curl -X POST http://localhost:8888/api/generate \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com"}'
   ```

### What image sizes are supported?

The default Open Graph size is 1200x630 pixels, but you can specify custom dimensions:

```json
{
  "url": "https://example.com",
  "width": 1200,
  "height": 630
}
```

Supported ranges:

- Width: 200-2400 pixels
- Height: 200-1600 pixels

### Can I customize the image templates?

Currently, OG Drip uses a default template that extracts and displays:

- Page title
- Meta description
- Favicon or logo
- URL

Template customization is planned for future releases. See our [TODO.md](../../TODO.md) for roadmap.

### How long are generated images stored?

Generated images are stored indefinitely by default. For production use, you should implement a
cleanup strategy:

1. **Manual cleanup:**

   ```bash
   # Remove images older than 30 days
   find ./backend/outputs -name "*.png" -mtime +30 -delete
   ```

2. **Automated cleanup:** Set up a cron job or use the admin API to manage storage.

## API Questions

### How do I authenticate with the API?

Most endpoints are public, but admin endpoints require authentication:

```bash
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  http://localhost:8888/api/admin/history
```

Set your admin token in the backend environment:

```env
ADMIN_TOKEN=your_secure_random_token_here
```

### What are the API rate limits?

- **Public endpoints:** 100 requests per hour per IP
- **Admin endpoints:** 1000 requests per hour per token

Rate limit headers are included in responses:

- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`

### How do I handle API errors?

The API returns standard HTTP status codes with JSON error responses:

```json
{
  "error": true,
  "message": "Invalid URL provided",
  "code": "INVALID_URL"
}
```

Common error codes:

- `400`: Bad request (invalid parameters)
- `401`: Unauthorized (missing/invalid token)
- `429`: Rate limit exceeded
- `500`: Internal server error

## Deployment Questions

### Which deployment method should I choose?

**Docker (Recommended for most users):**

- Easy to set up and manage
- Consistent across environments
- Good for self-hosting

**Coolify (Recommended for production):**

- Automated deployments
- Built-in SSL certificates
- Easy scaling and management

**Local development:**

- Best for development and testing
- Requires manual dependency management

### How do I deploy to production?

See our deployment guides:

- [Docker Deployment](../deployment/docker.md)
- [Coolify Deployment](../deployment/coolify.md)
- [Production Setup](../deployment/production.md)

### How do I set up HTTPS?

**With Docker:** Use a reverse proxy like nginx or Traefik with SSL certificates.

**With Coolify:** SSL is automatically configured when you add a domain.

### How do I backup my data?

**Database backup:**

```bash
# Copy the SQLite database
cp backend/data/ogdrip.db backup/ogdrip-$(date +%Y%m%d).db
```

**Generated images backup:**

```bash
# Backup images directory
tar -czf backup/images-$(date +%Y%m%d).tar.gz backend/outputs/
```

## Performance Questions

### Why is image generation slow?

Several factors can affect performance:

1. **Target website speed:** Slow websites take longer to load
2. **Image size:** Larger images take more time to generate
3. **System resources:** Low memory or CPU can slow generation
4. **Network latency:** Slow internet affects page loading

**Optimization tips:**

- Use reasonable image dimensions
- Ensure adequate system resources
- Consider implementing caching for frequently requested URLs

### How can I improve performance?

1. **Increase system resources:**

   - More RAM for browser instances
   - Faster CPU for image processing
   - SSD storage for better I/O

2. **Optimize configuration:**

   ```env
   # Reduce browser timeout for faster failures
   BROWSER_TIMEOUT=15

   # Limit concurrent generations
   MAX_CONCURRENT_GENERATIONS=3
   ```

3. **Implement caching:** Cache generated images for frequently requested URLs

### How many concurrent users can OG Drip handle?

This depends on your system resources and configuration:

- **Small instance (2GB RAM):** 5-10 concurrent generations
- **Medium instance (4GB RAM):** 10-20 concurrent generations
- **Large instance (8GB+ RAM):** 20+ concurrent generations

Each browser instance uses approximately 200-400MB of RAM.

## Troubleshooting

### The generated image is blank or corrupted

This usually indicates a problem with the target website:

1. **Check the URL:** Ensure it's accessible and returns valid HTML
2. **Test manually:** Try loading the URL in your browser
3. **Check logs:** Look for error messages in the backend logs
4. **Verify resources:** Ensure adequate memory and disk space

### I'm getting "browser launch failed" errors

This typically happens in containerized environments:

1. **Docker:** Ensure proper Chrome dependencies are installed
2. **Permissions:** Browser needs proper permissions to run
3. **Memory:** Insufficient memory can cause launch failures

**For Docker, try:**

```dockerfile
# Add Chrome dependencies
RUN apt-get update && apt-get install -y \
    chromium-browser \
    --no-install-recommends
```

### The frontend can't connect to the backend

Check your configuration:

1. **URLs match:** Ensure `PUBLIC_BACKEND_URL` is correct
2. **CORS:** Verify CORS configuration allows frontend origin
3. **Firewall:** Check if ports are blocked
4. **Network:** Ensure backend is accessible from frontend

### Database errors

Common database issues:

1. **Permissions:** Ensure write permissions to data directory
2. **Disk space:** Check available disk space
3. **Corruption:** Restore from backup if database is corrupted

## Getting Help

### Where can I get more help?

1. **Documentation:** Check our comprehensive [docs](../README.md)
2. **GitHub Issues:** Search existing issues or create a new one
3. **Community:** Join discussions in GitHub Discussions
4. **Support:** For urgent issues, contact the maintainers

### How do I report a bug?

When reporting bugs, please include:

1. **Description:** Clear description of the issue
2. **Steps to reproduce:** Exact steps to trigger the bug
3. **Environment:** OS, browser, versions
4. **Logs:** Relevant error messages or logs
5. **Screenshots:** If applicable

Use our
[bug report template](https://github.com/yourusername/ogdrip/issues/new?template=bug_report.md).

### How do I request a feature?

1. **Check existing requests:** Search for similar feature requests
2. **Use the template:** Use our feature request template
3. **Provide details:** Include use case and implementation ideas
4. **Engage with community:** Participate in discussions

---

_Don't see your question here? [Create an issue](https://github.com/yourusername/ogdrip/issues/new)
or check our [troubleshooting guide](common-issues.md)._
