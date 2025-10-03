#!/bin/bash

# Spark Park Cricket - Health Check Script
# Usage: ./health-check.sh [environment]

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
TIMEOUT=30
INTERVAL=5

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
    echo "  --timeout N    Timeout in seconds (default: 30)"
    echo "  --interval N   Check interval in seconds (default: 5)"
    echo "  --help         Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 prod"
    echo "  $0 dev --timeout 60"
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

check_service_health() {
    local service_name="$1"
    local url="$2"
    local expected_status="${3:-200}"
    
    log "Checking $service_name health at $url..."
    
    local response_code
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        response_code=$(curl -s -o /dev/null -w "%{http_code}" "$url" || echo "000")
    else
        response_code=$(ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "curl -s -o /dev/null -w '%{http_code}' '$url'" || echo "000")
    fi
    
    if [[ "$response_code" == "$expected_status" ]]; then
        success "$service_name is healthy (HTTP $response_code)"
        return 0
    else
        error "$service_name is unhealthy (HTTP $response_code)"
        return 1
    fi
}

check_docker_containers() {
    log "Checking Docker containers..."
    
    local containers=("cricket-backend" "cricket-redis" "cricket-frontend")
    local failed_containers=()
    
    for container in "${containers[@]}"; do
        if [[ "$ENVIRONMENT" == "dev" ]]; then
            if docker ps --format "table {{.Names}}" | grep -q "^$container$"; then
                if docker ps --format "table {{.Names}}\t{{.Status}}" | grep "$container" | grep -q "Up"; then
                    success "Container $container is running"
                else
                    error "Container $container is not healthy"
                    failed_containers+=("$container")
                fi
            else
                error "Container $container is not running"
                failed_containers+=("$container")
            fi
        else
            local container_status
            container_status=$(ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
                "$PROD_USER@$PROD_SERVER" \
                "docker ps --format 'table {{.Names}}\t{{.Status}}' | grep '$container' || echo 'not found'")
            
            if [[ "$container_status" == "not found" ]]; then
                error "Container $container is not running on production"
                failed_containers+=("$container")
            elif echo "$container_status" | grep -q "Up"; then
                success "Container $container is running on production"
            else
                error "Container $container is not healthy on production"
                failed_containers+=("$container")
            fi
        fi
    done
    
    if [[ ${#failed_containers[@]} -gt 0 ]]; then
        error "Failed containers: ${failed_containers[*]}"
        return 1
    fi
    
    return 0
}

wait_for_service() {
    local service_name="$1"
    local url="$2"
    local expected_status="${3:-200}"
    local timeout="$TIMEOUT"
    local interval="$INTERVAL"
    
    log "Waiting for $service_name to be healthy (timeout: ${timeout}s)..."
    
    local elapsed=0
    while [[ $elapsed -lt $timeout ]]; do
        if check_service_health "$service_name" "$url" "$expected_status"; then
            return 0
        fi
        
        sleep $interval
        elapsed=$((elapsed + interval))
        
        if [[ $elapsed -lt $timeout ]]; then
            log "Retrying in ${interval}s... (${elapsed}/${timeout}s elapsed)"
        fi
    done
    
    error "$service_name failed to become healthy within ${timeout}s"
    return 1
}

check_application_services() {
    log "Checking application services..."
    
    local services=()
    local failed_services=()
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        services=(
            "Backend:http://localhost:8080/health:200"
            "Frontend:http://localhost:3000:200"
        )
    else
        services=(
            "Backend:http://localhost:8080/health:200"
            "Frontend:http://localhost:3000:200"
        )
    fi
    
    for service in "${services[@]}"; do
        IFS=':' read -r name url expected_status <<< "$service"
        
        if ! check_service_health "$name" "$url" "$expected_status"; then
            failed_services+=("$name")
        fi
    done
    
    if [[ ${#failed_services[@]} -gt 0 ]]; then
        error "Failed services: ${failed_services[*]}"
        return 1
    fi
    
    return 0
}

check_database_connection() {
    log "Checking database connection..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        # Check Redis connection
        if docker exec cricket-redis redis-cli ping | grep -q "PONG"; then
            success "Redis connection is healthy"
        else
            error "Redis connection failed"
            return 1
        fi
    else
        # Check Redis connection on production
        local redis_status
        redis_status=$(ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "docker exec cricket-redis redis-cli ping || echo 'failed'")
        
        if [[ "$redis_status" == "PONG" ]]; then
            success "Redis connection is healthy on production"
        else
            error "Redis connection failed on production"
            return 1
        fi
    fi
    
    return 0
}

get_service_logs() {
    local service_name="$1"
    local lines="${2:-20}"
    
    log "Getting recent logs for $service_name..."
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        docker logs --tail "$lines" "$service_name" 2>&1
    else
        ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no \
            "$PROD_USER@$PROD_SERVER" \
            "docker logs --tail $lines $service_name 2>&1"
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
            --timeout)
                TIMEOUT="$2"
                shift 2
                ;;
            --interval)
                INTERVAL="$2"
                shift 2
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
    
    log "Starting health check for $ENVIRONMENT environment..."
    
    local overall_status=0
    
    # Check Docker containers
    if ! check_docker_containers; then
        overall_status=1
    fi
    
    # Check application services
    if ! check_application_services; then
        overall_status=1
    fi
    
    # Check database connections
    if ! check_database_connection; then
        overall_status=1
    fi
    
    if [[ $overall_status -eq 0 ]]; then
        success "All health checks passed!"
        log "Application is healthy and ready to serve requests"
    else
        error "Some health checks failed!"
        log "Getting recent logs for troubleshooting..."
        
        get_service_logs "cricket-backend" 50
        echo ""
        get_service_logs "cricket-redis" 20
    fi
    
    exit $overall_status
}

# Run main function with all arguments
main "$@"
