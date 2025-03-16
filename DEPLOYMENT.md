# Deployment Guide

This guide explains how to deploy the Open Graph Generator service using Nixpacks.

## Prerequisites

- Git
- Node.js 20.x or later
- pnpm 8.x or later
- Go 1.24 or later
- SQLite 3.x

## Deployment Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/ogdrip.git
   cd ogdrip
   ```

2. Install dependencies:

   ```bash
   pnpm install
   ```

3. Build the application:

   ```bash
   pnpm build
   ```

4. Set up environment variables:

   - Copy `.env.example` to `.env` in both frontend and backend directories
   - Update the values according to your environment
   - Required variables:

     ```
     # Frontend
     PUBLIC_BACKEND_URL=http://localhost:8888
     SENTRY_DSN=your-sentry-dsn (optional)

     # Backend
     PORT=8888
     BASE_URL=http://localhost:8888
     OUTPUT_DIR=./outputs
     ENABLE_CORS=true
     MAX_QUEUE_SIZE=10
     ```

5. Deploy using Nixpacks:
   ```bash
   nixpacks build . --name ogdrip
   ```

## Environment Variables

### Frontend Variables

- `PUBLIC_BACKEND_URL`: URL where the backend service is accessible
- `SENTRY_DSN`: Sentry DSN for error tracking (optional)
- `NODE_ENV`: Set to "production" for production builds
- `PORT`: Port for the frontend service (default: 3000)

### Backend Variables

- `PORT`: Port for the backend service (default: 8888)
- `BASE_URL`: Public URL where the backend service is accessible
- `OUTPUT_DIR`: Directory for storing generated images
- `ENABLE_CORS`: Enable CORS for cross-origin requests
- `MAX_QUEUE_SIZE`: Maximum number of concurrent image generations
- `SENTRY_DSN`: Sentry DSN for error tracking (optional)

## Health Checks

The service provides health check endpoints:

- Frontend: `GET /health`
- Backend: `GET /api/health`

## Monitoring

1. Application logs are available through your platform's logging interface

2. Error tracking is available through Sentry if configured

3. Key metrics to monitor:
   - HTTP response times
   - Error rates
   - Image generation queue size
   - Disk usage for image storage

## Troubleshooting

1. If the frontend can't connect to the backend:

   - Check that `PUBLIC_BACKEND_URL` is set correctly
   - Verify the backend service is running
   - Check network/firewall settings

2. If image generation fails:

   - Verify Chromium is installed and accessible
   - Check the `outputs` directory permissions
   - Review backend logs for specific errors

3. For performance issues:
   - Monitor the image generation queue size
   - Check disk space in the outputs directory
   - Review resource usage (CPU, memory)

## Backup and Recovery

1. Important data to backup:

   - SQLite database in the backend's data directory
   - Generated images in the outputs directory
   - Environment configuration files

2. Recovery steps:
   - Restore the SQLite database file
   - Restore the outputs directory
   - Verify environment variables
   - Restart the services
