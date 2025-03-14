#!/bin/bash
# Test Nginx configuration with pre-built static files

echo "Cleaning up any existing test environment..."
nerdctl compose -f docker-compose.nginx-test.yml down --volumes

# Build projects locally first
echo "Building projects locally..."
pnpm build

echo "Creating Docker test environment with pre-built files..."
# Create temporary Docker Compose file for testing
cat > docker-compose.nginx-static-test.yml << EOL
services:
  # Use a simple static file server for frontend
  frontend:
    image: nginx:alpine
    volumes:
      - ./frontend/dist:/usr/share/nginx/html
    restart: unless-stopped

  # Simple mock backend
  backend:
    image: busybox
    command: sh -c "mkdir -p /app/outputs && echo 'Starting mock API server on port 8888' && while true; do nc -l -p 8888 -e echo -e 'HTTP/1.1 200 OK\r\nContent-Length: 15\r\n\r\n{\"status\":\"ok\"}'; done"
    environment:
      - PORT=8888
      - BASE_URL=http://backend:8888
    volumes:
      - outputs:/app/outputs
    restart: unless-stopped

  # Main Nginx service for testing routing
  nginx:
    image: nginx:latest
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx.test.conf:/etc/nginx/conf.d/default.conf
      - outputs:/app/outputs
    depends_on:
      - frontend
      - backend
    # Command to create self-signed SSL certificates for testing
    command: /bin/bash -c "mkdir -p /etc/letsencrypt/live/og-drip.com &&
              openssl req -x509 -nodes -days 365 -newkey rsa:2048
              -keyout /etc/letsencrypt/live/og-drip.com/privkey.pem
              -out /etc/letsencrypt/live/og-drip.com/fullchain.pem
              -subj '/CN=og-drip.com' &&
              nginx -g 'daemon off;'"

volumes:
  outputs:
EOL

echo "Starting static test environment..."
nerdctl compose -f docker-compose.nginx-static-test.yml up

# The script will wait here while containers are running
# Press Ctrl+C to stop the containers when done testing

# Cleanup is handled by the trap below
trap 'echo "Cleaning up containers..."; nerdctl compose -f docker-compose.nginx-static-test.yml down --volumes; rm -f docker-compose.nginx-static-test.yml; echo "Test environment stopped"' EXIT
