#!/bin/bash

# VM Deployment Script for Dev and Prod Environments
# Uses TOML configuration for environment-specific settings

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$(dirname "$SCRIPT_DIR")/config"
CONFIG_FILE="$CONFIG_DIR/monitoring.toml"

# Check if TOML config file exists
if [[ ! -f "$CONFIG_FILE" ]]; then
    echo -e "${RED}Error: Configuration file not found: $CONFIG_FILE${NC}"
    exit 1
fi

# Function to parse TOML config (simple parser)
parse_toml() {
    local section="$1"
    local key="$2"
    local file="$3"
    
    # Extract value from TOML file
    awk -F' = ' '
    /^\['"$section"'\]/ { in_section=1; next }
    /^\[/ { in_section=0; next }
    in_section && /^'"$key"' = / { 
        gsub(/^"|"$/, "", $2)
        print $2
        exit
    }
    ' "$file"
}

# Function to get global config
get_global_config() {
    local key="$1"
    parse_toml "global" "$key" "$CONFIG_FILE"
}

# Function to get environment config
get_env_config() {
    local env="$1"
    local key="$2"
    parse_toml "$env" "$key" "$CONFIG_FILE"
}

# Function to check if environment is enabled
is_env_enabled() {
    local env="$1"
    local enabled=$(get_env_config "$env" "enabled")
    [[ "$enabled" == "true" ]]
}

# Logging functions
log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] SUCCESS: $1${NC}"
}

log_warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
}

# Function to deploy single environment
deploy_environment() {
    local env="$1"
    
    if ! is_env_enabled "$env"; then
        log_warn "Environment '$env' is disabled, skipping..."
        return 0
    fi
    
    log_info "Deploying $env environment..."
    
    # Get environment-specific config
    local name=$(get_env_config "$env" "name")
    local domain=$(get_env_config "$env" "domain")
    local backend_url=$(get_env_config "$env" "backend_url")
    local prometheus_port=$(get_env_config "$env" "prometheus_port")
    local grafana_port=$(get_env_config "$env" "grafana_port")
    local grafana_password=$(get_env_config "$env" "grafana_password")
    local data_retention=$(get_env_config "$env" "data_retention")
    local scrape_interval=$(get_env_config "$env" "scrape_interval")
    local container_prefix=$(get_env_config "$env" "container_prefix")
    
    # Get global config
    local vm_host=$(get_global_config "vm_host")
    local ssh_key=$(get_global_config "ssh_key")
    local ssh_user=$(get_global_config "ssh_user")
    local monitoring_dir=$(get_global_config "monitoring_dir")
    
    log_info "Environment: $name"
    log_info "Domain: $domain"
    log_info "Backend URL: $backend_url"
    log_info "Prometheus Port: $prometheus_port"
    log_info "Grafana Port: $grafana_port"
    
    # Create environment-specific directory structure on VM
    log_info "Creating directory structure on VM..."
    ssh -i "$ssh_key" "$ssh_user@$vm_host" "
        mkdir -p $monitoring_dir/environments/$name/{prometheus,grafana/{datasources,dashboards}}
        mkdir -p $monitoring_dir/environments/$name/prometheus/rules
    "
    
    # Generate environment-specific Prometheus config
    log_info "Generating Prometheus configuration for $name..."
    cat > /tmp/prometheus_${name}.yml << EOF
global:
  scrape_interval: ${scrape_interval}
  evaluation_interval: 15s
  external_labels:
    environment: '${name}'
    cluster: 'cricket-${name}-cluster'

scrape_configs:
  # Prometheus itself
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  # Backend metrics
  - job_name: 'cricket-backend-${name}'
    static_configs:
      - targets: ['${domain}:443']
    metrics_path: '/metrics'
    scheme: 'https'
    scrape_interval: ${scrape_interval}
    scrape_timeout: 10s
    tls_config:
      insecure_skip_verify: true

  # Frontend metrics (if available)
  - job_name: 'cricket-frontend-${name}'
    static_configs:
      - targets: ['${domain}:443']
    metrics_path: '/api/metrics'
    scheme: 'https'
    scrape_interval: 30s
    scrape_timeout: 10s
    tls_config:
      insecure_skip_verify: true

# Alerting rules
rule_files:
  - "rules/*.yml"

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
EOF

    # Generate environment-specific Grafana datasource config
    log_info "Generating Grafana datasource configuration for $name..."
    cat > /tmp/datasource_${name}.yml << EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
    jsonData:
      httpMethod: POST
      manageAlerts: true
      prometheusType: Prometheus
      prometheusVersion: 2.40.0
      cacheLevel: 'High'
      disableRecordingRules: false
      incrementalQueryOverlapWindow: 10m
      queryTimeout: 60s
      timeInterval: ${scrape_interval}
    version: 1
EOF

    # Generate environment-specific Grafana dashboard config
    log_info "Generating Grafana dashboard configuration for $name..."
    cat > /tmp/dashboard_${name}.yml << EOF
apiVersion: 1

providers:
  - name: 'cricket-dashboards'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
EOF

    # Generate environment-specific Docker Compose config
    log_info "Generating Docker Compose configuration for $name..."
    cat > /tmp/docker-compose_${name}.yml << EOF
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus-${container_prefix}
    restart: unless-stopped
    ports:
      - "${prometheus_port}:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - ./prometheus/rules:/etc/prometheus/rules:ro
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=${data_retention}'
      - '--web.enable-lifecycle'
      - '--web.enable-admin-api'
      - '--web.enable-remote-write-receiver'
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:latest
    container_name: grafana-${container_prefix}
    restart: unless-stopped
    ports:
      - "${grafana_port}:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/datasources:/etc/grafana/provisioning/datasources:ro
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards:ro
      - ./grafana/dashboards:/var/lib/grafana/dashboards:ro
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${grafana_password}
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_DOMAIN=${domain}
      - GF_SERVER_ROOT_URL=https://${domain}:${grafana_port}
      - GF_INSTALL_PLUGINS=grafana-piechart-panel,grafana-worldmap-panel
    networks:
      - monitoring

volumes:
  prometheus_data:
    driver: local
  grafana_data:
    driver: local

networks:
  monitoring:
    driver: bridge
EOF

    # Copy configuration files to VM
    log_info "Copying configuration files to VM..."
    scp -i "$ssh_key" /tmp/prometheus_${name}.yml "$ssh_user@$vm_host:$monitoring_dir/environments/$name/prometheus/prometheus.yml"
    scp -i "$ssh_key" /tmp/datasource_${name}.yml "$ssh_user@$vm_host:$monitoring_dir/environments/$name/grafana/datasources/prometheus.yml"
    scp -i "$ssh_key" /tmp/dashboard_${name}.yml "$ssh_user@$vm_host:$monitoring_dir/environments/$name/grafana/dashboards/dashboard.yml"
    scp -i "$ssh_key" /tmp/docker-compose_${name}.yml "$ssh_user@$vm_host:$monitoring_dir/environments/$name/docker-compose.yml"
    
    # Copy dashboard JSON file
    if [[ -f "$CONFIG_DIR/../grafana/dashboards/cricket-db-performance.json" ]]; then
        scp -i "$ssh_key" "$CONFIG_DIR/../grafana/dashboards/cricket-db-performance.json" "$ssh_user@$vm_host:$monitoring_dir/environments/$name/grafana/dashboards/"
    fi
    
    # Deploy on VM
    log_info "Deploying $name environment on VM..."
    ssh -i "$ssh_key" "$ssh_user@$vm_host" "
        cd $monitoring_dir/environments/$name
        
        # Stop existing containers
        docker compose down 2>/dev/null || true
        
        # Start monitoring stack
        docker compose up -d
        
        # Wait for services to start
        sleep 30
        
        # Check service health
        if curl -s -f http://localhost:${prometheus_port}/api/v1/status/config > /dev/null; then
            echo 'Prometheus is healthy'
        else
            echo 'Prometheus health check failed'
            exit 1
        fi
        
        if curl -s -f http://localhost:${grafana_port}/api/health > /dev/null; then
            echo 'Grafana is healthy'
        else
            echo 'Grafana health check failed'
            exit 1
        fi
    "
    
    # Clean up temporary files
    rm -f /tmp/prometheus_${name}.yml /tmp/datasource_${name}.yml /tmp/dashboard_${name}.yml /tmp/docker-compose_${name}.yml
    
    log_success "$name environment deployed successfully!"
    log_info "Access URLs:"
    log_info "  Prometheus: https://${domain}:${prometheus_port}"
    log_info "  Grafana: https://${domain}:${grafana_port} (admin/${grafana_password})"
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS] [ENVIRONMENT]"
    echo ""
    echo "Deploy monitoring environments to VM based on TOML configuration"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -l, --list     List available environments"
    echo "  -c, --config   Show current configuration"
    echo ""
    echo "Environments:"
    echo "  dev            Deploy development environment"
    echo "  prod           Deploy production environment"
    echo "  all            Deploy all enabled environments (default)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Deploy all enabled environments"
    echo "  $0 dev               # Deploy only dev environment"
    echo "  $0 prod              # Deploy only prod environment"
    echo "  $0 --list            # List available environments"
}

