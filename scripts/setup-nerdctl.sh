#!/bin/bash
# Setup script for nerdctl with macOS

set -e

echo "Setting up nerdctl for OGDrip..."

# Check if nerdctl is installed
if ! command -v nerdctl &> /dev/null; then
    echo "nerdctl is not installed."
    echo "Installing nerdctl with Homebrew..."
    brew install nerdctl
fi

# Check which container runtime is available
if command -v docker &> /dev/null; then
    echo "Docker is detected, checking if it's running..."
    if docker info &> /dev/null; then
        echo "Docker is running. You can use nerdctl with Docker."
        echo ""
        echo "You can now use the npm/pnpm scripts to build and run containers:"
        echo "  npm run nerdctl:build:frontend"
        echo "  npm run nerdctl:run:frontend"
        echo ""
        exit 0
    fi
fi

if command -v colima &> /dev/null; then
    echo "Colima is detected. Setting up nerdctl with Colima..."
    
    # Check if Colima is running
    if ! colima status &> /dev/null; then
        echo "Starting Colima with containerd runtime..."
        colima start --runtime containerd
    fi
    
    # Wait for containerd socket
    echo "Waiting for containerd socket..."
    max_attempts=10
    attempts=0
    while [ ! -S "$HOME/.colima/default/containerd/containerd.sock" ] && [ $attempts -lt $max_attempts ]; do
        sleep 2
        ((attempts++))
        echo "Waiting for containerd socket... Attempt $attempts/$max_attempts"
    done

    if [ ! -S "$HOME/.colima/default/containerd/containerd.sock" ]; then
        echo "Error: containerd socket not found after $max_attempts attempts"
        echo "Please check Colima status with: colima status"
        exit 1
    fi

    # Set up environment variable for nerdctl
    export CONTAINER_RUNTIME_ENDPOINT=unix:///$HOME/.colima/default/containerd/containerd.sock
    
    echo "Testing nerdctl connection..."
    if nerdctl info &> /dev/null; then
        echo "Success! nerdctl is properly connected to Colima's containerd"
        echo ""
        echo "To use nerdctl with Colima, you need to set this environment variable:"
        echo "export CONTAINER_RUNTIME_ENDPOINT=unix:///$HOME/.colima/default/containerd/containerd.sock"
        echo ""
        echo "You can add it to your shell profile (~/.zshrc or ~/.bash_profile)"
        echo ""
        echo "Then you can use the npm/pnpm scripts:"
        echo "npm run nerdctl:build:frontend"
        echo "npm run nerdctl:run:frontend"
        echo ""
    else
        echo "Error: nerdctl failed to connect to containerd"
        echo "Please check Colima status with: colima status"
        exit 1
    fi
elif command -v lima &> /dev/null; then
    echo "Lima is detected. Setting up nerdctl with Lima..."
    # TODO: Add Lima specific setup if needed
    echo "Please refer to Lima documentation for nerdctl setup."
elif command -v "Rancher Desktop" &> /dev/null || pgrep -f "Rancher Desktop" > /dev/null; then
    echo "Rancher Desktop is detected..."
    echo "Make sure Rancher Desktop is configured to use containerd as the runtime (not dockerd)"
    echo ""
    echo "Testing nerdctl connection..."
    if nerdctl info &> /dev/null; then
        echo "Success! nerdctl is properly connected to Rancher Desktop's containerd"
        echo ""
        echo "You can now use the npm/pnpm scripts to build and run containers:"
        echo "npm run nerdctl:build:frontend"
        echo "npm run nerdctl:run:frontend"
        echo ""
    else
        echo "Error: nerdctl failed to connect to containerd in Rancher Desktop"
        echo "Please ensure Rancher Desktop is running and using containerd runtime."
        exit 1
    fi
else
    echo "No supported container runtime found (Docker, Colima, Lima, or Rancher Desktop)"
    echo "Please install one of these container runtimes before using nerdctl."
    exit 1
fi