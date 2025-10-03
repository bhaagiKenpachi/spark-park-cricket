# Spark Park Cricket Deployment

This directory contains deployment scripts and configurations for the Spark Park Cricket application.

## Structure

```
deploy/
├── environments/
│   ├── dev/
│   │   ├── docker-compose.yml
│   │   ├── .env
│   │   └── deploy.sh
│   └── prod/
│       ├── docker-compose.yml
│       ├── .env
│       └── deploy.sh
├── scripts/
│   ├── deploy.sh
│   ├── build.sh
│   ├── health-check.sh
│   └── rollback.sh
└── configs/
    ├── nginx.conf
    └── docker-compose.base.yml
```

## Environments

### Development (dev)
- Local development environment
- Uses localhost URLs
- Development database configurations
- Hot reload enabled

### Production (prod)
- Production server deployment
- Uses production URLs and configurations
- Production database settings
- Optimized for performance

## Quick Start

1. **Setup SSH key** (one-time):
   ```bash
   # Copy your SSH key from tee-auth folder
   cp ~/dojima/tee-auth/your-ssh-key ~/.ssh/spark-cricket-prod
   chmod 600 ~/.ssh/spark-cricket-prod
   ```

2. **Deploy to production**:
   ```bash
   ./scripts/deploy.sh prod
   ```

3. **Deploy to development**:
   ```bash
   ./scripts/deploy.sh dev
   ```

## Configuration

Each environment has its own:
- `docker-compose.yml` - Environment-specific Docker configuration
- `.env` - Environment variables
- `deploy.sh` - Environment-specific deployment script

## Scripts

- `deploy.sh` - Main deployment script (environment-agnostic)
- `build.sh` - Build Docker images
- `health-check.sh` - Check application health
- `rollback.sh` - Rollback to previous version

## Production Server

- **IP**: 15.235.202.148
- **SSH Key**: Located in `~/dojima/tee-auth/` folder
- **User**: ubuntu (default for most cloud providers)

## Services

The deployment includes:
- **Backend API**: Go-based REST API with GraphQL support
- **Frontend**: Next.js React application
- **Redis**: Caching and session storage
- **Nginx**: Reverse proxy (production only)
