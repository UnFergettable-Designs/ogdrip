# ogdrip

Social presence that drips with style. A powerful Open Graph image and metadata generator for enhancing your content's appearance on social platforms.

## Features

- Generate Open Graph images from websites or custom content
- Create comprehensive meta tags for social media sharing
- Preview how content will look on different platforms
- API service for integration with other applications

## Project Structure

This is a monorepo with the following components:

- **Frontend**: Astro + Svelte 5 application with Deno runtime
- **Backend**: Go service with ChromeDP for screenshot and image generation
- **Shared**: Shared TypeScript types and utilities

## Running with Docker

The easiest way to run the application is with Docker Compose:

```bash
# Clone the repository
git clone https://github.com/yourusername/open-graph-generate.git
cd open-graph-generate

# Start the application with Docker Compose
docker-compose up -d
```

The application will be available at:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8888

## Development Setup

### Prerequisites

- [Deno](https://deno.land/) 1.39 or higher
- [Go](https://golang.org/) 1.23 or higher
- [Node.js](https://nodejs.org/) (optional, for npm packages)

### Running Locally

```bash
# Start both frontend and backend
deno task dev

# Or start them separately
deno task dev:frontend
deno task dev:backend
```

## Deployment on Coolify

This application is designed to work well with Coolify:

1. Connect your Git repository to Coolify
2. Choose "Docker Compose" as the deployment method
3. Use the docker-compose.yml from the repository
4. Configure the following environment variables:
   - `BASE_URL`: The public URL where your application will be hosted
   - `ENABLE_CORS`: Set to `true` if needed
   - `OUTPUT_DIR`: Directory for storing generated images

## Configuration

### Environment Variables

- Frontend

  - `BACKEND_URL`: URL of the backend API (default: http://localhost:8888)
  - `PORT`: Port for the frontend server (default: 3000)

- Backend
  - `PORT`: Port for the backend server (default: 8888)
  - `BASE_URL`: Public URL for the backend API
  - `OUTPUT_DIR`: Directory for storing generated images
  - `ENABLE_CORS`: Whether to enable CORS (default: true)
  - `MAX_QUEUE_SIZE`: Maximum number of concurrent tasks (default: 10)

## License

[MIT License](LICENSE)