# Function to list environments
list_environments() {
    echo "Available environments:"
    echo ""
    
    for env in dev prod; do
        if is_env_enabled "$env"; then
            local name=$(get_env_config "$env" "name")
            local domain=$(get_env_config "$env" "domain")
            echo "  $env: $name ($domain) - ENABLED"
        else
            echo "  $env: DISABLED"
        fi
    done
}

# Function to show configuration
show_config() {
    echo "Current configuration from $CONFIG_FILE:"
    echo ""
    
    # Global config
    echo "Global:"
    echo "  VM Host: $(get_global_config "vm_host")"
    echo "  SSH Key: $(get_global_config "ssh_key")"
    echo "  SSH User: $(get_global_config "ssh_user")"
    echo "  Monitoring Dir: $(get_global_config "monitoring_dir")"
    echo ""
    
    # Environment configs
    for env in dev prod; do
        if is_env_enabled "$env"; then
            echo "$env:"
            echo "  Name: $(get_env_config "$env" "name")"
            echo "  Domain: $(get_env_config "$env" "domain")"
            echo "  Backend URL: $(get_env_config "$env" "backend_url")"
            echo "  Prometheus Port: $(get_env_config "$env" "prometheus_port")"
            echo "  Grafana Port: $(get_env_config "$env" "grafana_port")"
            echo "  Data Retention: $(get_env_config "$env" "data_retention")"
            echo ""
        fi
    done
}

# Main execution
main() {
    local target_env="${1:-all}"
    
    case "$target_env" in
        -h|--help)
            show_help
            exit 0
            ;;
        -l|--list)
            list_environments
            exit 0
            ;;
        -c|--config)
            show_config
            exit 0
            ;;
        all)
            log_info "Deploying all enabled environments..."
            for env in dev prod; do
                if is_env_enabled "$env"; then
                    deploy_environment "$env"
                    echo ""
                fi
            done
            ;;
        dev|prod)
            deploy_environment "$target_env"
            ;;
        *)
            log_error "Unknown environment: $target_env"
            show_help
            exit 1
            ;;
    esac
    
    log_success "Deployment completed!"
}

# Run main function with all arguments
main "$@"
