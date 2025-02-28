FROM golang:1.23-alpine

# Install required dependencies
RUN apk add --no-cache chromium ca-certificates git

# Set environment variables
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_FLAGS=--headless,--disable-gpu,--no-sandbox

# Set the working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the executable with the main function
RUN go build -o og-generator main.go server.go service.go

# Create output directory
RUN mkdir -p /app/outputs

# Expose port
EXPOSE 8888

# Run the API service
CMD ["./og-generator", "-service", "-port", "8888"] 