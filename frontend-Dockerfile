FROM node:22-slim

# Note: This base image will be pulled from Docker Hub.
# If you encounter authentication issues, you may need to:
# 1. Login to Docker Hub: docker login
# 2. Or use a local development approach without Docker (see LOCAL_DEPLOYMENT.md)

# Set the working directory
WORKDIR /app

# Install pnpm and serve for static file serving
RUN npm install -g pnpm serve

# Copy key configuration files first
COPY package.json pnpm-lock.yaml* ./
COPY tsconfig.json ./
COPY svelte.config.js ./
COPY astro.config.mjs ./

# Copy source files with specific attention to content collections
COPY src/components src/components/
COPY src/layouts src/layouts/
COPY src/pages src/pages/
COPY src/content src/content/
COPY src/utils src/utils/
COPY src/test src/test/
COPY src/env.d.ts src/env.d.ts
COPY public public/

# Install dependencies using pnpm
RUN pnpm install

# Build the static site
RUN pnpm run build

# Runtime environment variables
ENV NODE_ENV=production
ENV ASTRO_TELEMETRY_DISABLED=1

# Expose the port
EXPOSE 3000

# Serve the static site
CMD ["serve", "-s", "dist", "-l", "3000"] 