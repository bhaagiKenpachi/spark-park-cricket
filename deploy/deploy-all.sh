#!/bin/bash

# Spark Park Cricket - Master Deployment Script
# Deploys both application and monitoring with proper error handling and security validation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log() { echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

usage() {
    echo "Spark Park Cricket - Master Deployment Script"
    echo ""
    echo "Usage: $0 [command] [environment] [options]"
    echo ""
    echo "Commands:"
    echo "  deploy-app      Deploy application only"
    echo "  deploy-mon      Deploy monitoring only"
    echo "  deploy-all      Deploy both application and monitoring"
    echo "  health-check    Run health checks"
    echo "  setup           Initial setup and security check"
    echo ""
    echo "Environments:"
    echo "  dev             Development environment"
    echo "  prod            Production environment"
    echo ""
    echo "Required Environment Variables:"
    echo "  GRAFANA_ADMIN_PASSWORD_PROD         Grafana admin password for production"
    echo "  GRAFANA_ADMIN_PASSWORD_DEV          Grafana admin password for development"
    echo ""
    echo "Security Options:"
    echo "  --with-logs     Include log forwarding setup"
    echo "  --force         Force deployment even if services are running"
    echo "  --clean         Clean build (remove existing images)"
    echo "  --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  export GRAFANA_ADMIN_PASSWORD_PROD='SecurePassword123!'"
    echo "  $0 deploy-all prod --with-logs"
    echo "  $0 setup"
}

# Security validation
validate_security() {
    log "Validating security configuration..."
    
    # Check for default passwords in environment variables
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        if [[ "${GRAFANA_ADMIN_PASSWORD_PROD:-}" == "prod123" || "${GRAFANA_ADMIN_PASSWORD_PROD:-}" == "admin" ]]; then
            log_error "Default password detected for production Grafana. Please set GRAFANA_ADMIN_PASSWORD_PROD environment variable with a strong password."
        fi
        if [[ -z "${GRAFANA_ADMIN_PASSWORD_PROD:-}" ]]; then
            log_warn "GRAFANA_ADMIN_PASSWORD_PROD not set. Using default password (CHANGE_ME_NOW). Please set this environment variable."
        fi
    fi
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        if [[ "${GRAFANA_ADMIN_PASSWORD_DEV:-}" == "dev123" || "${GRAFANA_ADMIN_PASSWORD_DEV:-}" == "admin" ]]; then
            log_error "Default password detected for development Grafana. Please set GRAFANA_ADMIN_PASSWORD_DEV environment variable with a strong password."
        fi
        if [[ -z "${GRAFANA_ADMIN_PASSWORD_DEV:-}" ]]; then
            log_warn "GRAFANA_ADMIN_PASSWORD_DEV not set. Using default password (CHANGE_ME_NOW). Please set this environment variable."
        fi
    fi
    
    log_success "Security validation completed"
}

# Parse arguments
COMMAND=""
ENVIRONMENT=""
WITH_LOGS=false
FORCE=false
CLEAN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        deploy-app|deploy-mon|deploy-all|health-check|setup)
            COMMAND="$1"
            shift
            ;;
        dev|prod)
            ENVIRONMENT="$1"
            shift
            ;;
        --with-logs)
            WITH_LOGS=true
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        --help)
            usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            ;;
    esac
done

# Validate command and environment
if [[ -z "$COMMAND" ]]; then
    log_error "Command is required"
    usage
fi

if [[ "$COMMAND" != "setup" && -z "$ENVIRONMENT" ]]; then
    log_error "Environment is required for command: $COMMAND"
    usage
fi

if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" && "$COMMAND" != "setup" ]]; then
    log_error "Invalid environment. Must be 'dev' or 'prod'"
fi

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker compose &> /dev/null; then
        log_error "Docker Compose is not installed"
    fi
    
    # Check if required directories exist
    if [[ ! -d "scripts" ]]; then
        log_error "scripts directory not found"
    fi
    
    if [[ ! -d "../monitoring" ]]; then
        log_error "monitoring directory not found"
    fi
    
    log_success "Prerequisites check passed"
}

