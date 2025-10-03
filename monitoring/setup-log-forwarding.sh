#!/bin/bash

# Unified Log Forwarding Setup Script
# Sets up Promtail on server instance to forward logs to Loki on monitoring instance
# Handles both dev and prod environments with proper error handling

set -e

# Configuration
SERVER_HOST="${SPARK_CRICKET_PROD_SERVER:-15.235.202.148}"
SERVER_USER="ubuntu"
SSH_KEY="${SPARK_CRICKET_SSH_KEY_PROD:-$HOME/.ssh/spark-cricket-prod}"
MONITORING_HOST="${SPARK_CRICKET_MONITORING_SERVER:-51.79.143.135}"

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
    echo "  dev     - Development environment (Loki port: 3101)"
    echo "  prod    - Production environment (Loki port: 3102)"
    echo ""
    echo "Options:"
    echo "  --force    Force setup even if Promtail is already running"
    echo "  --remove   Remove existing Promtail setup"
    echo "  --help     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 prod"
    echo "  $0 dev --force"
    echo "  $0 prod --remove"
}

# Parse arguments
ENVIRONMENT=""
FORCE_SETUP=false
REMOVE_SETUP=false

while [[ $# -gt 0 ]]; do
    case $1 in
        dev|prod)
            ENVIRONMENT="$1"
            shift
            ;;
        --force)
            FORCE_SETUP=true
            shift
            ;;
        --remove)
            REMOVE_SETUP=true
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

# Set Loki port based on environment
if [[ "$ENVIRONMENT" == "dev" ]]; then
    LOKI_PORT="3101"
else
    LOKI_PORT="3102"
fi

# Test SSH connection to server instance
test_server_connection() {
    log "Testing SSH connection to server instance..."
    if ! ssh -i "$SSH_KEY" -o ConnectTimeout=10 -o BatchMode=yes "$SERVER_USER@$SERVER_HOST" "echo 'SSH connection successful'" >/dev/null 2>&1; then
        log_error "Cannot connect to server instance $SERVER_HOST"
    fi
    log_success "SSH connection to server established"
}

# Test connectivity to monitoring instance
test_monitoring_connection() {
    log "Testing connectivity to monitoring instance..."
    if ! ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "nc -z $MONITORING_HOST $LOKI_PORT" >/dev/null 2>&1; then
        log_warn "Cannot connect to Loki on monitoring instance. Make sure Loki is running on $MONITORING_HOST:$LOKI_PORT"
        log "You may need to configure firewall rules or check Loki service status"
    else
        log_success "Connectivity to monitoring instance verified"
    fi
}

# Remove existing Promtail setup
remove_promtail() {
    log "Removing existing Promtail setup..."
    
    # Stop and remove Promtail container
    ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "
        sudo docker stop promtail-server 2>/dev/null || true &&
        sudo docker rm promtail-server 2>/dev/null || true
    "
    
    # Remove Promtail configuration
    ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "
        sudo rm -rf /opt/promtail 2>/dev/null || true
    "
    
    log_success "Existing Promtail setup removed"
}

# Check if Promtail is already running
check_existing_promtail() {
    if ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "sudo docker ps | grep promtail-server" >/dev/null 2>&1; then
        if [[ "$FORCE_SETUP" == "false" ]]; then
            log_warn "Promtail is already running on server instance"
            log "Use --force to reinstall or --remove to remove existing setup"
            exit 0
        else
            log "Force setup requested, removing existing Promtail..."
            remove_promtail
        fi
    fi
}

