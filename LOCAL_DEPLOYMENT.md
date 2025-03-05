# Local Deployment Guide for Open Graph Generator

This guide will help you set up and run the Open Graph Generator locally for development or personal use.

## Prerequisites

- [Node.js](https://nodejs.org/) (v16 or later)
- [PNPM](https://pnpm.io/) (v8 or later)
- [Go](https://golang.org/) (v1.23 or later)
- [Git](https://git-scm.com/)
- [Chrome/Chromium](https://www.google.com/chrome/) (required for the backend's headless browser functionality)

## Quick Start

The quickest way to get started is to use Docker Compose for local development:

```bash
git clone https://github.com/your-username/ogdrip.git
cd ogdrip
docker compose up
```

This will start both the frontend and backend services. The frontend will be available at http://localhost:3000 and the backend at http://localhost:8888.

### Docker Hub Authentication Issues

If you encounter authentication issues when pulling Docker images, you have a few options:

1. **Log in to Docker Hub**:

   ```bash
   docker login
   ```

   This will prompt you to enter your Docker Hub username and password.

2. **Use the manual setup instead**:
   If you prefer not to authenticate with Docker Hub, follow the manual setup instructions below.

3. **Create local images**:
   You can build the images locally without pulling from Docker Hub by first installing Node.js and Go, then running:

   ```bash
   # For the frontend
   cd frontend
   pnpm install
   pnpm dev

   # For the backend (in a separate terminal)
   cd backend
   go run .
   ```

## Manual Setup

If you prefer to run the services directly on your machine without Docker, follow these steps:

### Backend Setup

1. Clone the repository and navigate to the backend directory:

```bash
git clone https://github.com/your-username/ogdrip.git
cd ogdrip/backend
```

2. Create a `.env` file for local development:

```bash
cp .env.example .env
```

3. Update the `.env` file with your local settings:

```
PORT=8888
BASE_URL=http://localhost:8888
ENABLE_CORS=true
OUTPUT_DIR=./outputs
```

4. Run the backend:

```bash
go run .
```

The backend API will be accessible at http://localhost:8888.

### Frontend Setup

1. Open a new terminal and navigate to the frontend directory:

```bash
cd ogdrip/frontend
```

2. Create a `.env` file for local development:

```bash
cp .env.example .env
```

3. Update the `.env` file with your local settings:

```
BACKEND_URL=http://localhost:8888
PUBLIC_BACKEND_URL=http://localhost:8888
```

4. Install dependencies and start the development server:

```bash
pnpm install
pnpm dev
```

The frontend will be accessible at http://localhost:3000.

## Configuration Options

### Backend Configuration

| Environment Variable | Description                               | Default                 |
| -------------------- | ----------------------------------------- | ----------------------- |
| `PORT`               | Port for the backend API server           | `8888`                  |
| `BASE_URL`           | Public URL for accessing generated assets | `http://localhost:8888` |
| `OUTPUT_DIR`         | Directory for storing generated files     | `./outputs`             |
| `ENABLE_CORS`        | Enable Cross-Origin Resource Sharing      | `true`                  |
| `MAX_QUEUE_SIZE`     | Maximum queue size for concurrent jobs    | `10`                    |
| `CHROME_PATH`        | Optional path to Chrome executable        | System default          |
| `ADMIN_TOKEN`        | Token for accessing admin features        | `admin`                 |
| `DB_PATH`            | Path for the SQLite database              | `./data/generations.db` |

### Frontend Configuration

| Environment Variable | Description                   | Default                 |
| -------------------- | ----------------------------- | ----------------------- |
| `BACKEND_URL`        | Backend API URL (server-side) | `http://localhost:8888` |
| `PUBLIC_BACKEND_URL` | Backend API URL (client-side) | `http://localhost:8888` |
| `PORT`               | Port for the frontend server  | `3000`                  |

## Testing the API

Once the backend is running, you can test the API using curl:

```bash
# Test the health endpoint
curl http://localhost:8888/api/health

# Generate an Open Graph image
curl -X POST \
  -F "url=https://example.com" \
  -F "title=Example Title" \
  -F "description=Example Description" \
  http://localhost:8888/api/generate
```

## Development Workflow

1. **Backend Development**:

   - Make changes to the Go code
   - Run `go run .` to restart the server
   - For hot reloading, consider using [Air](https://github.com/cosmtrek/air)

2. **Frontend Development**:
   - Make changes to the Svelte/Astro code
   - The development server will automatically reload

## Common Issues

### Chrome/Chromium Not Found

If you encounter issues with Chrome not being found, set the `CHROME_PATH` environment variable to the location of your Chrome executable:

```bash
# Linux example
export CHROME_PATH=/usr/bin/google-chrome

# macOS example
export CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

# Windows example
set CHROME_PATH="C:\Program Files\Google\Chrome\Application\chrome.exe"
```

### CORS Issues

If you're experiencing CORS issues:

1. Ensure `ENABLE_CORS=true` is set in the backend `.env` file
2. Verify that `BACKEND_URL` and `PUBLIC_BACKEND_URL` in the frontend `.env` file match the URL where your backend is running

### Output Directory

The backend needs permission to write to the output directory. If you encounter permission issues:

```bash
mkdir -p outputs
chmod 755 outputs
```

## Admin Features

To access admin features locally:

1. Set the `ADMIN_TOKEN` in your backend `.env` file
2. Access the admin page at http://localhost:3000/admin
3. Enter the token you set to authenticate

## Building for Production

To build the application for production:

### Backend

```bash
cd backend
go build -o og-generator .
```

### Frontend

```bash
cd frontend
pnpm build
```

The frontend build will be available in the `dist` directory, which can be served by any static file server.
