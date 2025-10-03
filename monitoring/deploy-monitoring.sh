#!/bin/bash

# Unified Monitoring Deployment Script
# Deploys monitoring stack with Loki support to monitoring instance
# Handles both dev and prod environments with proper error handling

set -e

# Configuration
MONITORING_HOST="${SPARK_CRICKET_MONITORING_SERVER:-51.79.143.135}"
MONITORING_USER="ubuntu"
SSH_KEY="${SPARK_CRICKET_SSH_KEY_MONITORING:-$HOME/.ssh/spark-cricket-monitoring}"
MONITORING_PATH="/opt/monitoring"

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

usage() {
    echo "Usage: $0 [environment] [options]"
    echo ""
    echo "Environments:"
    echo "  dev     - Development environment (ports: 3001, 9091, 3101)"
    echo "  prod    - Production environment (ports: 3003, 9092, 3102)"
    echo ""
    echo "Options:"
    echo "  --setup-log-forwarding    Also setup log forwarding from server instance"
    echo "  --force                  Force deployment even if services are running"
    echo "  --help                   Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev"
    echo "  $0 prod --setup-log-forwarding"
    echo "  $0 prod --force"
}

# Parse arguments
ENVIRONMENT=""
SETUP_LOG_FORWARDING=false
FORCE_DEPLOYMENT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        dev|prod)
            ENVIRONMENT="$1"
            shift
            ;;
        --setup-log-forwarding)
            SETUP_LOG_FORWARDING=true
            shift
            ;;
        --force)
            FORCE_DEPLOYMENT=true
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

if [[ -z "$ENVIRONMENT" ]]; then
    log_error "Environment is required"
    usage
fi

if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" ]]; then
    log_error "Invalid environment. Must be 'dev' or 'prod'"
fi

# Test SSH connection
test_ssh_connection() {
    log "Testing SSH connection to monitoring instance..."
    if ! ssh -i "$SSH_KEY" -o ConnectTimeout=10 -o BatchMode=yes "$MONITORING_USER@$MONITORING_HOST" "echo 'SSH connection successful'" >/dev/null 2>&1; then
        log_error "Cannot connect to monitoring instance $MONITORING_HOST"
    fi
    log_success "SSH connection established"
}

# Install Docker and Docker Compose if needed
install_docker() {
    log "Checking Docker installation on monitoring instance..."
    
    if ! ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "command -v docker >/dev/null 2>&1"; then
        log "Installing Docker on monitoring instance..."
        ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "
            curl -fsSL https://get.docker.com -o get-docker.sh &&
            sudo sh get-docker.sh &&
            sudo systemctl start docker &&
            sudo systemctl enable docker &&
            sudo usermod -aG docker ubuntu &&
            rm get-docker.sh
        "
        log_success "Docker installed successfully"
    else
        log "Docker is already installed"
    fi

    # Install Docker Compose if not present
    log "Checking Docker Compose installation..."
    if ! ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1"; then
        log "Installing Docker Compose..."
        ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "
            sudo apt-get update &&
            sudo apt-get install -y docker-compose-plugin
        "
        log_success "Docker Compose installed successfully"
    else
        log "Docker Compose is already installed"
    fi
}

# Create directory structure
create_directories() {
    log "Creating directory structure on monitoring instance..."
    ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "
        sudo mkdir -p $MONITORING_PATH/environments/$ENVIRONMENT/{grafana/{dashboards,datasources},prometheus/{rules},loki,promtail}
    "
    log_success "Directory structure created"
}

# Copy configuration files
copy_configurations() {
    log "Copying monitoring configuration files..."
    
    # Copy the entire environment directory
    rsync -avz -e "ssh -i $SSH_KEY" \
        monitoring/environments/$ENVIRONMENT/ \
        "$MONITORING_USER@$MONITORING_HOST:$MONITORING_PATH/environments/$ENVIRONMENT/"
    
    log_success "Configuration files copied"
}

# Stop existing containers
stop_existing_containers() {
    log "Stopping existing monitoring containers..."
    ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "
        cd $MONITORING_PATH/environments/$ENVIRONMENT &&
        sudo docker compose down || true
    "
    log_success "Existing containers stopped"
}

# Start monitoring services
start_monitoring_services() {
    log "Starting monitoring services..."
    ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "
        cd $MONITORING_PATH/environments/$ENVIRONMENT &&
        sudo docker compose up -d
    "
    log_success "Monitoring services started"
}

# Wait for services to be ready
wait_for_services() {
    log "Waiting for services to be ready..."
    sleep 30
    
    # Check service status
    log "Checking service status..."
    ssh -i "$SSH_KEY" "$MONITORING_USER@$MONITORING_HOST" "
        cd $MONITORING_PATH/environments/$ENVIRONMENT &&
        sudo docker compose ps
    "
}

# Setup log forwarding if requested
setup_log_forwarding() {
    if [[ "$SETUP_LOG_FORWARDING" == "true" ]]; then
        log "Setting up log forwarding from server instance..."
        
        # Run the log forwarding setup script
        if [[ -f "monitoring/setup-log-forwarding.sh" ]]; then
            ./monitoring/setup-log-forwarding.sh "$ENVIRONMENT"
        else
            log_warn "Log forwarding script not found. Please run setup-log-forwarding.sh separately."
        fi
    fi
}

# Show deployment summary
show_summary() {
    local grafana_port
    local prometheus_port
    local loki_port
    
    if [[ "$ENVIRONMENT" == "dev" ]]; then
        grafana_port="3001"
        prometheus_port="9091"
        loki_port="3101"
    else
        grafana_port="3003"
        prometheus_port="9092"
        loki_port="3102"
    fi
    
    log "Monitoring deployment completed successfully!"
    log ""
    log "Access URLs:"
    log "  Grafana:     http://$MONITORING_HOST:$grafana_port"
    log "  Prometheus:  http://$MONITORING_HOST:$prometheus_port"
    log "  Loki:        http://$MONITORING_HOST:$loki_port"
    log ""
    log "Default credentials:"
    log "  Grafana: admin / ${ENVIRONMENT}123"
    log ""
    log "Dashboard URLs:"
    log "  Database Performance: http://$MONITORING_HOST:$grafana_port/d/cricket-db-performance"
    log "  Backend Logs:         http://$MONITORING_HOST:$grafana_port/d/cricket-backend-logs-$ENVIRONMENT"
}

# Main execution
main() {
    log "Starting monitoring deployment for environment: $ENVIRONMENT"
    
    test_ssh_connection
    install_docker
    create_directories
    copy_configurations
    stop_existing_containers
    start_monitoring_services
    wait_for_services
    setup_log_forwarding
    show_summary
    
    log_success "Monitoring deployment completed for $ENVIRONMENT environment!"
}

# Run main function
main "$@"
