#!/bin/bash

# Spark Park Cricket - Build Script
# Usage: ./build.sh [environment] [options]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Default values
ENVIRONMENT=""
PUSH_IMAGES=false
CLEAN_BUILD=false
VERBOSE=false

# Production server configuration
PROD_SERVER="15.235.202.148"
PROD_USER="ubuntu"
SSH_KEY_PATH="$HOME/.ssh/spark-cricket-prod"

# Functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

usage() {
    echo "Usage: $0 [environment] [options]"
    echo ""
    echo "Environments:"
    echo "  dev     - Development environment (local)"
    echo "  prod    - Production environment (15.235.202.148)"
    echo ""
    echo "Options:"
    echo "  --push          Push images to registry after building"
    echo "  --clean         Clean build (remove existing images)"
    echo "  --verbose       Enable verbose output"
    echo "  --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev --clean"
    echo "  $0 prod --push"
    echo "  $0 dev --verbose"
}

check_environment() {
    if [[ -z "$ENVIRONMENT" ]]; then
        error "Environment is required"
        usage
        exit 1
    fi
    
    if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" ]]; then
        error "Invalid environment. Must be 'dev' or 'prod'"
        exit 1
    fi
}

check_dependencies() {
    log "Checking dependencies..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
        exit 1
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker compose &> /dev/null; then
        error "Docker Compose is not installed"
        exit 1
    fi
    
    success "Dependencies check passed"
}

clean_build() {
    if [[ "$CLEAN_BUILD" == "true" ]]; then
        log "Cleaning existing images..."
        
        if [[ "$ENVIRONMENT" == "dev" ]]; then
            # Remove existing images
            docker compose -f "$PROJECT_ROOT/docker-compose.yml" down --rmi all 2>/dev/null || true
            docker system prune -f
        else
            # Clean on production server
            ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
                "$PROD_USER@$PROD_SERVER" \
                "cd /opt/spark-cricket && docker compose -f deploy/environments/prod/docker-compose.yml down --rmi all 2>/dev/null || true && docker system prune -f"
        fi
        
        success "Clean build completed"
    fi
}

build_backend() {
    log "Building backend image..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT/backend"
        
        # Build backend image
        docker build -t spark-cricket-backend:latest \
            --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
            --build-arg VCS_REF="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" \
            .
        
        success "Backend image built successfully"
    else
        # Build on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket/backend && docker build -t spark-cricket-backend:latest ."
        
        success "Backend image built on production server"
    fi
}

build_frontend() {
    log "Building frontend image..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT/web"
        
        # Build frontend image
        docker build -t spark-cricket-frontend:latest \
            --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
            --build-arg VCS_REF="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" \
            .
        
        success "Frontend image built successfully"
    else
        # Build on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket/web && docker build -t spark-cricket-frontend:latest ."
        
        success "Frontend image built on production server"
    fi
}

build_all_images() {
    log "Building all images for $ENVIRONMENT environment..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT"
        
        # Build all services
        docker compose -f docker-compose.yml build
        
        success "All images built successfully"
    else
        # Build all services on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket && docker compose -f deploy/environments/prod/docker-compose.yml build"
        
        success "All images built on production server"
    fi
}

push_images() {
    if [[ "$PUSH_IMAGES" == "true" ]]; then
        log "Pushing images to registry..."
        
        warning "Image pushing not implemented yet"
        warning "Configure your Docker registry in the deployment configuration"
        
        # TODO: Implement image pushing to registry
        # docker tag spark-cricket-backend:latest your-registry/spark-cricket-backend:latest
        # docker push your-registry/spark-cricket-backend:latest
    fi
}

show_build_info() {
    log "Build information:"
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        echo "  Environment: Development (Local)"
        echo "  Build Date: $(date -u +'%Y-%m-%d %H:%M:%S UTC')"
        echo "  Git Commit: $(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
        echo "  Git Branch: $(git branch --show-current 2>/dev/null || echo 'unknown')"
    else
        echo "  Environment: Production (15.235.202.148)"
        echo "  Build Date: $(date -u +'%Y-%m-%d %H:%M:%S UTC')"
        echo "  Git Commit: $(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
        echo "  Git Branch: $(git branch --show-current 2>/dev/null || echo 'unknown')"
    fi
    
    echo ""
    log "Built images:"
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        docker images | grep -E "(spark-cricket|cricket)" || echo "  No images found"
    else
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "docker images | grep -E '(spark-cricket|cricket)' || echo '  No images found'"
    fi
}

main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            dev|prod)
                ENVIRONMENT="$1"
                shift
                ;;
            --push)
                PUSH_IMAGES=true
                shift
                ;;
            --clean)
                CLEAN_BUILD=true
                shift
                ;;
            --verbose)
                VERBOSE=true
                set -x
                shift
                ;;
            --help)
                usage
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    check_environment
    check_dependencies
    clean_build
    build_all_images
    push_images
    show_build_info
    
    success "Build completed successfully!"
}

# Run main function with all arguments
main "$@"
