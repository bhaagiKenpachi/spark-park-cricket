# Monitoring Stack Deployment

This directory contains environment-agnostic scripts for deploying Prometheus, Grafana, and Alertmanager to a remote VM.

## Overview

The monitoring stack includes:
- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboards
- **Alertmanager**: Alert routing and notification
- **Node Exporter**: System metrics
- **cAdvisor**: Container metrics

## Quick Start

### 1. Deploy to Development Environment

```bash
# Quick deployment to your VM
./deploy/quick-deploy.sh deploy-dev
```

### 2. Check Service Status

```bash
./deploy/quick-deploy.sh status
```

### 3. Access Services

- **Prometheus**: http://51.79.143.135:9090
- **Grafana**: http://51.79.143.135:3000 (admin/admin123)
- **Alertmanager**: http://51.79.143.135:9093

## Directory Structure

```
vm-scripts/
├── install/                 # Installation scripts
│   └── install-monitoring.sh
├── deploy/                  # Deployment scripts
│   ├── deploy-monitoring.sh
│   └── quick-deploy.sh
├── config/                  # Configuration files
│   ├── docker-compose.yml
│   ├── prometheus/
│   │   ├── prometheus.yml
│   │   └── rules/
│   │       └── cricket-alerts.yml
│   ├── grafana/
│   │   ├── datasources/
│   │   └── dashboards/
│   └── alertmanager/
│       └── alertmanager.yml
├── templates/               # Configuration templates
│   └── env.template
└── README.md
```

## Configuration

### Environment Variables

The monitoring stack is configured using environment variables. Copy `templates/env.template` to `config/.env` and customize:

```bash
cp templates/env.template config/.env
```

Key configuration options:
- `ENVIRONMENT`: dev/prod
- `ENVIRONMENT_DOMAIN`: Domain for services
- `GRAFANA_ADMIN_PASSWORD`: Grafana admin password
- `BACKEND_HOST`: Backend service host
- `BACKEND_PORT`: Backend service port

### Environment-Specific Configuration

#### Development Environment
```bash
ENVIRONMENT=dev
ENVIRONMENT_DOMAIN=localhost
GRAFANA_ADMIN_PASSWORD=admin123
```

#### Production Environment
```bash
ENVIRONMENT=prod
ENVIRONMENT_DOMAIN=cricket.example.com
GRAFANA_ADMIN_PASSWORD=secure_password_here
CRITICAL_ALERT_EMAIL=admin@example.com
SMTP_HOST=smtp.example.com:587
```

## Deployment Scripts

### 1. Full Deployment Script

```bash
./deploy/deploy-monitoring.sh [OPTIONS]
```

Options:
- `-h, --host HOST`: Remote host IP (default: 51.79.143.135)
- `-u, --user USER`: Remote username (default: ubuntu)
- `-k, --key KEY`: SSH private key path
- `-e, --env ENV`: Environment (dev/prod)

Examples:
```bash
# Deploy to default settings
./deploy/deploy-monitoring.sh

# Deploy to custom host
./deploy/deploy-monitoring.sh -h 192.168.1.100 -u admin

# Deploy to production
./deploy/deploy-monitoring.sh -e prod -k /path/to/key.pem
```

### 2. Quick Deployment Script

```bash
./deploy/quick-deploy.sh [COMMAND] [OPTIONS]
```

Commands:
- `deploy-dev`: Deploy to development
- `deploy-prod`: Deploy to production
- `status`: Check service status
- `logs [service]`: View logs
- `restart`: Restart services
- `stop`: Stop services
- `update-config`: Update configuration

Examples:
```bash
# Deploy to development
./deploy/quick-deploy.sh deploy-dev

# Check status
./deploy/quick-deploy.sh status

# View Prometheus logs
./deploy/quick-deploy.sh logs prometheus

# Restart services
./deploy/quick-deploy.sh restart
```

## Service Management

### SSH Access
```bash
ssh -i /Users/luffybhaagi/dojima/tee-auth/runner_private_key.pem ubuntu@51.79.143.135
```

### Service Control
```bash
# Navigate to monitoring directory
cd /opt/monitoring

# Start services
docker-compose up -d

# Stop services
docker-compose down

# Restart services
docker-compose restart

# View logs
docker-compose logs -f [service_name]

# Check status
docker-compose ps
```

### Systemd Service
The monitoring stack is configured as a systemd service:
```bash
# Enable service
sudo systemctl enable monitoring-stack.service

# Start service
sudo systemctl start monitoring-stack.service

# Check status
sudo systemctl status monitoring-stack.service
```

## Monitoring Features

### Prometheus Metrics
- System metrics (CPU, memory, disk)
- Application metrics (HTTP requests, response times)
- Database metrics (connections, queries)
- Custom cricket application metrics

### Grafana Dashboards
- System performance dashboard
- Application metrics dashboard
- Database performance dashboard
- Cricket-specific metrics dashboard

### Alerting Rules
- High CPU/Memory usage
- Service downtime
- High error rates
- Database connection issues
- Application-specific alerts

## Troubleshooting

### Common Issues

1. **SSH Connection Failed**
   ```bash
   # Check SSH key permissions
   chmod 600 /Users/luffybhaagi/dojima/tee-auth/runner_private_key.pem
   
   # Test connection
   ssh -i /Users/luffybhaagi/dojima/tee-auth/runner_private_key.pem ubuntu@51.79.143.135
   ```

2. **Services Not Starting**
   ```bash
   # Check Docker status
   docker --version
   docker-compose --version
   
   # Check service logs
   ./deploy/quick-deploy.sh logs
   ```

3. **Port Access Issues**
   ```bash
   # Check if ports are open
   netstat -tlnp | grep -E ':(3000|9090|9093)'
   
   # Check firewall
   sudo ufw status
   ```

### Log Locations
- Docker logs: `docker-compose logs [service]`
- System logs: `/var/log/syslog`
- Service logs: `journalctl -u monitoring-stack.service`

## Security Considerations

### Production Deployment
- Change default passwords
- Configure SSL/TLS certificates
- Set up proper firewall rules
- Use environment-specific secrets
- Enable authentication for all services

### Access Control
- Restrict SSH access to specific IPs
- Use strong SSH keys
- Regularly rotate credentials
- Monitor access logs

## Maintenance

### Updates
```bash
# Update configuration
./deploy/quick-deploy.sh update-config

# Restart services
./deploy/quick-deploy.sh restart
```

### Backup
```bash
# Backup configuration
tar -czf monitoring-backup-$(date +%Y%m%d).tar.gz /opt/monitoring

# Backup data volumes
docker run --rm -v prometheus_data:/data -v $(pwd):/backup alpine tar czf /backup/prometheus-data-$(date +%Y%m%d).tar.gz -C /data .
```

## Support

For issues or questions:
1. Check service logs: `./deploy/quick-deploy.sh logs`
2. Verify service status: `./deploy/quick-deploy.sh status`
3. Review configuration files in `config/` directory
4. Check SSH connectivity and permissions
