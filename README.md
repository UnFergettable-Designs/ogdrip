# Open Graph Generator

This tool generates Open Graph images and meta tags for enhancing link previews on social media platforms like Facebook, Twitter, and LinkedIn.

## Features

- Generates high-quality screenshots of webpages for social sharing
- Creates properly formatted Open Graph meta tags
- Can be used as a command-line tool or as a REST API service
- Supports customization of image dimensions, title, description, and more
- Includes a preview mode to visualize the social media appearance

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/open-graph-generate.git
cd open-graph-generate

# Install dependencies
go get github.com/chromedp/chromedp
```

## Command-Line Usage

### Basic Usage

```bash
# Generate from a URL
go run server.go -url="https://example.com"

# Generate with custom title and description
go run server.go -url="https://example.com" -title="My Custom Title" -description="A great description for social sharing"

# Generate with preview server
go run server.go -url="https://example.com" -preview
```

### Command-Line Options

| Option          | Description                                 | Default               |
| --------------- | ------------------------------------------- | --------------------- |
| `-url`          | Webpage URL to capture                      | (required)            |
| `-output`       | Output file path for the screenshot         | `og_image.png`        |
| `-html`         | Output file for HTML with meta tags         | `og_meta.html`        |
| `-title`        | Title for Open Graph meta tags              | (extracted from URL)  |
| `-description`  | Description for Open Graph meta tags        | (extracted from URL)  |
| `-type`         | Type for Open Graph meta tags               | `website`             |
| `-site`         | Site name for Open Graph meta tags          | (extracted from URL)  |
| `-target-url`   | Target URL for the content                  | (same as URL)         |
| `-width`        | Width of the Open Graph image               | `1200`                |
| `-height`       | Height of the Open Graph image              | `630`                 |
| `-twitter-card` | Twitter card type                           | `summary_large_image` |
| `-preview`      | Start a local server to preview the OG tags | `false`               |
| `-port`         | Port for the preview server                 | `8080`                |
| `-wait`         | Wait time in milliseconds before capturing  | `8000`                |
| `-quality`      | Screenshot quality (1-100)                  | `90`                  |
| `-selector`     | CSS selector to wait for before capturing   | `body`                |
| `-debug`        | Enable debug mode with additional logging   | `false`               |
| `-verbose`      | Enable verbose logging                      | `false`               |

## Running as a Service

The project includes a service wrapper that exposes the Open Graph generator as a REST API.

### Starting the Service

```bash
go run service.go
```

By default, the service runs on port 8888. You can modify the port and other settings by editing the `config` variable in `service.go`.

### API Endpoints

#### Generate Open Graph Assets

**Endpoint:** `POST /api/generate`

**Form parameters:** Same as the command-line options (url, title, description, etc.)

**Example with curl:**

```bash
curl -X POST \
  http://localhost:8888/api/generate \
  -F 'url=https://example.com' \
  -F 'title=My Website Title' \
  -F 'description=A description for social media'
```

**Response:**

```json
{
  "success": true,
  "message": "Open Graph assets generated successfully",
  "image_url": "http://localhost:8888/outputs/abc123_og_image.png",
  "meta_tags_url": "http://localhost:8888/outputs/abc123_og_meta.html"
}
```

#### Health Check

**Endpoint:** `GET /api/health`

**Example:**

```bash
curl http://localhost:8888/api/health
```

**Response:**

```json
{
  "success": true,
  "message": "Open Graph Generator API is running"
}
```

## Testing

The project includes both unit tests and integration tests to ensure everything works correctly.

### Running Unit Tests

```bash
# Run all unit tests
go test -v

# Run specific test
go test -v -run TestHealthEndpoint
```

### Running Integration Tests

Integration tests require a working Chrome/Chromium installation and will actually generate images and HTML files.

```bash
# Run integration tests
go test -v -tags=integration

# Run a specific integration test
go test -v -tags=integration -run TestRunGeneratorDirectly
```

### Testing the API with the Test Script

A test script is included to quickly test the API service:

```bash
# Make the script executable
chmod +x test_api.sh

# Run the test with default URL (example.com)
./test_api.sh

# Run the test with a specific URL
./test_api.sh https://github.com
```

## Deployment

### Running as a System Service

You can run the API as a system service using systemd:

1. Create a service file at `/etc/systemd/system/opengraph-generator.service`:

```
[Unit]
Description=Open Graph Generator Service
After=network.target

[Service]
ExecStart=/path/to/opengraph-generator
WorkingDirectory=/path/to/opengraph-directory
User=yourusername
Group=yourusername
Restart=always

[Install]
WantedBy=multi-user.target
```

2. Build the executable:

```bash
go build -o opengraph-generator service.go
```

3. Start the service:

```bash
sudo systemctl start opengraph-generator
sudo systemctl enable opengraph-generator
```

### Docker Deployment

The project includes Docker support for easy deployment:

```bash
# Build the Docker image
docker build -t opengraph-generator .

# Run the container
docker run -p 8888:8888 opengraph-generator

# Using docker-compose
docker-compose up -d
```

## Troubleshooting

### SSL Certificate Errors

If you encounter SSL errors when accessing a site, check that:

- The domain is spelled correctly
- The site has a valid SSL certificate

If needed, you can use the `-debug=true` flag to get more information about errors.

### Image Not Generated

If the service reports that no image was generated:

- Check that the URL is accessible and valid
- Try increasing the wait time with `-wait=15000` (15 seconds)
- Check for any specific errors in the logs

## License

[MIT License](LICENSE)