# Setup SSH keys and configuration
setup_ssh() {
    log "Setting up SSH configuration..."
    
    if [[ -f "scripts/setup-ssh.sh" ]]; then
        ./scripts/setup-ssh.sh
        log_success "SSH setup completed"
    else
        log_warn "SSH setup script not found"
    fi
}

# Deploy application
deploy_application() {
    log "Deploying application to $ENVIRONMENT environment..."
    
    if [[ "$CLEAN" == "true" ]]; then
        log "Building with clean option..."
        ./scripts/build.sh "$ENVIRONMENT" --clean
    fi
    
    ./scripts/deploy.sh "$ENVIRONMENT"
    log_success "Application deployment completed"
}

# Deploy monitoring
deploy_monitoring() {
    log "Deploying monitoring to $ENVIRONMENT environment..."
    
    if [[ "$WITH_LOGS" == "true" ]]; then
        "$PROJECT_ROOT/monitoring/deploy-monitoring.sh" "$ENVIRONMENT" --setup-log-forwarding
    else
        "$PROJECT_ROOT/monitoring/deploy-monitoring.sh" "$ENVIRONMENT"
    fi
    
    log_success "Monitoring deployment completed"
}

# Deploy both application and monitoring
deploy_all() {
    log "Deploying complete stack to $ENVIRONMENT environment..."
    
    deploy_application
    
    # Wait a bit for application to stabilize
    sleep 10
    
    deploy_monitoring
    
    log_success "Complete stack deployment completed"
}

# Run health checks
run_health_checks() {
    log "Running health checks for $ENVIRONMENT environment..."
    
    if [[ -f "scripts/health-check.sh" ]]; then
        ./scripts/health-check.sh "$ENVIRONMENT"
        log_success "Health checks completed"
    else
        log_warn "Health check script not found"
    fi
}

# Rollback deployment
rollback_deployment() {
    log "Rolling back $ENVIRONMENT environment..."
    
    if [[ -f "scripts/rollback.sh" ]]; then
        echo "Available versions:"
        ./scripts/rollback.sh "$ENVIRONMENT" --help
        log_warn "Please specify a version to rollback to"
        log "Usage: ./scripts/rollback.sh $ENVIRONMENT --version VERSION"
    else
        log_warn "Rollback script not found"
    fi
}

# Show deployment summary
show_summary() {
    log "Deployment Summary"
    echo "=================="
    echo "Command: $COMMAND"
    echo "Environment: $ENVIRONMENT"
    echo "With Logs: $WITH_LOGS"
    echo "Force: $FORCE"
    echo "Clean: $CLEAN"
    echo ""
    
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        echo "Production URLs:"
        echo "  Application: http://15.235.202.148"
        echo "  Health Check: http://15.235.202.148:8080/health"
        echo "  Grafana: http://51.79.143.135:3003 (admin/[PASSWORD])"
        echo "  Prometheus: http://51.79.143.135:9092"
        echo "  Loki: http://51.79.143.135:3102"
    elif [[ "$ENVIRONMENT" == "dev" ]]; then
        echo "Development URLs:"
        echo "  Application: http://localhost:3000"
        echo "  Health Check: http://localhost:8080/health"
        echo "  Grafana: http://51.79.143.135:3001 (admin/[PASSWORD])"
        echo "  Prometheus: http://51.79.143.135:9091"
        echo "  Loki: http://51.79.143.135:3101"
    fi
}

# Main execution
main() {
    log "Starting Spark Park Cricket deployment..."
    
    # Change to deploy directory
    cd "$SCRIPT_DIR"
    
    check_prerequisites
    validate_security
    
    case "$COMMAND" in
        setup)
            setup_ssh
            ;;
        deploy-app)
            deploy_application
            run_health_checks
            ;;
        deploy-mon)
            deploy_monitoring
            ;;
        deploy-all)
            deploy_all
            run_health_checks
            ;;
        health-check)
            run_health_checks
            ;;
        rollback)
            rollback_deployment
            ;;
        *)
            log_error "Unknown command: $COMMAND"
            ;;
    esac
    
    show_summary
    log_success "Deployment completed successfully!"
}

# Run main function
main "$@"