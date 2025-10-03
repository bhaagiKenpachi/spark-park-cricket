#!/bin/bash

# Spark Park Cricket - Environment-Agnostic Deployment Script
# Usage: ./deploy.sh [environment] [options]
# Example: ./deploy.sh prod --build --restart

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
DEPLOY_DIR="$PROJECT_ROOT/deploy"

# Default values
ENVIRONMENT=""
BUILD_IMAGES=false
RESTART_SERVICES=false
SKIP_TESTS=false
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
    echo "  --build         Build Docker images before deployment"
    echo "  --restart       Restart all services after deployment"
    echo "  --skip-tests    Skip running tests before deployment"
    echo "  --verbose       Enable verbose output"
    echo "  --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 prod --build --restart"
    echo "  $0 dev --build"
    echo "  $0 prod --skip-tests"
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
    
    success "Deploying to $ENVIRONMENT environment"
}

check_dependencies() {
    log "Checking dependencies..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
        exit 1
    fi
    
    # Check if Docker Compose is installed
    if ! docker compose version &> /dev/null; then
        error "Docker Compose is not installed"
        exit 1
    fi
    
    # Check if SSH key exists for production
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        if [[ ! -f "$SSH_KEY_PATH" ]]; then
            warning "SSH key not found at $SSH_KEY_PATH"
            warning "Please copy your SSH key from ~/dojima/tee-auth/ folder"
            warning "Example: cp ~/dojima/tee-auth/your-key ~/.ssh/spark-cricket-prod && chmod 600 ~/.ssh/spark-cricket-prod"
            exit 1
        fi
    fi
    
    success "Dependencies check passed"
}

setup_ssh_config() {
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        log "Setting up SSH configuration..."
        
        # Create SSH config if it doesn't exist
        SSH_CONFIG="$HOME/.ssh/config"
        if [[ ! -f "$SSH_CONFIG" ]]; then
            touch "$SSH_CONFIG"
            chmod 600 "$SSH_CONFIG"
        fi
        
        # Add host configuration if not exists
        if ! grep -q "Host spark-cricket-prod" "$SSH_CONFIG"; then
            cat >> "$SSH_CONFIG" << EOF

Host spark-cricket-prod
    HostName $PROD_SERVER
    User $PROD_USER
    IdentityFile $SSH_KEY_PATH
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
EOF
            success "SSH configuration added"
        fi
    fi
}

run_tests() {
    if [[ "$SKIP_TESTS" == "true" ]]; then
        warning "Skipping tests as requested"
        return 0
    fi
    
    log "Running tests..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT/backend"
        
        # Run unit tests
        log "Running unit tests..."
        if ! go test ./... -v; then
            error "Unit tests failed"
            exit 1
        fi
        
        # Run integration tests
        log "Running integration tests..."
        if ! make test-integration; then
            error "Integration tests failed"
            exit 1
        fi
        
        success "All tests passed"
    else
        warning "Skipping tests for production deployment"
    fi
}

build_images() {
    if [[ "$BUILD_IMAGES" == "true" ]]; then
        log "Building Docker images..."
        
          if [[ "$ENVIRONMENT" == "dev" ]]; then
              cd "$PROJECT_ROOT"
              docker compose -f docker-compose.yml build
          else
            # For production, we'll build on the remote server
            log "Images will be built on production server"
        fi
        
        success "Docker images built successfully"
    fi
}

deploy_to_dev() {
    log "Deploying to development environment..."
    
    cd "$PROJECT_ROOT"
    
      # Stop existing containers
      log "Stopping existing containers..."
      docker compose -f docker-compose.yml down
      
      # Start services
      log "Starting services..."
      docker compose -f docker-compose.yml up -d
    
    # Wait for services to be healthy
    log "Waiting for services to be healthy..."
    sleep 10
    
    # Check health
    if ! ./deploy/scripts/health-check.sh dev; then
        error "Health check failed"
        exit 1
    fi
    
    success "Development deployment completed"
}

deploy_to_prod() {
    log "Deploying to production environment..."
    
    # Copy deployment files to production server
    log "Copying deployment files to production server..."
    rsync -avz --delete \
        -e "ssh -i $SSH_KEY_PATH -o StrictHostKeyChecking=no" \
        "$DEPLOY_DIR/" \
        "$PROD_USER@$PROD_SERVER:/opt/spark-cricket/deploy/"
    
    rsync -avz --delete \
        -e "ssh -i $SSH_KEY_PATH -o StrictHostKeyChecking=no" \
        --exclude=".git" \
        --exclude="node_modules" \
        --exclude=".next" \
        --exclude="coverage" \
        --exclude="test-results" \
        "$PROJECT_ROOT/" \
        "$PROD_USER@$PROD_SERVER:/opt/spark-cricket/"
    
    # Execute deployment on production server
    log "Executing deployment on production server..."
    ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
        "$PROD_USER@$PROD_SERVER" \
        "cd /opt/spark-cricket && ./deploy/scripts/deploy.sh prod --build --restart"
    
    # Wait for services to be healthy
    log "Waiting for services to be healthy..."
    sleep 30
    
    # Check health
    if ! ./deploy/scripts/health-check.sh prod; then
        error "Production health check failed"
        exit 1
    fi
    
    success "Production deployment completed"
}

main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            dev|prod)
                ENVIRONMENT="$1"
                shift
                ;;
            --build)
                BUILD_IMAGES=true
                shift
                ;;
            --restart)
                RESTART_SERVICES=true
                shift
                ;;
            --skip-tests)
                SKIP_TESTS=true
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
    
    # Main deployment flow
    check_environment
    check_dependencies
    setup_ssh_config
    run_tests
    build_images
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        deploy_to_dev
    else
        deploy_to_prod
    fi
    
    log "Deployment completed successfully!"
    log "Application URLs:"
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        echo "  Frontend: http://localhost:3000"
        echo "  Backend:  http://localhost:8080"
    else
        echo "  Frontend: http://15.235.202.148:3000"
        echo "  Backend:  http://15.235.202.148:8080"
    fi
}

# Run main function with all arguments
main "$@"
