#!/bin/bash
# Test Nginx configuration with production setup

echo "Cleaning up any existing test environment..."
nerdctl compose -f docker-compose.nginx-test.yml down --volumes

echo "Building and starting containers for Nginx configuration testing..."
nerdctl compose -f docker-compose.nginx-test.yml build
nerdctl compose -f docker-compose.nginx-test.yml up

# The script will wait here while containers are running
# Press Ctrl+C to stop the containers when done testing

# Cleanup is handled by the trap below
trap 'echo "Cleaning up containers..."; nerdctl compose -f docker-compose.nginx-test.yml down --volumes; echo "Test environment stopped"' EXIT
