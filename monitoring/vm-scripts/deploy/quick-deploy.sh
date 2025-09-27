#!/bin/bash

# Quick deployment script for monitoring stack
# This script provides a simplified interface for common deployment scenarios

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SSH_KEY="/Users/luffybhaagi/dojima/tee-auth/terraform/runner_private_key.pem"
REMOTE_HOST="51.79.143.135"
REMOTE_USER="ubuntu"

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Function to deploy to development environment
deploy_dev() {
    log "Deploying to development environment..."
    "$SCRIPT_DIR/deploy-monitoring.sh" -h "$REMOTE_HOST" -u "$REMOTE_USER" -k "$SSH_KEY" -e dev
}

# Function to deploy to production environment
deploy_prod() {
    log "Deploying to production environment..."
    warn "Production deployment requires additional configuration!"
    read -p "Are you sure you want to deploy to production? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        "$SCRIPT_DIR/deploy-monitoring.sh" -h "$REMOTE_HOST" -u "$REMOTE_USER" -k "$SSH_KEY" -e prod
    else
        log "Production deployment cancelled"
    fi
}

# Function to check service status
check_status() {
    log "Checking service status on $REMOTE_HOST..."
    
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        cd /opt/monitoring
        echo '=== Docker Compose Status ==='
        docker-compose ps
        echo ''
        echo '=== Service Health Checks ==='
        echo 'Prometheus:'
        curl -s -o /dev/null -w 'HTTP Status: %{http_code}\n' http://localhost:9090 || echo 'Not accessible'
        echo 'Grafana:'
        curl -s -o /dev/null -w 'HTTP Status: %{http_code}\n' http://localhost:3000 || echo 'Not accessible'
        echo 'Alertmanager:'
        curl -s -o /dev/null -w 'HTTP Status: %{http_code}\n' http://localhost:9093 || echo 'Not accessible'
    "
}

# Function to view logs
view_logs() {
    local service="${1:-all}"
    
    if [[ "$service" == "all" ]]; then
        log "Viewing logs for all services..."
        ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
            cd /opt/monitoring
            docker-compose logs --tail=50
        "
    else
        log "Viewing logs for $service..."
        ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
            cd /opt/monitoring
            docker-compose logs --tail=50 $service
        "
    fi
}

# Function to restart services
restart_services() {
    log "Restarting monitoring services..."
    
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        cd /opt/monitoring
        docker-compose restart
    "
    
    log "Services restarted"
}

# Function to stop services
stop_services() {
    log "Stopping monitoring services..."
    
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        cd /opt/monitoring
        docker-compose down
    "
    
    log "Services stopped"
}

# Function to update configuration
update_config() {
    log "Updating configuration..."
    
    # Deploy only configuration files
    scp -i "$SSH_KEY" -r "$SCRIPT_DIR/../config/"* "$REMOTE_USER@$REMOTE_HOST:/tmp/monitoring-config/"
    
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        sudo cp -r /tmp/monitoring-config/* /opt/monitoring/
        sudo chown -R $REMOTE_USER:$REMOTE_USER /opt/monitoring
        rm -rf /tmp/monitoring-config
        cd /opt/monitoring
        docker-compose restart
    "
    
    log "Configuration updated and services restarted"
}

# Function to show help
show_help() {
    echo "Quick Monitoring Deployment Script"
    echo "================================="
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  deploy-dev     Deploy to development environment"
    echo "  deploy-prod    Deploy to production environment"
    echo "  status         Check service status"
    echo "  logs [service] View logs (service: prometheus, grafana, alertmanager, or all)"
    echo "  restart        Restart all services"
    echo "  stop           Stop all services"
    echo "  update-config  Update configuration files"
    echo "  help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy-dev              # Deploy to dev environment"
    echo "  $0 status                  # Check service status"
    echo "  $0 logs prometheus         # View Prometheus logs"
    echo "  $0 logs                    # View all logs"
    echo "  $0 restart                 # Restart services"
    echo ""
    echo "Access URLs:"
    echo "  Prometheus: http://$REMOTE_HOST:9090"
    echo "  Grafana: http://$REMOTE_HOST:3000 (admin/admin123)"
    echo "  Alertmanager: http://$REMOTE_HOST:9093"
}

# Main function
main() {
    local command="${1:-help}"
    
    case "$command" in
        deploy-dev)
            deploy_dev
            ;;
        deploy-prod)
            deploy_prod
            ;;
        status)
            check_status
            ;;
        logs)
            view_logs "$2"
            ;;
        restart)
            restart_services
            ;;
        stop)
            stop_services
            ;;
        update-config)
            update_config
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
