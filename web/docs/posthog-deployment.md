# PostHog Self-Hosted Deployment Guide

This guide provides step-by-step instructions for deploying self-hosted PostHog instances for both development and production environments.

## Table of Contents

1. [Overview](#overview)
2. [VM Requirements](#vm-requirements)
3. [Production Setup](#production-setup)
4. [Development Setup](#development-setup)
5. [Configuration](#configuration)
6. [Monitoring & Maintenance](#monitoring--maintenance)
7. [Backup & Recovery](#backup--recovery)
8. [Troubleshooting](#troubleshooting)

## Overview

For a system with 50-100 users per week, we recommend a lightweight self-hosted setup optimized for cost and simplicity.

### Architecture

```
┌─────────────────┐
│  Frontend App   │ (localhost:3000 or production URL)
│   (Next.js)     │
└────────┬────────┘
         │ HTTPS
         ▼
┌─────────────────┐
│  PostHog VM     │ (Your VM IP or domain)
│  ┌───────────┐  │
│  │ Nginx     │  │ Port 80/443
│  │ (Reverse  │  │
│  │  Proxy)   │  │
│  └─────┬─────┘  │
│        │        │
│  ┌─────▼─────┐  │
│  │ PostHog   │  │ Port 8000
│  │  Web      │  │
│  └───────────┘  │
│  ┌───────────┐  │
│  │ PostHog   │  │
│  │  Worker   │  │
│  └───────────┘  │
│  ┌───────────┐  │
│  │PostgreSQL │  │ Port 5432
│  └───────────┘  │
│  ┌───────────┐  │
│  │ClickHouse │  │ Port 8123/9000
│  └───────────┘  │
│  ┌───────────┐  │
│  │  Redis    │  │ Port 6379
│  └───────────┘  │
└─────────────────┘
```

## VM Requirements

### Minimum Specifications (Suitable for 50-100 users/week)

- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 20GB SSD
- **OS**: Ubuntu 22.04 LTS
- **Network**: 100 Mbps

### Recommended Providers

| Provider | Plan | Monthly Cost | Specs |
|----------|------|-------------|-------|
| Hetzner | CX21 | €4.90 | 2 vCPU, 4GB RAM, 40GB SSD |
| DigitalOcean | Basic | $24 | 2 vCPU, 4GB RAM, 80GB SSD |
| Linode | Nanode 4GB | $12 | 2 vCPU, 4GB RAM, 80GB SSD |
| AWS EC2 | t3.medium | ~$30 | 2 vCPU, 4GB RAM |

**Recommendation**: Hetzner for best value

## Production Setup

### Step 1: Provision VM

```bash
# SSH into your production VM
ssh root@your-production-ip

# Update system
apt update && apt upgrade -y

# Install required packages
apt install -y curl git ufw
```

### Step 2: Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
apt install docker-compose-plugin

# Verify installation
docker --version
docker compose version
```

### Step 3: Create PostHog Directory

```bash
# Create directory for PostHog
mkdir -p /opt/posthog
cd /opt/posthog
```

### Step 4: Create Docker Compose File

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    container_name: posthog-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: posthog
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: posthog
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U posthog"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:7-alpine
    container_name: posthog-redis
    restart: unless-stopped
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  clickhouse:
    image: clickhouse/clickhouse-server:23.11-alpine
    container_name: posthog-clickhouse
    restart: unless-stopped
    environment:
      CLICKHOUSE_DB: posthog
    volumes:
      - clickhouse-data:/var/lib/clickhouse
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "localhost:8123/ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  posthog:
    image: posthog/posthog:latest
    container_name: posthog-web
    restart: unless-stopped
    ports:
      - "8000:8000"
    environment:
      DATABASE_URL: postgres://posthog:${POSTGRES_PASSWORD}@postgres:5432/posthog
      REDIS_URL: redis://redis:6379/
      CLICKHOUSE_HOST: clickhouse
      SECRET_KEY: ${SECRET_KEY}
      SITE_URL: ${SITE_URL}
      IS_BEHIND_PROXY: true
      DISABLE_SECURE_SSL_REDIRECT: true
      SECURE_COOKIES: false
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      clickhouse:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8000/_health"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 60s

  worker:
    image: posthog/posthog:latest
    container_name: posthog-worker
    restart: unless-stopped
    command: ./bin/docker-worker-celery --with-scheduler
    environment:
      DATABASE_URL: postgres://posthog:${POSTGRES_PASSWORD}@postgres:5432/posthog
      REDIS_URL: redis://redis:6379/
      CLICKHOUSE_HOST: clickhouse
      SECRET_KEY: ${SECRET_KEY}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      clickhouse:
        condition: service_healthy

volumes:
  postgres-data:
  clickhouse-data:
```

### Step 5: Create Environment File

Create `.env`:

```bash
# Generate secure password
POSTGRES_PASSWORD=$(openssl rand -base64 32)

# Generate secret key
SECRET_KEY=$(openssl rand -base64 64)

# Set your production URL
SITE_URL=https://posthog.yourdomain.com

# Save to .env file
cat > .env <<EOF
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
SECRET_KEY=${SECRET_KEY}
SITE_URL=https://posthog.yourdomain.com
EOF

# Secure the file
chmod 600 .env
```

### Step 6: Configure Firewall

```bash
# Allow SSH, HTTP, and HTTPS
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

### Step 7: Install and Configure Nginx

```bash
# Install Nginx
apt install -y nginx certbot python3-certbot-nginx

# Create Nginx configuration
cat > /etc/nginx/sites-available/posthog <<'EOF'
server {
    listen 80;
    server_name posthog.yourdomain.com;

    location / {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
    }

    client_max_body_size 100M;
}
EOF

# Enable site
ln -s /etc/nginx/sites-available/posthog /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

### Step 8: Install SSL Certificate

```bash
# Get Let's Encrypt certificate
certbot --nginx -d posthog.yourdomain.com

# Test auto-renewal
certbot renew --dry-run
```

### Step 9: Start PostHog

```bash
cd /opt/posthog
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f posthog
```

### Step 10: Initial Setup

1. Open `https://posthog.yourdomain.com` in your browser
2. Create admin account
3. Set up organization
4. Create project
5. Copy Project API Key

### Step 11: Update Application Configuration

Update your application's environment variables:

```env
# Production
NEXT_PUBLIC_POSTHOG_KEY=phc_your_production_api_key
NEXT_PUBLIC_POSTHOG_HOST=https://posthog.yourdomain.com
NEXT_PUBLIC_POSTHOG_DEBUG=false
```

## Development Setup

Follow the same steps as production but with different configurations:

### Dev Environment File

```env
POSTGRES_PASSWORD=dev_password_change_in_prod
SECRET_KEY=dev_secret_key_change_in_prod
SITE_URL=https://posthog-dev.yourdomain.com
```

### Dev Application Config

```env
# Development
NEXT_PUBLIC_POSTHOG_KEY=phc_your_dev_api_key
NEXT_PUBLIC_POSTHOG_HOST=https://posthog-dev.yourdomain.com
NEXT_PUBLIC_POSTHOG_DEBUG=true
```

### Recommendations for Dev

- Use separate domain: `posthog-dev.yourdomain.com`
- Create separate PostHog project for dev data
- Enable debug mode
- Consider smaller VM (2GB RAM may suffice)

## Configuration

### Resource Limits

For 50-100 users/week, adjust resource limits:

```yaml
# In docker-compose.yml, add deploy limits
posthog:
  deploy:
    resources:
      limits:
        memory: 1.5G
      reservations:
        memory: 512M
```

### Data Retention

Configure in PostHog dashboard:
1. Go to Settings > Project Settings > Data Management
2. Set retention period: 90 days (recommended for your scale)
3. Enable auto-delete of old events

### Performance Tuning

```yaml
# Optimize ClickHouse for small scale
clickhouse:
  environment:
    CLICKHOUSE_MAX_MEMORY_USAGE: 2000000000  # 2GB
```

## Monitoring & Maintenance

### Health Checks

```bash
# Check container health
docker compose ps

# Check disk usage
df -h

# Check memory usage
free -h

# Check logs
docker compose logs --tail=100 posthog
```

### Automated Monitoring Script

Create `/opt/posthog/monitor.sh`:

```bash
#!/bin/bash

# Check if PostHog is responding
if ! curl -f http://localhost:8000/_health > /dev/null 2>&1; then
    echo "PostHog health check failed!"
    docker compose restart posthog
fi

# Check disk space
DISK_USAGE=$(df -h /opt/posthog | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $DISK_USAGE -gt 80 ]; then
    echo "Disk usage is ${DISK_USAGE}%!"
fi
```

Add to crontab:

```bash
# Check every 5 minutes
*/5 * * * * /opt/posthog/monitor.sh
```

### Updates

```bash
cd /opt/posthog

# Pull latest images
docker compose pull

# Restart services
docker compose up -d

# Clean up old images
docker image prune -a -f
```

## Backup & Recovery

### Automated Backup Script

Create `/opt/posthog/backup.sh`:

```bash
#!/bin/bash

BACKUP_DIR="/opt/posthog/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup PostgreSQL
docker compose exec -T postgres pg_dump -U posthog posthog | gzip > "$BACKUP_DIR/postgres_$DATE.sql.gz"

# Backup environment file
cp .env "$BACKUP_DIR/env_$DATE"

# Keep only last 7 days of backups
find $BACKUP_DIR -name "postgres_*.sql.gz" -mtime +7 -delete

echo "Backup completed: $DATE"
```

Add to crontab:

```bash
# Daily backup at 2 AM
0 2 * * * /opt/posthog/backup.sh
```

### Restore from Backup

```bash
# Stop PostHog
docker compose down

# Restore database
gunzip < postgres_backup.sql.gz | docker compose exec -T postgres psql -U posthog posthog

# Start PostHog
docker compose up -d
```

## Troubleshooting

### PostHog Not Starting

```bash
# Check logs
docker compose logs posthog

# Common issues:
# 1. Database not ready - wait 60 seconds
# 2. Port conflict - check if 8000 is used
sudo lsof -i :8000

# 3. Memory issues
docker compose restart posthog
```

### High Memory Usage

```bash
# Check memory
docker stats

# Restart services
docker compose restart

# If persistent, increase VM memory
```

### Disk Space Full

```bash
# Check usage
docker system df

# Clean up
docker system prune -a -f

# Clean old ClickHouse data (if needed)
# In PostHog dashboard: Settings > Data Management
```

### Can't Connect from Application

1. Check firewall rules
2. Verify SSL certificate
3. Test PostHog endpoint:
   ```bash
   curl https://posthog.yourdomain.com/_health
   ```
4. Check CSP headers in application

## Cost Optimization

For 50-100 users/week:

1. **Use Hetzner**: €4.90/month vs $24+ elsewhere
2. **Single VM**: Don't separate services across VMs
3. **90-day retention**: Reduces storage needs
4. **Disable unused features**: In PostHog settings
5. **Schedule downtime**: If only active during business hours

## Next Steps

After deployment:

1. ✅ Test event tracking from application
2. ✅ Set up feature flags
3. ✅ Configure alerts in PostHog
4. ✅ Set up backups
5. ✅ Monitor resource usage
6. ✅ Document your setup
7. ✅ Train team on PostHog dashboard

## Support & Resources

- PostHog Docs: https://posthog.com/docs
- PostHog Community: https://posthog.com/community
- GitHub Issues: https://github.com/PostHog/posthog/issues


