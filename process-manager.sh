#!/bin/bash

# Process management utility for OG Drip deployment
# Provides utilities for managing the application processes

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_FILE="/tmp/ogdrip.pid"
LOG_FILE="/tmp/ogdrip.log"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

show_usage() {
    echo "Usage: $0 {start|stop|restart|status|health|logs|cleanup}"
    echo ""
    echo "Commands:"
    echo "  start    - Start the OG Drip service"
    echo "  stop     - Stop the OG Drip service"
    echo "  restart  - Restart the OG Drip service"
    echo "  status   - Show service status"
    echo "  health   - Run health check"
    echo "  logs     - Show service logs"
    echo "  cleanup  - Clean up temporary files and processes"
}

is_running() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            return 0
        else
            rm -f "$PID_FILE"
            return 1
        fi
    fi
    return 1
}

start_service() {
    if is_running; then
        log "${YELLOW}Service is already running (PID: $(cat $PID_FILE))${NC}"
        return 0
    fi
    
    log "${BLUE}Starting OG Drip service...${NC}"
    
    # Start the service in background
    cd "$SCRIPT_DIR"
    nohup ./start.sh > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    
    # Wait a moment and check if it started successfully
    sleep 3
    if is_running; then
        log "${GREEN}✓ Service started successfully (PID: $(cat $PID_FILE))${NC}"
        return 0
    else
        log "${RED}✗ Failed to start service${NC}"
        return 1
    fi
}

stop_service() {
    if ! is_running; then
        log "${YELLOW}Service is not running${NC}"
        return 0
    fi
    
    PID=$(cat "$PID_FILE")
    log "${BLUE}Stopping OG Drip service (PID: $PID)...${NC}"
    
    # Send SIGTERM for graceful shutdown
    kill -TERM "$PID" 2>/dev/null || true
    
    # Wait up to 30 seconds for graceful shutdown
    for i in {1..30}; do
        if ! kill -0 "$PID" 2>/dev/null; then
            rm -f "$PID_FILE"
            log "${GREEN}✓ Service stopped gracefully${NC}"
            return 0
        fi
        sleep 1
    done
    
    # Force kill if still running
    log "${YELLOW}Forcing service shutdown...${NC}"
    kill -KILL "$PID" 2>/dev/null || true
    rm -f "$PID_FILE"
    log "${GREEN}✓ Service stopped (forced)${NC}"
}

restart_service() {
    log "${BLUE}Restarting OG Drip service...${NC}"
    stop_service
    sleep 2
    start_service
}

show_status() {
    if is_running; then
        PID=$(cat "$PID_FILE")
        log "${GREEN}✓ Service is running (PID: $PID)${NC}"
        
        # Show process details
        ps -p "$PID" -o pid,ppid,cmd,etime,pcpu,pmem 2>/dev/null || true
        
        return 0
    else
        log "${RED}✗ Service is not running${NC}"
        return 1
    fi
}

run_health_check() {
    log "${BLUE}Running health check...${NC}"
    if [ -f "$SCRIPT_DIR/healthcheck.sh" ]; then
        "$SCRIPT_DIR/healthcheck.sh"
    else
        log "${YELLOW}Health check script not found${NC}"
        return 1
    fi
}

show_logs() {
    if [ -f "$LOG_FILE" ]; then
        tail -f "$LOG_FILE"
    else
        log "${YELLOW}No log file found${NC}"
    fi
}

cleanup() {
    log "${BLUE}Cleaning up...${NC}"
    
    # Stop service if running
    if is_running; then
        stop_service
    fi
    
    # Clean up temporary files
    rm -f "$PID_FILE"
    
    # Clean up old logs (keep last 100 lines)
    if [ -f "$LOG_FILE" ]; then
        tail -n 100 "$LOG_FILE" > "${LOG_FILE}.tmp"
        mv "${LOG_FILE}.tmp" "$LOG_FILE"
    fi
    
    # Kill any orphaned processes
    pkill -f "^ogdrip-backend$" 2>/dev/null || true
    pgrep -f "chromium.*ogdrip" | while read -r pid; do
        cmdline=$(tr '\0' ' ' < /proc/$pid/cmdline)
        if [[ "$cmdline" == *"--ogdrip-specific-flag"* ]]; then
            kill "$pid"
        fi
    done 2>/dev/null || true
    
    log "${GREEN}✓ Cleanup completed${NC}"
}

# Main command processing
case "${1:-}" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    status)
        show_status
        ;;
    health)
        run_health_check
        ;;
    logs)
        show_logs
        ;;
    cleanup)
        cleanup
        ;;
    *)
        show_usage
        exit 1
        ;;
esac