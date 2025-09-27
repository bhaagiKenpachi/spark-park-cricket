#!/bin/bash

# Production Environment Monitoring Deployment Script
# Deploys Prometheus and Grafana for cricket.dojima.foundation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENVIRONMENT="prod"
DOMAIN="cricket.dojima.foundation"
PROMETHEUS_PORT="9092"
GRAFANA_PORT="3003"
GRAFANA_PASSWORD="prod123"

# Logging functions
log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    log_success "Docker and Docker Compose are available"
}

# Deploy monitoring stack
deploy_monitoring() {
    log_info "Deploying monitoring stack for $ENVIRONMENT environment..."
    
    # Stop existing containers
    log_info "Stopping existing containers..."
    docker compose down 2>/dev/null || true
    
    # Start monitoring stack
    log_info "Starting monitoring stack..."
    docker compose up -d
    
    # Wait for services to start
    log_info "Waiting for services to start..."
    sleep 30
    
    # Check service health
    check_service_health
}

# Check service health
check_service_health() {
    log_info "Checking service health..."
    
    # Check Prometheus
    if curl -s -f "http://localhost:$PROMETHEUS_PORT/-/healthy" > /dev/null; then
        log_success "Prometheus is healthy"
    else
        log_error "Prometheus health check failed"
        return 1
    fi
    
    # Check Grafana
    if curl -s -f "http://localhost:$GRAFANA_PORT/api/health" > /dev/null; then
        log_success "Grafana is healthy"
    else
        log_error "Grafana health check failed"
        return 1
    fi
}

# Display access information
show_access_info() {
    log_success "Deployment completed successfully!"
    echo ""
    echo "=========================================="
    echo "Production Environment Monitoring"
    echo "=========================================="
    echo "Environment: $ENVIRONMENT"
    echo "Domain: $DOMAIN"
    echo ""
    echo "Prometheus:"
    echo "  Local: http://localhost:$PROMETHEUS_PORT"
    echo "  External: https://$DOMAIN:9090"
    echo ""
    echo "Grafana:"
    echo "  Local: http://localhost:$GRAFANA_PORT"
    echo "  External: https://$DOMAIN:3000"
    echo "  Username: admin"
    echo "  Password: $GRAFANA_PASSWORD"
    echo ""
    echo "Service Management:"
    echo "  Start: docker compose up -d"
    echo "  Stop: docker compose down"
    echo "  Logs: docker compose logs -f"
    echo "  Restart: docker compose restart"
    echo "=========================================="
}

# Main execution
main() {
    log_info "Starting $ENVIRONMENT environment monitoring deployment..."
    
    check_docker
    deploy_monitoring
    show_access_info
}

# Run main function
main "$@"
