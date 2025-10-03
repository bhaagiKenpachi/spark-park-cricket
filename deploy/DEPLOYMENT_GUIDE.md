# Spark Park Cricket - Deployment Guide

This guide covers the complete deployment process for the Spark Park Cricket application, including fixes for all identified issues.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed
- SSH access to production servers
- SSH keys configured (see `scripts/setup-ssh.sh`)

### Security Setup
```bash
# Set secure passwords (REQUIRED)
export GRAFANA_ADMIN_PASSWORD_PROD='YourStrongPassword2024!'
export GRAFANA_ADMIN_PASSWORD_DEV='YourDevPassword2024!'

# Deploy with security validation
./deploy/deploy-all.sh setup
```

### Basic Deployment
```bash
# Deploy to development
./deploy/deploy-all.sh deploy-app dev

# Deploy to production
./deploy/deploy-all.sh deploy-app prod
```

## 📋 Deployment Scripts

### Main Deployment Scripts
- `deploy-all.sh` - Master deployment script (handles both app and monitoring)
- `scripts/deploy.sh` - Main deployment script (fixed Docker Compose syntax)
- `scripts/build.sh` - Build Docker images (fixed Docker Compose syntax)
- `scripts/rollback.sh` - Rollback deployments (fixed Docker Compose syntax)
- `scripts/health-check.sh` - Health check services
- `scripts/setup-ssh.sh` - Setup SSH keys and configuration

### Environment-Specific Scripts
- `environments/dev/docker-compose.yml` - Development environment configuration
- `environments/prod/docker-compose.yml` - Production environment configuration

## 🔧 Fixed Issues

### 1. Docker Compose Syntax
**Issue**: Used `docker-compose` instead of `docker compose` (V2 syntax)
**Fix**: Updated all scripts to use `docker compose` consistently

### 2. SSH User Issues
**Issue**: Assumed `root` user but servers use `ubuntu` user
**Fix**: Updated all scripts to use `ubuntu` user with proper `sudo` handling

### 3. Permission Issues
**Issue**: File creation and system operations required elevated privileges
**Fix**: Added proper `sudo` usage in production deployment scripts

## 🌐 Server Configuration

### Production Server (15.235.202.148)
- User: `ubuntu`
- SSH Key: `~/.ssh/spark-cricket-prod`
- Docker: Installed via deployment script
- Services: Backend, Frontend, Redis, Nginx

### Monitoring Server (51.79.143.135)
- User: `ubuntu`
- SSH Key: `~/.ssh/spark-cricket-monitoring`
- Services: Prometheus, Grafana, Loki, Promtail

## 📊 Monitoring Setup

### Deploy Monitoring Stack
```bash
# Deploy monitoring to development
./monitoring/deploy-monitoring.sh dev

# Deploy monitoring to production
./monitoring/deploy-monitoring.sh prod

# Deploy with log forwarding
./monitoring/deploy-monitoring.sh prod --setup-log-forwarding
```

### Log Forwarding
```bash
# Setup log forwarding for production
./monitoring/setup-log-forwarding.sh prod

# Remove log forwarding
./monitoring/setup-log-forwarding.sh prod --remove
```

## 🔍 Troubleshooting

### Common Issues and Solutions

#### 1. Docker Compose Not Found
```bash
# Install Docker Compose plugin
sudo apt-get update
sudo apt-get install docker-compose-plugin
```

#### 2. SSH Permission Denied
```bash
# Check SSH key permissions
chmod 600 ~/.ssh/spark-cricket-prod
chmod 600 ~/.ssh/spark-cricket-monitoring

# Test SSH connection
ssh -i ~/.ssh/spark-cricket-prod ubuntu@15.235.202.148
ssh -i ~/.ssh/spark-cricket-monitoring ubuntu@51.79.143.135
```

#### 3. Loki Configuration Issues
- Fixed: Removed deprecated `enforce_metric_name` field
- Fixed: Added `allow_structured_metadata: false` for schema compatibility

#### 4. Log Forwarding Issues
- Fixed: Replaced complex shell script with Promtail container
- Fixed: Proper Docker socket access and volume mounts

### Health Checks
```bash
# Check application health
./scripts/health-check.sh prod

# Check monitoring services
curl -I http://51.79.143.135:3003  # Grafana
curl -I http://51.79.143.135:9092  # Prometheus
curl -I http://51.79.143.135:3102  # Loki
```

## 📈 Access URLs

### Application
- **Production**: http://15.235.202.148
- **Health Check**: http://15.235.202.148:8080/health

### Monitoring (Production)
- **Grafana**: http://51.79.143.135:3003 (admin/prod123)
- **Prometheus**: http://51.79.143.135:9092
- **Loki**: http://51.79.143.135:3102

### Monitoring (Development)
- **Grafana**: http://51.79.143.135:3001 (admin/dev123)
- **Prometheus**: http://51.79.143.135:9091
- **Loki**: http://51.79.143.135:3101

## 🎯 Best Practices

1. **Always test in development first**
2. **Use version control for all configuration changes**
3. **Monitor logs during deployment**
4. **Run health checks after deployment**
5. **Keep SSH keys secure and properly configured**
6. **Use the standard scripts for consistency**
7. **Set strong passwords for all services**
8. **Never commit .env files to version control**
9. **Use environment variables for sensitive data**
10. **Regularly rotate credentials and SSH keys**

## 📞 Support

For deployment issues:
1. Check the troubleshooting section
2. Review script logs for error messages
3. Verify SSH connectivity and permissions
4. Ensure all prerequisites are met
