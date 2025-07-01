#!/bin/bash

# Deployment validation script
# This script validates that all deployment configurations are correct

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

errors=0
warnings=0

log_error() {
    echo -e "${RED}✗ ERROR: $1${NC}"
    errors=$((errors + 1))
}

log_warning() {
    echo -e "${YELLOW}⚠ WARNING: $1${NC}"
    warnings=$((warnings + 1))
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

check_file_exists() {
    local file="$1"
    local description="$2"
    
    if [ -f "$file" ]; then
        log_success "$description exists: $file"
        return 0
    else
        log_error "$description missing: $file"
        return 1
    fi
}

check_file_executable() {
    local file="$1"
    local description="$2"
    
    if [ -x "$file" ]; then
        log_success "$description is executable: $file"
        return 0
    else
        log_error "$description is not executable: $file"
        return 1
    fi
}

check_go_build() {
    log_info "Checking Go build process..."
    
    cd "$SCRIPT_DIR/backend"
    
    # Check if Go files exist
    if ! ls *.go >/dev/null 2>&1; then
        log_error "No Go source files found in backend/"
        return 1
    fi
    
    # Test build
    if go build -o /tmp/test-ogdrip-backend *.go; then
        log_success "Go build successful"
        rm -f /tmp/test-ogdrip-backend
        return 0
    else
        log_error "Go build failed"
        return 1
    fi
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check if Go is available
    if command -v go >/dev/null 2>&1; then
        log_success "Go is available: $(go version)"
    else
        log_error "Go is not available"
    fi
    
    # Check if Node.js is available
    if command -v node >/dev/null 2>&1; then
        log_success "Node.js is available: $(node --version)"
    else
        log_warning "Node.js is not available (required for frontend)"
    fi
    
    # Check if pnpm is available
    if command -v pnpm >/dev/null 2>&1; then
        log_success "pnpm is available: $(pnpm --version)"
    else
        log_warning "pnpm is not available (required for frontend)"
    fi
}

check_configuration_files() {
    log_info "Checking configuration files..."
    
    # Core deployment files
    check_file_exists "$SCRIPT_DIR/start.sh" "Main startup script"
    check_file_exists "$SCRIPT_DIR/nixpacks.toml" "Nixpacks configuration"
    check_file_exists "$SCRIPT_DIR/coolify.json" "Coolify configuration"
    
    # Utility scripts
    check_file_exists "$SCRIPT_DIR/healthcheck.sh" "Health check script"
    check_file_exists "$SCRIPT_DIR/process-manager.sh" "Process manager script"
    
    # Environment files
    check_file_exists "$SCRIPT_DIR/backend/.env.production" "Backend production environment"
    check_file_exists "$SCRIPT_DIR/frontend/.env.production" "Frontend production environment"
    
    # Check if scripts are executable
    check_file_executable "$SCRIPT_DIR/start.sh" "Main startup script"
    check_file_executable "$SCRIPT_DIR/healthcheck.sh" "Health check script"
    check_file_executable "$SCRIPT_DIR/process-manager.sh" "Process manager script"
}

check_environment_files() {
    log_info "Checking environment file content..."
    
    # Check backend environment
    if [ -f "$SCRIPT_DIR/backend/.env.production" ]; then
        if grep -q "BASE_URL=" "$SCRIPT_DIR/backend/.env.production"; then
            log_success "Backend environment has BASE_URL configured"
        else
            log_warning "Backend environment missing BASE_URL configuration"
        fi
        
        if grep -q "CHROME_PATH=" "$SCRIPT_DIR/backend/.env.production"; then
            log_success "Backend environment has CHROME_PATH configured"
        else
            log_warning "Backend environment missing CHROME_PATH configuration"
        fi
    fi
    
    # Check frontend environment  
    if [ -f "$SCRIPT_DIR/frontend/.env.production" ]; then
        if grep -q "PUBLIC_BACKEND_URL=" "$SCRIPT_DIR/frontend/.env.production"; then
            log_success "Frontend environment has PUBLIC_BACKEND_URL configured"
        else
            log_warning "Frontend environment missing PUBLIC_BACKEND_URL configuration"
        fi
    fi
}

check_nixpacks_config() {
    log_info "Checking nixpacks configuration..."
    
    if [ -f "$SCRIPT_DIR/nixpacks.toml" ]; then
        # Check for Chromium
        if grep -q "chromium" "$SCRIPT_DIR/nixpacks.toml"; then
            log_success "Nixpacks includes Chromium dependency"
        else
            log_error "Nixpacks missing Chromium dependency"
        fi
        
        # Check for font support
        if grep -q "font" "$SCRIPT_DIR/nixpacks.toml"; then
            log_success "Nixpacks includes font dependencies"
        else
            log_warning "Nixpacks may be missing font dependencies"
        fi
        
        # Check build commands
        if grep -q "go build.*ogdrip-backend" "$SCRIPT_DIR/nixpacks.toml"; then
            log_success "Nixpacks has correct Go build command"
        else
            log_error "Nixpacks missing or incorrect Go build command"
        fi
    fi
}

check_coolify_config() {
    log_info "Checking Coolify configuration..."
    
    if [ -f "$SCRIPT_DIR/coolify.json" ]; then
        # Check for nixpacks builder
        if grep -q '"builder".*"nixpacks"' "$SCRIPT_DIR/coolify.json"; then
            log_success "Coolify configured to use nixpacks"
        else
            log_warning "Coolify may not be using nixpacks builder"
        fi
        
        # Check for health check
        if grep -q "/api/health" "$SCRIPT_DIR/coolify.json"; then
            log_success "Coolify has health check configured"
        else
            log_warning "Coolify missing health check configuration"
        fi
        
        # Check for persistent volumes
        if grep -q "persistent.*true" "$SCRIPT_DIR/coolify.json"; then
            log_success "Coolify has persistent volumes configured"
        else
            log_warning "Coolify may be missing persistent volume configuration"
        fi
    fi
}

check_start_script() {
    log_info "Checking start script..."
    
    if [ -f "$SCRIPT_DIR/start.sh" ]; then
        # Check for compiled binary usage
        if grep -q "ogdrip-backend.*-service" "$SCRIPT_DIR/start.sh"; then
            log_success "Start script uses compiled binary"
        else
            log_error "Start script may not be using compiled binary"
        fi
        
        # Check for signal handling
        if grep -q "trap.*cleanup" "$SCRIPT_DIR/start.sh"; then
            log_success "Start script has signal handling"
        else
            log_warning "Start script may lack proper signal handling"
        fi
        
        # Check for error handling
        if grep -q "set -e" "$SCRIPT_DIR/start.sh"; then
            log_success "Start script has error handling"
        else
            log_warning "Start script may lack error handling"
        fi
    fi
}

run_tests() {
    log_info "Running backend tests..."
    
    cd "$SCRIPT_DIR/backend"
    if go test -v; then
        log_success "Backend tests passed"
    else
        log_error "Backend tests failed"
    fi
}

main() {
    echo "=================================================="
    echo "          OG Drip Deployment Validation"
    echo "=================================================="
    echo ""
    
    check_dependencies
    echo ""
    
    check_configuration_files
    echo ""
    
    check_environment_files
    echo ""
    
    check_nixpacks_config
    echo ""
    
    check_coolify_config
    echo ""
    
    check_start_script
    echo ""
    
    check_go_build
    echo ""
    
    run_tests
    echo ""
    
    echo "=================================================="
    echo "                    SUMMARY"
    echo "=================================================="
    
    if [ $errors -eq 0 ] && [ $warnings -eq 0 ]; then
        echo -e "${GREEN}✓ All checks passed! Deployment is ready.${NC}"
        exit 0
    elif [ $errors -eq 0 ]; then
        echo -e "${YELLOW}⚠ Deployment is mostly ready with $warnings warnings.${NC}"
        echo "Consider addressing the warnings for optimal deployment."
        exit 0
    else
        echo -e "${RED}✗ Deployment has $errors errors and $warnings warnings.${NC}"
        echo "Please fix the errors before deploying."
        exit 1
    fi
}

main "$@"