#!/bin/bash

# Environment-agnostic monitoring deployment script
# This script deploys Prometheus, Grafana, and Alertmanager to a remote VM

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
MONITORING_DIR="/opt/monitoring"
REMOTE_USER="${REMOTE_USER:-ubuntu}"
REMOTE_HOST="${REMOTE_HOST:-51.79.143.135}"
SSH_KEY="${SSH_KEY:-/Users/luffybhaagi/dojima/tee-auth/terraform/runner_private_key.pem}"
ENVIRONMENT="${ENVIRONMENT:-dev}"

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# Function to check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if SSH key exists
    if [[ ! -f "$SSH_KEY" ]]; then
        error "SSH key not found at $SSH_KEY"
    fi
    
    # Check if SSH key has correct permissions
    if [[ $(stat -c %a "$SSH_KEY" 2>/dev/null || stat -f %A "$SSH_KEY") != "600" ]]; then
        warn "Setting SSH key permissions to 600"
        chmod 600 "$SSH_KEY"
    fi
    
    # Check if required files exist
    local required_files=(
        "$SCRIPT_DIR/../install/install-monitoring.sh"
        "$SCRIPT_DIR/../config/docker-compose.yml"
        "$SCRIPT_DIR/../config/prometheus/prometheus.yml"
        "$SCRIPT_DIR/../config/alertmanager/alertmanager.yml"
        "$SCRIPT_DIR/../config/grafana/datasources/prometheus.yml"
        "$SCRIPT_DIR/../config/grafana/dashboards/dashboard.yml"
        "$SCRIPT_DIR/../config/prometheus/rules/cricket-alerts.yml"
    )
    
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            error "Required file not found: $file"
        fi
    done
    
    log "Prerequisites check passed"
}

# Function to test SSH connection
test_ssh_connection() {
    log "Testing SSH connection to $REMOTE_HOST..."
    
    if ssh -i "$SSH_KEY" -o ConnectTimeout=10 -o BatchMode=yes "$REMOTE_USER@$REMOTE_HOST" "echo 'SSH connection successful'" 2>/dev/null; then
        log "SSH connection successful"
    else
        error "SSH connection failed. Please check your SSH key and host details."
    fi
}

# Function to create environment file
create_env_file() {
    log "Creating environment configuration..."
    
    local env_file="$SCRIPT_DIR/../config/.env"
    
    # Copy template if .env doesn't exist
    if [[ ! -f "$env_file" ]]; then
        cp "$SCRIPT_DIR/../templates/env.template" "$env_file"
        log "Created .env file from template"
    fi
    
    # Update environment-specific values
    sed -i.bak "s/ENVIRONMENT=.*/ENVIRONMENT=$ENVIRONMENT/" "$env_file"
    sed -i.bak "s/ENVIRONMENT_DOMAIN=.*/ENVIRONMENT_DOMAIN=$REMOTE_HOST/" "$env_file"
    
    # Remove backup file
    rm -f "$env_file.bak"
    
    log "Environment configuration updated"
}

# Function to install monitoring stack on remote host
install_monitoring() {
    log "Installing monitoring stack on remote host..."
    
    # Copy installation script
    scp -i "$SSH_KEY" "$SCRIPT_DIR/../install/install-monitoring.sh" "$REMOTE_USER@$REMOTE_HOST:/tmp/"
    
    # Make script executable and run
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        chmod +x /tmp/install-monitoring.sh
        /tmp/install-monitoring.sh
    "
    
    log "Monitoring stack installation completed"
}

# Function to deploy configuration files
deploy_config() {
    log "Deploying configuration files..."
    
    # Create remote monitoring directory
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "sudo mkdir -p $MONITORING_DIR/{prometheus,grafana,alertmanager}"
    
    # Copy configuration files
    scp -i "$SSH_KEY" -r "$SCRIPT_DIR/../config/"* "$REMOTE_USER@$REMOTE_HOST:/tmp/monitoring-config/"
    
    # Move configuration files to proper location
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        sudo cp -r /tmp/monitoring-config/* $MONITORING_DIR/
        sudo chown -R $REMOTE_USER:$REMOTE_USER $MONITORING_DIR
        rm -rf /tmp/monitoring-config
    "
    
    log "Configuration files deployed"
}

# Function to start monitoring services
start_services() {
    log "Starting monitoring services..."
    
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        cd $MONITORING_DIR
        docker-compose down || true
        docker-compose up -d
    "
    
    # Wait for services to start
    log "Waiting for services to start..."
    sleep 30
    
    # Check service status
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "
        cd $MONITORING_DIR
        docker-compose ps
    "
    
    log "Monitoring services started"
}

# Function to verify deployment
verify_deployment() {
    log "Verifying deployment..."
    
    # Check if services are running
    local services=("prometheus" "grafana" "alertmanager")
    local ports=("9090" "3000" "9093")
    
    for i in "${!services[@]}"; do
        local service="${services[$i]}"
        local port="${ports[$i]}"
        
        log "Checking $service on port $port..."
        if ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "curl -s -o /dev/null -w '%{http_code}' http://localhost:$port" | grep -q "200\|302"; then
            log "$service is running and accessible"
        else
            warn "$service might not be running properly"
        fi
    done
    
    log "Deployment verification completed"
}

# Function to show access information
show_access_info() {
    log "Deployment completed successfully!"
    echo ""
    info "Access Information:"
    echo "==================="
    echo "Prometheus: http://$REMOTE_HOST:9090"
    echo "Grafana: http://$REMOTE_HOST:3000 (admin/admin123)"
    echo "Alertmanager: http://$REMOTE_HOST:9093"
    echo ""
    echo "SSH Access:"
    echo "ssh -i $SSH_KEY $REMOTE_USER@$REMOTE_HOST"
    echo ""
    echo "Service Management:"
    echo "cd $MONITORING_DIR && docker-compose [start|stop|restart|logs]"
    echo ""
    echo "Logs:"
    echo "cd $MONITORING_DIR && docker-compose logs -f [service_name]"
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --host HOST        Remote host IP address (default: 51.79.143.135)"
    echo "  -u, --user USER        Remote username (default: ubuntu)"
    echo "  -k, --key KEY          SSH private key path (default: /Users/luffybhaagi/dojima/tee-auth/runner_private_key.pem)"
    echo "  -e, --env ENV          Environment (dev/prod) (default: dev)"
    echo "  --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Deploy to default settings"
    echo "  $0 -h 192.168.1.100 -u admin         # Deploy to custom host"
    echo "  $0 -e prod -k /path/to/key.pem       # Deploy to production"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--host)
            REMOTE_HOST="$2"
            shift 2
            ;;
        -u|--user)
            REMOTE_USER="$2"
            shift 2
            ;;
        -k|--key)
            SSH_KEY="$2"
            shift 2
            ;;
        -e|--env)
            ENVIRONMENT="$2"
            shift 2
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Main deployment function
main() {
    log "Starting monitoring deployment..."
    log "Target: $REMOTE_USER@$REMOTE_HOST"
    log "Environment: $ENVIRONMENT"
    log "SSH Key: $SSH_KEY"
    
    check_prerequisites
    test_ssh_connection
    create_env_file
    install_monitoring
    deploy_config
    start_services
    verify_deployment
    show_access_info
    
    log "Deployment completed successfully!"
}

# Run main function
main "$@"
