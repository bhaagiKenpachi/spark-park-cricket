# Production Environment Secrets

This document lists all required GitHub Secrets for the **production** deployment workflow.

## 📍 Where to Add These Secrets

1. Go to your GitHub repository
2. Click on **Settings** → **Environments**
3. Click on **prod** environment (create it if it doesn't exist)
4. Add secrets under **Environment secrets**

## ✅ Required Production Secrets

### 🖥️ Production VM Access

| Secret Name | Description | Example |
|------------|-------------|---------|
| `PROD_VM_HOST` | Production VM IP address or hostname | `15.235.202.148` |
| `PROD_VM_USER` | SSH username for production VM | `ubuntu` |
| `PROD_VM_SSH_KEY` | Private SSH key for production VM | Contents of your prod SSH key file |

### 🗄️ Database (Supabase - Production)

| Secret Name | Description |
|------------|-------------|
| `SUPABASE_URL` | Production Supabase project URL |
| `SUPABASE_API_KEY` | Production Supabase anon/public API key |
| `SUPABASE_PUBLISHABLE_KEY` | Production Supabase publishable key |
| `SUPABASE_SECRET_KEY` | Production service role secret key |

**Important**: Use **different Supabase projects** for dev and production!

### 🔐 OAuth (Google - Production)

| Secret Name | Description |
|------------|-------------|
| `GOOGLE_CLIENT_ID` | Production Google OAuth Client ID |
| `GOOGLE_CLIENT_SECRET` | Production Google OAuth Client Secret |
| `GOOGLE_REDIRECT_URL` | Production OAuth callback URL |

**Example**: `https://spark-park.dojima.foundation/api/v1/auth/google/callback`

**Important**: Configure production OAuth credentials in Google Cloud Console with production redirect URI.

### 🔑 Security & Session

| Secret Name | Description |
|------------|-------------|
| `SESSION_SECRET` | Production session encryption key (use different from dev!) |
| `GRAFANA_PASSWORD` | Production Grafana admin password (use strong password!) |

**Generate session secret**:
```bash
openssl rand -base64 64
```

### 🌐 Frontend & CORS Configuration

| Secret Name | Description | Example |
|------------|-------------|---------|
| `FRONTEND_URL` | Production frontend base URL | `https://spark-park.dojima.foundation` |
| `ALLOWED_ORIGINS` | Production CORS allowed origins | `https://spark-park.dojima.foundation` |

**Note**: Use HTTPS for production!

## 📋 Production Secrets Checklist

Before running production deployment:

### VM Access (3)
- [ ] `PROD_VM_HOST`
- [ ] `PROD_VM_USER`
- [ ] `PROD_VM_SSH_KEY`

### Database (4)
- [ ] `SUPABASE_URL` (production project)
- [ ] `SUPABASE_API_KEY`
- [ ] `SUPABASE_PUBLISHABLE_KEY`
- [ ] `SUPABASE_SECRET_KEY`

### OAuth (3)
- [ ] `GOOGLE_CLIENT_ID` (production credentials)
- [ ] `GOOGLE_CLIENT_SECRET`
- [ ] `GOOGLE_REDIRECT_URL`

### Security (2)
- [ ] `SESSION_SECRET` (unique for prod!)
- [ ] `GRAFANA_PASSWORD` (strong password!)

### Frontend (2)
- [ ] `FRONTEND_URL`
- [ ] `ALLOWED_ORIGINS`

**Total: 14 secrets**

## 🔒 Production Security Best Practices

1. ✅ **Use different credentials** for production and dev
2. ✅ **Use HTTPS** for all production URLs
3. ✅ **Strong passwords** - minimum 16 characters
4. ✅ **Rotate secrets** every 90 days
5. ✅ **Separate Supabase projects** - never use dev database in production
6. ✅ **Limit access** - only give production secrets access to necessary team members
7. ✅ **Enable SSL/TLS** - configure Cloudflare or Let's Encrypt
8. ✅ **Backup strategy** - ensure database backups are configured
9. ✅ **Monitoring** - set up alerts in Grafana for production
10. ✅ **Review logs** - monitor deployment logs for any issues

## 🚀 Deploying to Production

### First Time Setup

1. **Add all 14 secrets** to prod environment
2. **Verify DNS** points to production VM
3. **Configure SSL** (recommended)
4. **Test in dev first** before deploying to production

### Running Production Deployment

1. Go to **Actions** tab
2. Select **"Deploy to Production VM"** workflow
3. Click **"Run workflow"**
4. **⚠️ IMPORTANT**: Double-check you're deploying the correct branch!
5. Click **"Run workflow"** to confirm

### After Deployment

Verify all services:
```bash
# Check health endpoints
curl https://your-prod-domain/health
curl https://your-prod-domain/api/v1/health

# Monitor in Grafana
https://your-prod-domain:3002
```

## ⚠️ Production Deployment Warnings

- ⚠️ **No automatic deployments** - Production workflow is **manual trigger only**
- ⚠️ **Review changes** before deploying - check code changes and tests
- ⚠️ **Backup first** - ensure data is backed up before major updates
- ⚠️ **Monitor closely** - watch logs and metrics after deployment
- ⚠️ **Rollback plan** - keep previous image tags for quick rollback
- ⚠️ **Off-hours deployment** - deploy during low-traffic periods

## 🔄 Rollback Process

If production deployment fails or has issues:

```bash
# SSH to production VM
ssh -i ~/.ssh/prod-key user@prod-vm

# Check available images
docker images | grep cricket

# Rollback to previous version
cd ~/cricket-deployment/deploy/environments/prod

# Edit .env to use previous image tag
# BACKEND_IMAGE=ghcr.io/.../cricket-backend:prod-abc1234
# WEB_IMAGE=ghcr.io/.../cricket-web:prod-abc1234

# Restart services
docker compose up -d
```

## 📞 Support

For production deployment issues:
- Check workflow logs in GitHub Actions
- Review container logs: `docker logs container-name`
- Monitor Grafana dashboards
- Check `.github/RUNNER_SETUP.md` for runner issues

