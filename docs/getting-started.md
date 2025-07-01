# Getting Started with OG Drip

This guide will help you get OG Drip up and running quickly, whether you're setting up for
development or production use.

## What is OG Drip?

OG Drip is a modern Open Graph image generator that creates beautiful, customizable images for
social media sharing. It consists of:

- **Frontend**: Modern web interface built with Astro + Svelte 5
- **Backend**: High-performance Go service with ChromeDP for image generation
- **Database**: SQLite for data persistence and history tracking

## Prerequisites

Before you begin, ensure you have the following installed:

### For Development

- **Node.js** >= 22.13.0
- **pnpm** >= 10.5.2 (recommended package manager)
- **Go** >= 1.24
- **Git** for version control

### For Production

- **Docker** and **Docker Compose** (recommended)
- Or a **Coolify** instance for platform deployment

## Quick Start (5 minutes)

### Option 1: Docker (Recommended for Production)

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/ogdrip.git
   cd ogdrip
   ```

2. **Start with Docker Compose**:

   ```bash
   docker-compose up -d
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8888
   - API Documentation: http://localhost:8888/docs/

### Option 2: Local Development

1. **Clone and install dependencies**:

   ```bash
   git clone https://github.com/yourusername/ogdrip.git
   cd ogdrip
   pnpm install
   ```

2. **Set up environment variables**:

   ```bash
   # Copy example environment files
   cp frontend/.env.example frontend/.env
   cp backend/.env.example backend/.env

   # Edit the files with your configuration
   nano frontend/.env
   nano backend/.env
   ```

3. **Start development servers**:

   ```bash
   pnpm dev
   ```

4. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8888

## Basic Usage

### Generating Your First Open Graph Image

1. **Open the web interface** at http://localhost:3000

2. **Enter a URL** in the input field (e.g., `https://example.com`)

3. **Click "Generate Image"** to create your Open Graph image

4. **Download the result** or copy the generated image URL

### Using the API

You can also generate images programmatically using the REST API:

```bash
curl -X POST http://localhost:8888/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "width": 1200,
    "height": 630
  }'
```

## Configuration

### Environment Variables

#### Frontend Configuration (`frontend/.env`)

```env
# Backend API URL
PUBLIC_BACKEND_URL=http://localhost:8888
BACKEND_URL=http://localhost:8888

# Optional: Analytics and monitoring
PUBLIC_SENTRY_DSN=your_sentry_dsn_here
```

#### Backend Configuration (`backend/.env`)

```env
# Server configuration
PORT=8888
HOST=0.0.0.0

# Database
DATABASE_PATH=./data/ogdrip.db

# Admin access
ADMIN_TOKEN=your_secure_admin_token_here

# Optional: External services
SENTRY_DSN=your_sentry_dsn_here
```

### Advanced Configuration

For more advanced configuration options, see:

- [Configuration Guide](configuration.md)
- [Deployment Documentation](deployment/)
- [API Documentation](api/)

## Next Steps

Now that you have OG Drip running, you might want to:

### For Developers

1. **Explore the codebase**: Check out [Development Setup](development/setup.md)
2. **Run tests**: `pnpm test`
3. **Read contributing guidelines**: [CONTRIBUTING.md](../CONTRIBUTING.md)
4. **Set up your IDE**: See [Development Setup](development/setup.md)

### For Users

1. **Customize templates**: Learn about template customization
2. **Integrate with your site**: See [API Examples](api/examples.md)
3. **Set up monitoring**: Check [Monitoring Guide](deployment/monitoring.md)
4. **Scale for production**: Read [Production Setup](deployment/production.md)

### For Administrators

1. **Set up authentication**: Configure admin access
2. **Monitor performance**: Set up logging and metrics
3. **Configure backups**: Set up database backups
4. **Security hardening**: Review [Security Guide](security/)

## Common Use Cases

### 1. Blog or Website Integration

Generate Open Graph images automatically for your blog posts or web pages.

### 2. Social Media Management

Create consistent, branded images for social media sharing.

### 3. E-commerce Product Images

Generate product preview images for social sharing.

### 4. News and Content Sites

Automatically create engaging images for articles and news stories.

## Troubleshooting

### Common Issues

**Port already in use**:

```bash
# Check what's using the port
lsof -i :8888
# Kill the process or use a different port
```

**Permission denied (Docker)**:

```bash
# Add your user to the docker group
sudo usermod -aG docker $USER
# Log out and back in
```

**Go module issues**:

```bash
cd backend
go mod tidy
go mod download
```

For more troubleshooting help, see:

- [Common Issues](troubleshooting/common-issues.md)
- [FAQ](troubleshooting/faq.md)

## Getting Help

If you need help:

1. **Check the documentation**: Browse the [docs](.) directory
2. **Search existing issues**: Look through
   [GitHub Issues](https://github.com/yourusername/ogdrip/issues)
3. **Create a new issue**: If you can't find a solution
4. **Join the community**: Participate in discussions

## What's Next?

- **API Reference**: Learn about all available endpoints in [API Documentation](api/)
- **Architecture**: Understand the system design in [Architecture Guide](architecture/)
- **Deployment**: Deploy to production with [Deployment Guides](deployment/)
- **Contributing**: Help improve OG Drip with [Contributing Guidelines](../CONTRIBUTING.md)

---

_Need help? Check our [FAQ](troubleshooting/faq.md) or
[create an issue](https://github.com/yourusername/ogdrip/issues/new)._
