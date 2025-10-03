#!/bin/bash

# SSH Setup Script for Production Deployment
# This script helps set up SSH keys for production deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROD_SERVER="15.235.202.148"
PROD_USER="ubuntu"
TEE_AUTH_DIR="$HOME/dojima/tee-auth"
SSH_DIR="$HOME/.ssh"
SSH_KEY_NAME="spark-cricket-prod"
SSH_CONFIG_FILE="$SSH_DIR/config"

# Functions
log() {
    echo -e "${BLUE}[SSH Setup]${NC} $1"
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
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --copy-key KEY_PATH    Copy existing SSH key from tee-auth directory"
    echo "  --generate-new         Generate a new SSH key pair"
    echo "  --test-connection      Test SSH connection to production server"
    echo "  --setup-config         Setup SSH config file"
    echo "  --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --copy-key id_rsa"
    echo "  $0 --generate-new"
    echo "  $0 --test-connection"
}

find_ssh_keys() {
    log "Looking for SSH keys in tee-auth directory..."
    
    if [[ ! -d "$TEE_AUTH_DIR" ]]; then
        error "tee-auth directory not found at $TEE_AUTH_DIR"
        return 1
    fi
    
    echo "Available SSH keys in $TEE_AUTH_DIR:"
    find "$TEE_AUTH_DIR" -name "id_rsa*" -o -name "*.pem" -o -name "*.key" 2>/dev/null | while read -r key_file; do
        if [[ -f "$key_file" ]]; then
            local key_type
            key_type=$(file "$key_file" 2>/dev/null | grep -o "OpenSSH\|RSA\|DSA\|ECDSA\|ED25519" | head -1 || echo "Unknown")
            echo "  $key_file ($key_type)"
        fi
    done
}

copy_ssh_key() {
    local source_key="$1"
    
    if [[ -z "$source_key" ]]; then
        error "Source key path is required"
        find_ssh_keys
        exit 1
    fi
    
    local source_path="$TEE_AUTH_DIR/$source_key"
    
    if [[ ! -f "$source_path" ]]; then
        error "SSH key not found at $source_path"
        find_ssh_keys
        exit 1
    fi
    
    log "Copying SSH key from $source_path..."
    
    # Create SSH directory if it doesn't exist
    mkdir -p "$SSH_DIR"
    chmod 700 "$SSH_DIR"
    
    # Copy the key
    cp "$source_path" "$SSH_DIR/$SSH_KEY_NAME"
    chmod 600 "$SSH_DIR/$SSH_KEY_NAME"
    
    # Copy public key if it exists
    if [[ -f "$source_path.pub" ]]; then
        cp "$source_path.pub" "$SSH_DIR/$SSH_KEY_NAME.pub"
        chmod 644 "$SSH_DIR/$SSH_KEY_NAME.pub"
        success "SSH key and public key copied successfully"
    else
        success "SSH key copied successfully"
        warning "Public key not found at $source_path.pub"
    fi
    
    success "SSH key setup completed: $SSH_DIR/$SSH_KEY_NAME"
}

generate_new_key() {
    log "Generating new SSH key pair..."
    
    # Create SSH directory if it doesn't exist
    mkdir -p "$SSH_DIR"
    chmod 700 "$SSH_DIR"
    
    # Generate new SSH key
    ssh-keygen -t ed25519 -f "$SSH_DIR/$SSH_KEY_NAME" -N "" -C "spark-cricket-prod-$(date +%Y%m%d)"
    
    # Set proper permissions
    chmod 600 "$SSH_DIR/$SSH_KEY_NAME"
    chmod 644 "$SSH_DIR/$SSH_KEY_NAME.pub"
    
    success "New SSH key pair generated: $SSH_DIR/$SSH_KEY_NAME"
    
    echo ""
    log "Public key content:"
    cat "$SSH_DIR/$SSH_KEY_NAME.pub"
    echo ""
    
    warning "You need to add this public key to the production server:"
    echo "  1. Copy the public key above"
    echo "  2. SSH into the production server: ssh $PROD_USER@$PROD_SERVER"
    echo "  3. Add the key to ~/.ssh/authorized_keys"
    echo "  4. Or run: ssh-copy-id -i $SSH_DIR/$SSH_KEY_NAME.pub $PROD_USER@$PROD_SERVER"
}

