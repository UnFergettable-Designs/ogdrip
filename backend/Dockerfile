FROM golang:1.23-alpine

# Install required dependencies
RUN apk add --no-cache chromium ca-certificates git gcc musl-dev

# Set environment variables
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_FLAGS=--headless,--disable-gpu,--no-sandbox
ENV CGO_ENABLED=1

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Create output and data directories
RUN mkdir -p /app/outputs
RUN mkdir -p /app/data

# Build the executable
RUN go build -o og-generator .

# Expose port
EXPOSE 8888

# Command to run
CMD ["./og-generator", "-service"]