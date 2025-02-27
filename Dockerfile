# Build stage
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the main executable
RUN CGO_ENABLED=0 GOOS=linux go build -o og-generator .

# Runtime stage
FROM zenika/alpine-chrome:with-node

# Set working directory
WORKDIR /app

# Copy the binary
COPY --from=builder /app/og-generator .

# Create output directory
RUN mkdir -p /app/outputs

# Set environment variables
ENV PORT=8888
ENV BASE_URL=http://localhost:8888
ENV OUTPUT_DIR=/app/outputs

# Expose port
EXPOSE 8888

# Run the service in API mode
CMD ["./og-generator", "-service"] 