#!/bin/bash
# Helper script to set up the best container engine for development

# Check if script is being run with color support
if [ -t 1 ]; then
  BOLD="\033[1m"
  GREEN="\033[0;32m"
  YELLOW="\033[0;33m"
  RED="\033[0;31m"
  BLUE="\033[0;34m"
  NC="\033[0m" # No Color
else
  BOLD=""
  GREEN=""
  YELLOW=""
  RED=""
  BLUE=""
  NC=""
fi

echo -e "${BOLD}Container Runtime Setup Helper${NC}"
echo "This script will help you choose and set up the best container engine for your development needs."
echo ""

# Check for Homebrew
if ! command -v brew &> /dev/null; then
  echo -e "${RED}Homebrew not found. It's needed to install container tools.${NC}"
  echo "Visit https://brew.sh to install Homebrew first."
  exit 1
fi

echo -e "${BOLD}Available Container Engines:${NC}"
echo -e "1. ${GREEN}Docker Desktop${NC} - Full featured, user friendly, but requires license for commercial use"
echo -e "2. ${GREEN}Rancher Desktop${NC} - Open source, supports both Docker and containerd/nerdctl APIs"
echo -e "3. ${GREEN}Podman${NC} - Daemonless alternative to Docker, OCI-compliant"
echo -e "4. ${GREEN}Colima${NC} - Lightweight Docker/containerd for macOS (you've been using this)"
echo ""
echo -e "${YELLOW}Which container engine would you like to set up? (1-4)${NC}"

read -p "Enter choice (1-4): " CHOICE

case $CHOICE in
  1)
    echo -e "${BLUE}Setting up Docker Desktop...${NC}"
    echo "Checking if Docker Desktop is already installed..."
    if [ -d "/Applications/Docker.app" ]; then
      echo "Docker Desktop is already installed."
    else
      echo "Installing Docker Desktop via Homebrew..."
      brew install --cask docker
    fi
    echo ""
    echo "Docker Desktop requires manual startup. Please:"
    echo "1. Launch Docker Desktop from your Applications folder"
    echo "2. Complete any first-time setup steps"
    echo -e "${GREEN}Once Docker Desktop is running, you can use the standard Docker commands or npm scripts:${NC}"
    echo "  docker build -t ogdrip-frontend:local ./frontend"
    echo "  npm run docker:build"
    ;;
    
  2)
    echo -e "${BLUE}Setting up Rancher Desktop...${NC}"
    if [ -d "/Applications/Rancher Desktop.app" ]; then
      echo "Rancher Desktop is already installed."
    else
      echo "Installing Rancher Desktop via Homebrew..."
      brew install --cask rancher
    fi
    echo ""
    echo "Rancher Desktop requires manual startup. Please:"
    echo "1. Launch Rancher Desktop from your Applications folder"
    echo "2. In preferences, select your preferred container runtime (containerd recommended)"
    echo -e "${GREEN}Once Rancher Desktop is running, you can use Docker or nerdctl commands:${NC}"
    echo "  docker build -t ogdrip-frontend:local ./frontend"
    echo "  nerdctl build -t ogdrip-frontend:local ./frontend"
    ;;
    
  3)
    echo -e "${BLUE}Setting up Podman...${NC}"
    if command -v podman &> /dev/null; then
      echo "Podman is already installed."
    else
      echo "Installing Podman via Homebrew..."
      brew install podman
    fi
    
    echo "Initializing Podman machine..."
    podman machine init
    podman machine start
    
    echo -e "${GREEN}Podman is ready to use with the following commands:${NC}"
    echo "  podman build -t ogdrip-frontend:local ./frontend"
    echo "  npm run podman:build:frontend"
    ;;
    
  4)
    echo -e "${BLUE}Setting up Colima...${NC}"
    if command -v colima &> /dev/null; then
      echo "Colima is already installed."
    else
      echo "Installing Colima and nerdctl via Homebrew..."
      brew install colima nerdctl
    fi
    
    echo "Starting Colima with containerd runtime..."
    colima stop containerd 2>/dev/null || true
    colima start containerd --runtime containerd --cpu 2 --memory 4
    
    # Wait for containerd socket
    SOCKET_PATH="$HOME/.colima/containerd/containerd/containerd.sock"
    max_attempts=10
    attempts=0
    while [ ! -S "$SOCKET_PATH" ] && [ $attempts -lt $max_attempts ]; do
        sleep 3
        ((attempts++))
        echo "Waiting for containerd socket... Attempt $attempts/$max_attempts"
    done
    
    if [ -S "$SOCKET_PATH" ]; then
      export CONTAINER_RUNTIME_ENDPOINT="unix:///$SOCKET_PATH"
      echo -e "${GREEN}Colima is ready to use with the following commands:${NC}"
      echo "  CONTAINER_RUNTIME_ENDPOINT=unix:///$SOCKET_PATH nerdctl build -t ogdrip-frontend:local ./frontend"
      echo "  npm run colima:nerdctl:build:frontend"
    else
      echo -e "${RED}Failed to start Colima with containerd runtime.${NC}"
    fi
    ;;
    
  *)
    echo -e "${RED}Invalid choice. Please run the script again and select a number from 1-4.${NC}"
    exit 1
    ;;
esac

echo ""
echo -e "${BOLD}Setup Complete${NC}"
echo "For more information on container tools, see the documentation."
