# OG Drip - Open Graph Generator

A modern, efficient service for generating Open Graph images and metadata for your web pages.

## Features

- Generate beautiful Open Graph images from URLs
- Customizable templates and styles
- API-first design
- Built with Go and Astro+Svelte
- Automatic image optimization
- Secure admin interface

## Tech Stack

- **Frontend**: Astro + Svelte 5
- **Backend**: Go with ChromeDP
- **Database**: SQLite
- **Build System**: Turborepo + pnpm
- **Deployment**: Coolify with nixpacks

## Quick Start

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/ogdrip.git
   cd ogdrip
   ```

2. Install dependencies:

   ```bash
   pnpm install
   ```

3. Set up environment variables:

   ```bash
   cp frontend/.env.example frontend/.env
   cp backend/.env.example backend/.env
   ```

4. Start development servers:
   ```bash
   pnpm dev
   ```

The frontend will be available at http://localhost:3000 and the backend at http://localhost:8888.

## Development

This is a monorepo managed with Turborepo and pnpm workspaces. The main components are:

- `frontend/`: Astro + Svelte frontend
- `backend/`: Go backend service
- `shared/`: Shared TypeScript types

### Available Scripts

- `pnpm dev` - Start all development servers
- `pnpm build` - Build all packages
- `pnpm test` - Run tests
- `pnpm lint` - Lint all packages

## Deployment

This project is designed to be deployed on Coolify using nixpacks.

### Prerequisites

- A Coolify instance
- Git repository connected to Coolify
- Domain name (optional but recommended)

### Deployment Steps

1. In Coolify:

   - Create a new service
   - Choose "Source: GitHub"
   - Select your repository
   - Choose "Build Pack: Nixpacks"

2. Configure environment variables:

   ```
   PUBLIC_BACKEND_URL=https://your-domain.com
   BACKEND_URL=https://your-domain.com
   ADMIN_TOKEN=your-secure-admin-token
   ```

3. Add your domain in Coolify settings

Coolify will handle:

- Building the application
- Setting up SSL certificates
- Managing the reverse proxy
- Automatic deployments on push

## API Documentation

The API documentation is available through Swagger UI when the backend is running:

- Interactive documentation: `/docs/`
- OpenAPI specification: `/api/openapi.yaml`

## Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Run tests: `pnpm test`
5. Submit a pull request

## License

[MIT License](LICENSE)

## Acknowledgments

- [Chromedp](https://github.com/chromedp/chromedp) for headless browser automation
- [Astro](https://astro.build/) for the frontend framework
- [Svelte](https://svelte.dev/) for reactive UI components
- [Coolify](https://coolify.io/) for deployment platform
