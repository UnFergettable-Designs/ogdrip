#!/bin/bash

# Health check script for the OG Drip service
# This script checks if the backend API is responding properly

set -e

# Configuration
BACKEND_URL="${BACKEND_URL:-http://localhost:8888}"
HEALTH_ENDPOINT="/api/health"
TIMEOUT=10
MAX_RETRIES=3

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

check_health() {
    local url="$1"
    local retry_count=0
    
    while [ $retry_count -lt $MAX_RETRIES ]; do
        if curl -f -s --max-time $TIMEOUT "$url" >/dev/null 2>&1; then
            return 0
        fi
        
        retry_count=$((retry_count + 1))
        if [ $retry_count -lt $MAX_RETRIES ]; then
            log "${YELLOW}Health check attempt $retry_count failed, retrying...${NC}"
            sleep 2
        fi
    done
    
    return 1
}

# Main health check
log "Starting health check for OG Drip service..."

# Check backend health
HEALTH_URL="${BACKEND_URL}${HEALTH_ENDPOINT}"
log "Checking backend health at: $HEALTH_URL"

if check_health "$HEALTH_URL"; then
    log "${GREEN}✓ Backend is healthy${NC}"
    
    # Additional checks
    # Check if we can get a proper JSON response
    RESPONSE=$(curl -s --max-time $TIMEOUT "$HEALTH_URL" 2>/dev/null || echo "")
    if echo "$RESPONSE" | grep -q "success.*true"; then
        log "${GREEN}✓ Backend API is responding correctly${NC}"
        exit 0
    else
        log "${YELLOW}⚠ Backend is responding but may have issues${NC}"
        log "Response: $RESPONSE"
        exit 1
    fi
else
    log "${RED}✗ Backend health check failed${NC}"
    
    # Try to get more information about the failure
    log "Attempting to get more information..."
    curl -v --max-time $TIMEOUT "$HEALTH_URL" 2>&1 | head -20 || true
    
    exit 1
fi