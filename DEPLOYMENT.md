# Deployment Guide for Open Graph Generator

This guide explains how to deploy the Open Graph Generator application to production environments.

## Architecture Overview

The Open Graph Generator consists of two main components:

1. **Frontend** (Astro + Svelte): The user interface for creating Open Graph cards
2. **Backend** (Go): The API service that generates Open Graph images and metadata

## Environment Variables

### Frontend Environment Variables

The frontend requires the following environment variables:

| Variable             | Description                    | Default                 |
| -------------------- | ------------------------------ | ----------------------- |
| `PUBLIC_BACKEND_URL` | URL of the backend API service | `http://localhost:8888` |
| `PORT`               | Port for the frontend server   | `3000`                  |

### Backend Environment Variables

The backend requires the following environment variables:

| Variable         | Description                               | Default                 |
| ---------------- | ----------------------------------------- | ----------------------- |
| `PORT`           | Port for the backend API server           | `8888`                  |
| `BASE_URL`       | Public URL for accessing generated assets | `http://localhost:8888` |
| `OUTPUT_DIR`     | Directory for storing generated files     | `outputs`               |
| `ENABLE_CORS`    | Enable Cross-Origin Resource Sharing      | `true`                  |
| `MAX_QUEUE_SIZE` | Maximum queue size for concurrent jobs    | `10`                    |
| `CHROME_PATH`    | Optional path to Chrome executable        | System default          |
| `LOG_LEVEL`      | Logging level (debug, info, warn, error)  | `info`                  |

## Deployment Options

### Docker Compose (Recommended)

1. Copy `.env.example` to `.env` in both frontend and backend directories
2. Update the environment variables in the `.env` files
3. Run `docker-compose up -d`

### Manual Deployment

#### Frontend

1. Set environment variables:

   ```sh
   export PUBLIC_BACKEND_URL=https://api.yourdomain.com
   export PORT=3000
   ```

2. Build the frontend:

   ```sh
   cd frontend
   npm install
   npm run build
   ```

3. Start the frontend server:
   ```sh
   npm run start
   ```

#### Backend

1. Set environment variables:

   ```sh
   export PORT=8888
   export BASE_URL=https://api.yourdomain.com
   export OUTPUT_DIR=outputs
   export ENABLE_CORS=true
   ```

2. Build the backend:

   ```sh
   cd backend
   go build -o og-service
   ```

3. Start the backend service:
   ```sh
   ./og-service
   ```

## Cloud Deployment

### Vercel (Frontend)

1. Connect your GitHub repository to Vercel
2. Configure environment variables in the Vercel dashboard
3. Deploy the frontend

### Fly.io (Backend)

1. Install the Fly.io CLI
2. Create a `fly.toml` file in the backend directory
3. Configure secrets:
   ```sh
   fly secrets set BASE_URL=https://api.yourdomain.com
   fly secrets set OUTPUT_DIR=/app/outputs
   fly secrets set ENABLE_CORS=true
   ```
4. Deploy the backend:
   ```sh
   fly deploy
   ```

## Troubleshooting

- If the frontend cannot connect to the backend, check that:

  - The `PUBLIC_BACKEND_URL` is set correctly
  - CORS is enabled on the backend
  - Network access is allowed between frontend and backend

- If image generation fails, check that:
  - Chrome is installed and accessible
  - The `OUTPUT_DIR` is writable
  - The backend has sufficient memory and CPU resources
