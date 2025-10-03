# Spark Park Cricket Monitoring Stack

This directory contains the monitoring infrastructure for the Spark Park Cricket application, including Prometheus, Grafana, Loki, and Promtail for metrics collection, visualization, and log aggregation.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Server Instance │    │ Monitoring Instance │    │   Application   │
│  15.235.202.148 │    │  51.79.143.135   │    │     Logs        │
│                 │    │                  │    │                 │
│ ┌─────────────┐ │    │ ┌──────────────┐ │    │ ┌─────────────┐ │
│ │   Backend   │ │───▶│ │   Promtail   │ │    │ │   Backend   │ │
│ │  Container  │ │    │ │              │ │    │ │  Container  │ │
│ └─────────────┘ │    │ └──────────────┘ │    │ └─────────────┘ │
│                 │    │        │         │    │                 │
│ ┌─────────────┐ │    │ ┌──────────────┐ │    │                 │
│ │ Log Files   │ │───▶│ │     Loki     │ │◀───│                 │
│ │ /var/log/   │ │    │ │              │ │    │                 │
│ └─────────────┘ │    │ └──────────────┘ │    │                 │
└─────────────────┘    │        │         │    └─────────────────┘
                       │ ┌──────────────┐ │
                       │ │   Grafana    │ │
                       │ │              │ │
                       │ └──────────────┘ │
                       │        │         │
                       │ ┌──────────────┐ │
                       │ │  Prometheus  │ │
                       │ │              │ │
                       │ └──────────────┘ │
                       └──────────────────┘
```

## Components

### Prometheus
- **Purpose**: Metrics collection and storage
- **Ports**: 
  - Dev: 9091
  - Prod: 9092
- **Configuration**: Scrapes metrics from backend containers and system metrics

### Grafana
- **Purpose**: Metrics and logs visualization
- **Ports**:
  - Dev: 3001
  - Prod: 3003
- **Credentials**:
  - Dev: admin / dev123
  - Prod: admin / prod123

### Loki
- **Purpose**: Log aggregation and storage
- **Ports**:
  - Dev: 3101
  - Prod: 3102
- **Configuration**: Stores logs from Promtail and forwarded logs from server instance

### Promtail
- **Purpose**: Log collection and forwarding
- **Configuration**: Collects logs from Docker containers and forwards to Loki

## Environments

### Development Environment
- **Monitoring Instance**: 51.79.143.135
- **Services**: Prometheus (9091), Grafana (3001), Loki (3101)
- **Configuration**: `monitoring/environments/dev/`

### Production Environment
- **Monitoring Instance**: 51.79.143.135
- **Services**: Prometheus (9092), Grafana (3003), Loki (3102)
- **Configuration**: `monitoring/environments/prod/`

## Deployment

### Prerequisites
1. SSH access to monitoring instance (51.79.143.135)
2. SSH access to server instance (15.235.202.148)
3. SSH keys configured:
   - `~/.ssh/spark-cricket-monitoring` (for monitoring instance)
   - `~/.ssh/spark-cricket-prod` (for server instance)

### Deploy Monitoring Stack

```bash
# Deploy development environment
./monitoring/deploy-monitoring.sh dev

# Deploy production environment
./monitoring/deploy-monitoring.sh prod
```

### Setup Log Forwarding

```bash
# Setup log forwarding from server instance to monitoring instance
./monitoring/setup-log-forwarding.sh
```

## Access URLs

### Development Environment
- **Grafana**: http://51.79.143.135:3001
- **Prometheus**: http://51.79.143.135:9091
- **Loki**: http://51.79.143.135:3101

### Production Environment
- **Grafana**: http://51.79.143.135:3003
- **Prometheus**: http://51.79.143.135:9092
- **Loki**: http://51.79.143.135:3102

## Dashboards

### Available Dashboards
1. **Database Performance Dashboard**
   - URL: `/d/cricket-db-performance`
   - Shows database operation metrics, response times, and API performance

2. **Backend Logs Dashboard**
   - Dev: `/d/cricket-backend-logs-dev`
   - Prod: `/d/cricket-backend-logs-prod`
   - Shows real-time logs from backend containers

## Log Collection

### Log Sources
1. **Docker Container Logs**: Collected by Promtail from Docker containers
2. **Server Instance Logs**: Forwarded from server instance (15.235.202.148) via log forwarding service

### Log Forwarding Service
- **Service Name**: `cricket-log-forwarder`
- **Location**: `/opt/forward-logs.sh`
- **Purpose**: Forwards backend container logs from server instance to Loki on monitoring instance
- **Frequency**: Every 30 seconds

### Service Management
```bash
# Check service status
systemctl status cricket-log-forwarder

# View service logs
journalctl -u cricket-log-forwarder -f

# Restart service
systemctl restart cricket-log-forwarder
```

## Configuration Files

### Loki Configuration
- **Dev**: `monitoring/environments/dev/loki/loki.yml`
- **Prod**: `monitoring/environments/prod/loki/loki.yml`

### Promtail Configuration
- **Dev**: `monitoring/environments/dev/promtail/promtail.yml`
- **Prod**: `monitoring/environments/prod/promtail/promtail.yml`

### Prometheus Configuration
- **Dev**: `monitoring/environments/dev/prometheus/prometheus.yml`
- **Prod**: `monitoring/environments/prod/prometheus/prometheus.yml`

### Grafana Configuration
- **Datasources**: `monitoring/environments/*/grafana/datasources/prometheus.yml`
- **Dashboards**: `monitoring/environments/*/grafana/dashboards/`

## Monitoring and Troubleshooting

### Check Service Status
```bash
# On monitoring instance
cd /opt/monitoring/environments/{dev|prod}
docker compose ps

# Check individual container logs
docker compose logs prometheus
docker compose logs grafana
docker compose logs loki
docker compose logs promtail
```

### Common Issues

1. **Loki not receiving logs**
   - Check Promtail configuration
   - Verify log forwarding service is running on server instance
   - Check network connectivity between instances

2. **Grafana not showing data**
   - Verify datasource configuration
   - Check Prometheus targets are up
   - Verify Loki datasource is configured correctly

3. **High memory usage**
   - Check Loki retention settings
   - Monitor Prometheus storage usage
   - Adjust log collection frequency if needed

### Log Queries

#### Grafana Log Queries (Loki)
```logql
# All backend logs
{container_name=~"cricket-backend.*"}

# Logs from specific instance
{instance="15.235.202.148"}

# Error logs
{container_name=~"cricket-backend.*"} |= "ERROR"

# Logs with specific text
{container_name=~"cricket-backend.*"} |= "database"
```

## Security Considerations

1. **Network Access**: Monitoring services are exposed on public IPs with basic authentication
2. **Credentials**: Default passwords should be changed in production
3. **Firewall**: Consider restricting access to monitoring ports
4. **Log Data**: Logs may contain sensitive information - ensure proper access controls

## Maintenance

### Backup
- Grafana dashboards and datasources are stored in configuration files
- Loki data is stored in Docker volumes
- Regular backups of configuration and data volumes recommended

### Updates
- Update Docker images regularly
- Monitor for security updates
- Test updates in development environment first

### Scaling
- For high-volume environments, consider:
  - Loki clustering
  - Prometheus federation
  - Horizontal scaling of Promtail instances

## Support

For issues or questions regarding the monitoring setup:
1. Check service logs and status
2. Verify network connectivity
3. Review configuration files
4. Check Grafana and Prometheus targets status
