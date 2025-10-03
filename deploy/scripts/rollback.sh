#!/bin/bash

# Spark Park Cricket - Rollback Script
# Usage: ./rollback.sh [environment] [options]

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
BACKUP_VERSION=""
FORCE_ROLLBACK=false
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
    echo "  --version VERSION    Rollback to specific version"
    echo "  --force             Force rollback without confirmation"
    echo "  --verbose           Enable verbose output"
    echo "  --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 prod --version v1.2.3"
    echo "  $0 dev --force"
    echo "  $0 prod --verbose"
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

get_available_versions() {
    log "Getting available versions..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        # Check local git tags
        git tag --sort=-version:refname | head -10
    else
        # Check git tags on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket && git tag --sort=-version:refname | head -10"
    fi
}

get_current_version() {
    log "Getting current version..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        git describe --tags --always 2>/dev/null || echo "unknown"
    else
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket && git describe --tags --always 2>/dev/null || echo 'unknown'"
    fi
}

confirm_rollback() {
    if [[ "$FORCE_ROLLBACK" == "true" ]]; then
        return 0
    fi
    
    local current_version="$1"
    local target_version="$2"
    
    echo ""
    warning "Rollback Confirmation"
    echo "======================"
    echo "Environment: $ENVIRONMENT"
    echo "Current version: $current_version"
    echo "Target version: $target_version"
    echo ""
    echo "This will:"
    echo "  1. Stop the current deployment"
    echo "  2. Checkout the target version"
    echo "  3. Rebuild and restart services"
    echo "  4. Run health checks"
    echo ""
    
    read -p "Are you sure you want to proceed? (y/N): " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log "Rollback cancelled by user"
        exit 0
    fi
}

create_backup() {
    log "Creating backup of current deployment..."
    
    local backup_dir="/tmp/spark-cricket-backup-$(date +%Y%m%d-%H%M%S)"
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        # Create local backup
        mkdir -p "$backup_dir"
        
        # Backup current docker compose state
        docker compose -f "$PROJECT_ROOT/docker-compose.yml" config > "$backup_dir/docker-compose-backup.yml"
        
        # Backup current environment
        cp "$PROJECT_ROOT/.env" "$backup_dir/.env" 2>/dev/null || true
        
        log "Backup created at: $backup_dir"
    else
        # Create backup on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "mkdir -p '$backup_dir' && \
             docker compose -f /opt/spark-cricket/deploy/environments/prod/docker-compose.yml config > '$backup_dir/docker-compose-backup.yml' && \
             cp /opt/spark-cricket/.env '$backup_dir/.env' 2>/dev/null || true && \
             echo 'Backup created at: $backup_dir'"
    fi
}

stop_services() {
    log "Stopping current services..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT"
        docker compose down
    else
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket && docker compose -f deploy/environments/prod/docker-compose.yml down"
    fi
    
    success "Services stopped"
}

checkout_version() {
    local target_version="$1"
    
    log "Checking out version: $target_version"
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT"
        
        # Checkout the target version
        if ! git checkout "$target_version"; then
            error "Failed to checkout version $target_version"
            exit 1
        fi
        
        # Pull latest changes if it's a branch
        if git show-ref --verify --quiet refs/remotes/origin/"$target_version"; then
            git pull origin "$target_version"
        fi
    else
        # Checkout on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket && \
             if ! git checkout '$target_version'; then \
                 echo 'Failed to checkout version $target_version'; \
                 exit 1; \
             fi && \
             if git show-ref --verify --quiet refs/remotes/origin/'$target_version'; then \
                 git pull origin '$target_version'; \
             fi"
    fi
    
    success "Checked out version: $target_version"
}

rebuild_and_start() {
    log "Rebuilding and starting services..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        cd "$PROJECT_ROOT"
        
        # Rebuild images
        docker compose build
        
        # Start services
        docker compose up -d
        
        # Wait for services to be healthy
        sleep 15
    else
        # Rebuild and start on production server
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "cd /opt/spark-cricket && \
             docker compose -f deploy/environments/prod/docker-compose.yml build && \
             docker compose -f deploy/environments/prod/docker-compose.yml up -d && \
             sleep 15"
    fi
    
    success "Services rebuilt and started"
}

verify_rollback() {
    log "Verifying rollback..."
    
    # Wait a bit more for services to stabilize
    sleep 10
    
    # Run health checks
    if ! "$SCRIPT_DIR/health-check.sh" "$ENVIRONMENT"; then
        error "Health checks failed after rollback"
        return 1
    fi
    
    success "Rollback verification successful"
}

rollback_to_version() {
    local target_version="$1"
    local current_version
    local available_versions
    
    # Get current version
    current_version=$(get_current_version)
    
    # Get available versions
    available_versions=$(get_available_versions)
    
    # Validate target version
    if [[ -z "$target_version" ]]; then
        error "Target version is required"
        echo "Available versions:"
        echo "$available_versions"
        exit 1
    fi
    
    # Check if target version exists
    if ! echo "$available_versions" | grep -q "^$target_version$"; then
        error "Target version $target_version not found"
        echo "Available versions:"
        echo "$available_versions"
        exit 1
    fi
    
    # Confirm rollback
    confirm_rollback "$current_version" "$target_version"
    
    # Create backup
    create_backup
    
    # Stop services
    stop_services
    
    # Checkout version
    checkout_version "$target_version"
    
    # Rebuild and start
    rebuild_and_start
    
    # Verify rollback
    verify_rollback
    
    success "Rollback to $target_version completed successfully!"
    
    log "Rollback Summary:"
    echo "  Environment: $ENVIRONMENT"
    echo "  Previous version: $current_version"
    echo "  Current version: $target_version"
    echo "  Status: Success"
}

rollback_to_previous() {
    log "Rolling back to previous version..."
    
    # This would require maintaining a deployment history
    # For now, we'll list available versions and let user choose
    warning "Automatic rollback to previous version not implemented"
    warning "Please specify a version using --version option"
    
    echo "Available versions:"
    get_available_versions
}

main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            dev|prod)
                ENVIRONMENT="$1"
                shift
                ;;
            --version)
                BACKUP_VERSION="$2"
                shift 2
                ;;
            --force)
                FORCE_ROLLBACK=true
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
    
    if [[ -n "$BACKUP_VERSION" ]]; then
        rollback_to_version "$BACKUP_VERSION"
    else
        rollback_to_previous
    fi
}

# Run main function with all arguments
main "$@"
