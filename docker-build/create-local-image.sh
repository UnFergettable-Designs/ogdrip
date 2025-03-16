#!/bin/bash
set -e

# This script creates a local image from the Ubuntu base image
# It avoids pulling from a registry by using debootstrap to create a local tar file

echo "Creating local base image for container testing..."

# Check if debootstrap is installed
if ! command -v debootstrap &> /dev/null; then
    echo "Error: debootstrap is not installed"
    echo "Please install it with: brew install debootstrap"
    exit 1
fi

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
echo "Working in temporary directory: $TEMP_DIR"

# Create a minimal Ubuntu file system using debootstrap
echo "Creating minimal Ubuntu file system (this may take a few minutes)..."
sudo debootstrap --variant=minbase jammy "$TEMP_DIR/rootfs"

# Add necessary packages
sudo chroot "$TEMP_DIR/rootfs" apt-get update
sudo chroot "$TEMP_DIR/rootfs" apt-get install -y --no-install-recommends \
    ca-certificates \
    chromium-browser \
    libsqlite3-0 \
    curl \
    nodejs \
    npm \
    supervisor

# Create a tar file from the file system
echo "Creating tar file from file system..."
sudo tar -C "$TEMP_DIR/rootfs" -c . | sudo nerdctl load

# Import the tar file as a local image
echo "Importing as local image..."
nerdctl tag library/import:latest ubuntu:local

# Clean up
echo "Cleaning up..."
sudo rm -rf "$TEMP_DIR"

echo "Base image ubuntu:local is now available for local use"
echo "You can now build your container using:"
echo "nerdctl build -f docker-build/Dockerfile.local -t ogdrip-local:test ."
