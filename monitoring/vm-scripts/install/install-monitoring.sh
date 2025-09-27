#!/bin/bash

# Environment-agnostic Prometheus and Grafana installation script
# This script installs Docker, Docker Compose, and sets up monitoring stack

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   warn "Running as root. This is not recommended but will proceed with caution."
fi

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if [ -f /etc/debian_version ]; then
        OS="debian"
    elif [ -f /etc/redhat-release ]; then
        OS="redhat"
    else
        error "Unsupported Linux distribution"
    fi
else
    error "This script only supports Linux systems"
fi

log "Detected OS: $OS"

# Function to install Docker
install_docker() {
    log "Installing Docker..."
    
    if command -v docker &> /dev/null; then
        log "Docker is already installed"
        return
    fi
    
    if [[ "$OS" == "debian" ]]; then
        # Update package index
        sudo apt-get update
        
        # Install prerequisites
        sudo apt-get install -y \
            ca-certificates \
            curl \
            gnupg \
            lsb-release
        
        # Add Docker's official GPG key
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        
        # Set up repository
        echo \
            "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
            $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        # Install Docker Engine
        sudo apt-get update
        sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
    elif [[ "$OS" == "redhat" ]]; then
        # Install prerequisites
        sudo yum install -y yum-utils
        
        # Add Docker repository
        sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
        
        # Install Docker Engine
        sudo yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    fi
    
    # Start and enable Docker
    sudo systemctl start docker
    sudo systemctl enable docker
    
    # Add current user to docker group
    sudo usermod -aG docker $USER
    
    log "Docker installed successfully"
}

# Function to install Docker Compose
install_docker_compose() {
    log "Installing Docker Compose..."
    
    if command -v docker-compose &> /dev/null; then
        log "Docker Compose is already installed"
        return
    fi
    
    # Install Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    # Create symlink
    sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    
    log "Docker Compose installed successfully"
}

# Function to create monitoring directory structure
create_monitoring_structure() {
    log "Creating monitoring directory structure..."
    
    MONITORING_DIR="/opt/monitoring"
    
    sudo mkdir -p $MONITORING_DIR/{prometheus,grafana,alertmanager}
    sudo mkdir -p $MONITORING_DIR/grafana/{dashboards,datasources}
    sudo mkdir -p $MONITORING_DIR/prometheus/rules
    
    # Set permissions
    sudo chown -R $USER:$USER $MONITORING_DIR
    
    log "Monitoring directory structure created at $MONITORING_DIR"
}

# Function to create systemd service for monitoring
create_systemd_service() {
    log "Creating systemd service for monitoring stack..."
    
    cat << EOF | sudo tee /etc/systemd/system/monitoring-stack.service
[Unit]
Description=Monitoring Stack (Prometheus, Grafana, Alertmanager)
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/monitoring
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
User=$USER
Group=$USER

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable monitoring-stack.service
    
    log "Systemd service created and enabled"
}

# Main installation function
main() {
    log "Starting monitoring stack installation..."
    
    # Update system packages
    log "Updating system packages..."
    if [[ "$OS" == "debian" ]]; then
        sudo apt-get update
        sudo apt-get upgrade -y
    elif [[ "$OS" == "redhat" ]]; then
        sudo yum update -y
    fi
    
    # Install required packages
    install_docker
    install_docker_compose
    create_monitoring_structure
    create_systemd_service
    
    log "Installation completed successfully!"
    log "Please log out and log back in to ensure Docker group permissions are applied."
    log "Then run the deployment script to configure and start the monitoring stack."
}

# Run main function
main "$@"
