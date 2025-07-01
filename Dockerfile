# Multi-stage Dockerfile for OG Drip Application
# Builds both Go backend and Astro frontend in a single container

# Stage 1: Frontend Builder
FROM node:22-alpine AS frontend-builder

# Set working directory
WORKDIR /app

# Copy workspace configuration files
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json ./

# Enable corepack for pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy source code for shared and frontend
COPY shared/ ./shared/
COPY frontend/ ./frontend/

# Install dependencies
RUN pnpm install --frozen-lockfile

# Build shared types and frontend
RUN pnpm turbo build --filter=@ogdrip/shared --filter=@ogdrip/frontend

# Stage 2: Backend Builder  
FROM golang:1.24-alpine AS backend-builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./backend/

# Download go dependencies
WORKDIR /app/backend
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the backend binary
RUN mkdir -p build && go build -o build/ogdrip-backend *.go

# Stage 3: Final Runtime Image
FROM node:22-alpine AS runtime

# Set environment variables for Chromium (when available)
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_FLAGS="--headless --disable-gpu --disable-software-rasterizer --disable-dev-shm-usage --no-sandbox --disable-background-timer-throttling --disable-backgrounding-occluded-windows --disable-renderer-backgrounding"

# Create app directory and necessary subdirectories
WORKDIR /app
RUN mkdir -p outputs data

# Enable corepack for pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy built frontend from frontend-builder stage
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
COPY --from=frontend-builder /app/frontend/package.json ./frontend/package.json
COPY --from=frontend-builder /app/shared/dist ./shared/dist

# Copy built backend from backend-builder stage
COPY --from=backend-builder /app/backend/build/ogdrip-backend ./backend/build/ogdrip-backend

# Copy startup script and make it executable
COPY start.sh ./
RUN chmod +x start.sh

# Install only production dependencies for frontend runtime
WORKDIR /app/frontend
RUN pnpm install --prod --frozen-lockfile

# Set working directory back to app root
WORKDIR /app

# Set production environment
ENV NODE_ENV=production
ENV GO111MODULE=on

# Expose ports (backend: 8888, frontend: 3000)
EXPOSE 8888 3000

# Start the application using the existing start script
CMD ["./start.sh"]