# Create Promtail configuration
create_promtail_config() {
    log "Creating Promtail configuration..."
    
    local promtail_config=$(cat <<EOF
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://${MONITORING_HOST}:${LOKI_PORT}/loki/api/v1/push

scrape_configs:
  # Collect logs from backend container
  - job_name: cricket-backend-${ENVIRONMENT}
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
        filters:
          - name: name
            values: ["cricket-backend-${ENVIRONMENT}"]
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/?(.*)'
        target_label: 'container_name'
      - source_labels: ['__meta_docker_container_log_stream']
        target_label: 'logstream'
      - source_labels: ['__meta_docker_container_id']
        target_label: 'container_id'
      - target_label: 'environment'
        replacement: '${ENVIRONMENT}'
      - target_label: 'instance'
        replacement: '${SERVER_HOST}'
    pipeline_stages:
      - json:
          expressions:
            output: log
            stream: stream
            attrs:
      - json:
          expressions:
            tag:
          source: attrs
      - regex:
          expression: (?P<container_name>(?:[^|]*))\|
          source: tag
      - timestamp:
          format: RFC3339Nano
          source: time
      - labels:
          stream:
          container_name:
          environment:
          instance:
      - output:
          source: output

  # Collect logs from frontend container
  - job_name: cricket-frontend-${ENVIRONMENT}
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
        filters:
          - name: name
            values: ["cricket-frontend-${ENVIRONMENT}"]
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/?(.*)'
        target_label: 'container_name'
      - source_labels: ['__meta_docker_container_log_stream']
        target_label: 'logstream'
      - source_labels: ['__meta_docker_container_id']
        target_label: 'container_id'
      - target_label: 'environment'
        replacement: '${ENVIRONMENT}'
      - target_label: 'instance'
        replacement: '${SERVER_HOST}'
    pipeline_stages:
      - json:
          expressions:
            output: log
            stream: stream
            attrs:
      - json:
          expressions:
            tag:
          source: attrs
      - regex:
          expression: (?P<container_name>(?:[^|]*))\|
          source: tag
      - timestamp:
          format: RFC3339Nano
          source: time
      - labels:
          stream:
          container_name:
          environment:
          instance:
      - output:
          source: output
EOF
)
    
    # Create Promtail directory and config
    ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "
        sudo mkdir -p /opt/promtail
    "
    
    echo "$promtail_config" | ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "sudo tee /opt/promtail/promtail.yml > /dev/null"
    
    log_success "Promtail configuration created"
}

# Deploy Promtail container
deploy_promtail() {
    log "Deploying Promtail container..."
    
    ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "
        sudo docker run -d \
            --name promtail-server \
            --restart unless-stopped \
            -v /opt/promtail/promtail.yml:/etc/promtail/config.yml:ro \
            -v /var/log:/var/log:ro \
            -v /var/lib/docker/containers:/var/lib/docker/containers:ro \
            -v /var/run/docker.sock:/var/run/docker.sock:ro \
            grafana/promtail:latest \
            -config.file=/etc/promtail/config.yml
    "
    
    log_success "Promtail container deployed"
}

# Verify Promtail is working
verify_promtail() {
    log "Verifying Promtail is working..."
    
    # Wait for Promtail to start
    sleep 10
    
    # Check if container is running
    if ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "sudo docker ps | grep promtail-server" >/dev/null 2>&1; then
        log_success "Promtail container is running"
    else
        log_error "Promtail container failed to start"
    fi
    
    # Check Promtail logs
    log "Checking Promtail logs..."
    ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "sudo docker logs promtail-server --tail=5"
    
    # Generate test logs
    log "Generating test logs..."
    ssh -i "$SSH_KEY" "$SERVER_USER@$SERVER_HOST" "
        sudo docker logs --tail=3 cricket-backend-${ENVIRONMENT} 2>&1
    "
    
    log_success "Promtail verification completed"
}

# Show setup summary
show_summary() {
    log "Log forwarding setup completed successfully!"
    log ""
    log "Configuration Summary:"
    log "  Environment: $ENVIRONMENT"
    log "  Server Instance: $SERVER_HOST"
    log "  Monitoring Instance: $MONITORING_HOST"
    log "  Loki Port: $LOKI_PORT"
    log ""
    log "Log Sources:"
    log "  - cricket-backend-${ENVIRONMENT}"
    log "  - cricket-frontend-${ENVIRONMENT}"
    log ""
    log "Management Commands:"
    log "  View Promtail logs: ssh -i $SSH_KEY $SERVER_USER@$SERVER_HOST 'sudo docker logs promtail-server -f'"
    log "  Restart Promtail: ssh -i $SSH_KEY $SERVER_USER@$SERVER_HOST 'sudo docker restart promtail-server'"
    log "  Remove Promtail: $0 $ENVIRONMENT --remove"
    log ""
    log "Logs are now being forwarded to Loki at: http://$MONITORING_HOST:$LOKI_PORT"
    log "View logs in Grafana at the Backend Logs dashboard"
}

# Main execution
main() {
    log "Setting up log forwarding for $ENVIRONMENT environment..."
    
    test_server_connection
    test_monitoring_connection
    
    if [[ "$REMOVE_SETUP" == "true" ]]; then
        remove_promtail
        log_success "Promtail setup removed"
        exit 0
    fi
    
    check_existing_promtail
    create_promtail_config
    deploy_promtail
    verify_promtail
    show_summary
    
    log_success "Log forwarding setup completed successfully!"
}

# Run main function
main "$@"