setup_ssh_config() {
    log "Setting up SSH config file..."
    
    # Create SSH config file if it doesn't exist
    if [[ ! -f "$SSH_CONFIG_FILE" ]]; then
        touch "$SSH_CONFIG_FILE"
        chmod 600 "$SSH_CONFIG_FILE"
    fi
    
    # Check if host already exists
    if grep -q "Host spark-cricket-prod" "$SSH_CONFIG_FILE"; then
        warning "SSH config for spark-cricket-prod already exists"
        return 0
    fi
    
    # Add host configuration
    cat >> "$SSH_CONFIG_FILE" << EOF

Host spark-cricket-prod
    HostName $PROD_SERVER
    User $PROD_USER
    IdentityFile $SSH_DIR/$SSH_KEY_NAME
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
    ServerAliveInterval 60
    ServerAliveCountMax 3
EOF

    success "SSH config updated: $SSH_CONFIG_FILE"
}

test_connection() {
    log "Testing SSH connection to production server..."
    
    if [[ ! -f "$SSH_DIR/$SSH_KEY_NAME" ]]; then
        error "SSH key not found at $SSH_DIR/$SSH_KEY_NAME"
        error "Please run: $0 --copy-key KEY_NAME or $0 --generate-new"
        exit 1
    fi
    
    # Test SSH connection
    if ssh -i "$SSH_DIR/$SSH_KEY_NAME" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$PROD_USER@$PROD_SERVER" "echo 'SSH connection successful'"; then
        success "SSH connection to production server is working"
        
        # Test Docker access
        if ssh -i "$SSH_DIR/$SSH_KEY_NAME" -o StrictHostKeyChecking=no "$PROD_USER@$PROD_SERVER" "docker --version" >/dev/null 2>&1; then
            success "Docker is available on production server"
        else
            warning "Docker is not available on production server"
        fi
        
        # Test directory access
        if ssh -i "$SSH_DIR/$SSH_KEY_NAME" -o StrictHostKeyChecking=no "$PROD_USER@$PROD_SERVER" "ls -la /opt" >/dev/null 2>&1; then
            success "Directory access is working"
        else
            warning "Cannot access /opt directory (may need sudo)"
        fi
        
    else
        error "SSH connection to production server failed"
        echo ""
        log "Troubleshooting steps:"
        echo "  1. Verify the server IP: $PROD_SERVER"
        echo "  2. Check if the server is running"
        echo "  3. Verify SSH key is correct"
        echo "  4. Check firewall settings"
        echo "  5. Try: ssh -v -i $SSH_DIR/$SSH_KEY_NAME $PROD_USER@$PROD_SERVER"
        exit 1
    fi
}

show_ssh_info() {
    log "SSH Configuration Summary:"
    echo ""
    echo "  Server: $PROD_SERVER"
    echo "  User: $PROD_USER"
    echo "  SSH Key: $SSH_DIR/$SSH_KEY_NAME"
    echo "  SSH Config: $SSH_CONFIG_FILE"
    echo ""
    
    if [[ -f "$SSH_DIR/$SSH_KEY_NAME" ]]; then
        echo "  Key exists: ✓"
        echo "  Key permissions: $(ls -l "$SSH_DIR/$SSH_KEY_NAME" | awk '{print $1}')"
    else
        echo "  Key exists: ✗"
    fi
    
    if [[ -f "$SSH_CONFIG_FILE" ]] && grep -q "Host spark-cricket-prod" "$SSH_CONFIG_FILE"; then
        echo "  SSH config: ✓"
    else
        echo "  SSH config: ✗"
    fi
    
    echo ""
    log "Quick commands:"
    echo "  Test connection: $0 --test-connection"
    echo "  SSH to server: ssh spark-cricket-prod"
    echo "  SSH with key: ssh -i $SSH_DIR/$SSH_KEY_NAME $PROD_USER@$PROD_SERVER"
}

main() {
    local copy_key=""
    local generate_new=false
    local test_connection=false
    local setup_config=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --copy-key)
                copy_key="$2"
                shift 2
                ;;
            --generate-new)
                generate_new=true
                shift
                ;;
            --test-connection)
                test_connection=true
                shift
                ;;
            --setup-config)
                setup_config=true
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
    
    # Show current status
    show_ssh_info
    
    # Execute requested actions
    if [[ -n "$copy_key" ]]; then
        copy_ssh_key "$copy_key"
        setup_ssh_config
    fi
    
    if [[ "$generate_new" == "true" ]]; then
        generate_new_key
        setup_ssh_config
    fi
    
    if [[ "$setup_config" == "true" ]]; then
        setup_ssh_config
    fi
    
    if [[ "$test_connection" == "true" ]]; then
        test_connection
    fi
    
    # If no specific action was requested, show help
    if [[ -z "$copy_key" && "$generate_new" == "false" && "$test_connection" == "false" && "$setup_config" == "false" ]]; then
        echo ""
        usage
        echo ""
        log "Available SSH keys in tee-auth directory:"
        find_ssh_keys
    fi
}

# Run main function with all arguments
main "$@"
