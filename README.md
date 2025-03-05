# Open Graph Generator

A powerful tool for generating Open Graph images and meta tags for your websites and applications.

## Features

- Generate beautiful Open Graph images from URLs
- Create meta tags for social media sharing
- Customize titles, descriptions, and visual elements
- Admin dashboard for tracking generations
- API for integration with your applications
- Downloadable assets (images, HTML, and ZIP packages)

## Quick Start

The quickest way to get started is to use Docker Compose:

```bash
git clone https://github.com/your-username/ogdrip.git
cd ogdrip
docker compose up
```

Then visit http://localhost:3000 in your browser.

## Requirements

- Docker and Docker Compose (for containerized setup)
- Go v1.23 or later (for manual backend setup)
- Node.js v18 or later (for manual frontend setup)
- pnpm (for manual frontend setup)
- Chrome/Chromium (for headless browser functionality)

## Documentation

- [Local Deployment Guide](LOCAL_DEPLOYMENT.md) - How to run the application locally
- [Production Deployment Guide](DEPLOYMENT.md) - How to deploy to production environments

## Architecture

The Open Graph Generator consists of two main components:

1. **Frontend**: Built with Astro and Svelte

   - User interface for creating Open Graph assets
   - Preview functionality
   - Admin dashboard

2. **Backend**: Built with Go
   - API for generating Open Graph images
   - Headless Chrome integration for rendering
   - SQLite database for tracking generations

## API Usage

### Generate Open Graph Assets

```bash
curl -X POST \
  -F "url=https://example.com" \
  -F "title=Example Title" \
  -F "description=Example Description" \
  http://localhost:8888/api/generate
```

Response:

```json
{
  "success": true,
  "message": "Open Graph assets generated successfully",
  "image_url": "http://localhost:8888/outputs/abc123_og_image.png",
  "meta_tags_url": "http://localhost:8888/outputs/abc123_og_meta.html",
  "preview_url": "http://localhost:8888/preview/abc123",
  "zip_url": "http://localhost:8888/api/download/abc123_og_package.zip",
  "id": "abc123"
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Chromedp](https://github.com/chromedp/chromedp) for headless browser automation
- [Astro](https://astro.build/) for the frontend framework
- [Svelte](https://svelte.dev/) for reactive UI components
