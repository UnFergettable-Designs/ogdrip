FROM golang:1.20-alpine

# Install required dependencies
RUN apk add --no-cache chromium ca-certificates git gcc musl-dev

# Set environment variables
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_FLAGS=--headless,--disable-gpu,--no-sandbox
ENV CGO_ENABLED=1

# Set the working directory
WORKDIR /app

# Create necessary files for minimal build
COPY <<EOF /app/go.mod
module og-generator

go 1.20
EOF

COPY <<EOF /app/main.go
package main

import "fmt"

func main() {
  fmt.Println("Open Graph Generator API")
}
EOF

# Build the executable
RUN go build -o og-generator .

# Create output directory
RUN mkdir -p /app/outputs

# Expose port
EXPOSE 8888

# Command that will be overridden at runtime
CMD ["./og-generator"] 