# GitHub Secrets Configuration

This document lists all required GitHub Secrets for the automated deployment workflow.

## ⚙️ Self-Hosted Runner Requirements

This workflow uses **self-hosted GitHub Actions runners**. Ensure your runner has:

- ✅ **Docker** installed and running (`docker` and `docker compose` commands available)
- ✅ **Docker Buildx** plugin for multi-platform builds
- ✅ **Network access** to GitHub Container Registry (ghcr.io)
- ✅ **SSH access** to your dev VM
- ✅ Sufficient disk space for building Docker images (~5-10GB recommended)

**Verify runner prerequisites:**
```bash
docker --version
docker compose version
docker buildx version
```

## How to Add Secrets

1. Go to your GitHub repository
2. Click on **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret listed below

## Required Secrets

### 🖥️ VM Access Secrets

| Secret Name | Description | Example |
|------------|-------------|---------|
| `DEV_VM_HOST` | Dev VM IP address or hostname | `192.168.1.100` or `dev.example.com` |
| `DEV_VM_USER` | SSH username for VM access | `ubuntu` or `ec2-user` |
| `DEV_VM_SSH_KEY` | Private SSH key for VM authentication | Contents of your private key file |

**Note:** For `DEV_VM_SSH_KEY`, paste the entire contents of your private SSH key file (e.g., `~/.ssh/id_rsa`), including the `-----BEGIN RSA PRIVATE KEY-----` and `-----END RSA PRIVATE KEY-----` lines.

### 🗄️ Database Secrets (Supabase)

| Secret Name | Description | Example |
|------------|-------------|---------|
| `SUPABASE_URL` | Your Supabase project URL | `https://xxxxx.supabase.co` |
| `SUPABASE_API_KEY` | Supabase anon/public API key | `eyJhbGc...` |
| `SUPABASE_PUBLISHABLE_KEY` | Supabase publishable key | `eyJhbGc...` |
| `SUPABASE_SECRET_KEY` | Supabase service role secret key | `eyJhbGc...` |

**Where to find:** 
- Go to your Supabase project dashboard
- Navigate to **Settings** → **API**
- Copy the URL and keys

### 🔐 OAuth Secrets (Google)

| Secret Name | Description | Example |
|------------|-------------|---------|
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | `123456789-xxxxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | `GOCSPX-xxxxx` |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | `http://your-dev-vm-ip:8080/api/v1/auth/google/callback` |

**Where to find:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project or create a new one
3. Navigate to **APIs & Services** → **Credentials**
4. Create or select OAuth 2.0 Client ID
5. Add authorized redirect URIs

### 🔑 Session & Security

| Secret Name | Description | Example |
|------------|-------------|---------|
| `SESSION_SECRET` | Secret key for session encryption | Generate with: `openssl rand -base64 32` |

**How to generate:**
```bash
openssl rand -base64 32
```

### 📊 Monitoring Secrets

| Secret Name | Description | Example |
|------------|-------------|---------|
| `GRAFANA_PASSWORD` | Admin password for Grafana dashboard | `YourStrongPassword123!` |

**Note:** Use a strong password for production environments.

### 🌐 Frontend Configuration

| Secret Name | Description | Example |
|------------|-------------|---------|
| `FRONTEND_URL` | Public URL of your application (base URL) | `http://your-dev-vm-ip` or `https://spark-park-dev.dojima.foundation` |
| `ALLOWED_ORIGINS` | Comma-separated CORS allowed origins | `http://your-dev-vm-ip:3001,http://your-dev-vm-ip,http://localhost:3000` |

**Note:** `FRONTEND_URL` should be the base URL without port. The API endpoints will be constructed as `${FRONTEND_URL}/api/v1`.

## Secrets Checklist

Before running the deployment workflow, ensure you have added all these secrets:

- [ ] `DEV_VM_HOST`
- [ ] `DEV_VM_USER`
- [ ] `DEV_VM_SSH_KEY`
- [ ] `SUPABASE_URL`
- [ ] `SUPABASE_API_KEY`
- [ ] `SUPABASE_PUBLISHABLE_KEY`
- [ ] `SUPABASE_SECRET_KEY`
- [ ] `GOOGLE_CLIENT_ID`
- [ ] `GOOGLE_CLIENT_SECRET`
- [ ] `GOOGLE_REDIRECT_URL`
- [ ] `SESSION_SECRET`
- [ ] `GRAFANA_PASSWORD`
- [ ] `FRONTEND_URL`
- [ ] `ALLOWED_ORIGINS`

## Deployment Workflow

The deployment process:
1. Builds Docker images for backend and web in GitHub Actions
2. Pushes images to GitHub Container Registry (GHCR)
3. SSH into your dev VM
4. Copies deployment configurations from `deploy/environments/dev/`
5. Creates environment files from secrets
6. Pulls images from GHCR
7. Deploys all containers using `docker compose`
8. Verifies health of all services

## Testing Secrets

After adding all secrets, you can test the deployment by:

1. Go to **Actions** tab in your GitHub repository
2. Select **Deploy to Dev VM** workflow
3. Click **Run workflow**
4. Select the branch (dev) and click **Run workflow**
5. Monitor the deployment progress in the Actions logs

## 📈 Access URLs (After Deployment)

### Application Services
- **Backend API**: `http://your-dev-vm-ip:8081`
- **Frontend**: `http://your-dev-vm-ip:3001`
- **Nginx Reverse Proxy**: `http://your-dev-vm-ip` (port 80)
- **Health Check**: `http://your-dev-vm-ip:8081/health`

### Monitoring Stack
- **Grafana**: `http://your-dev-vm-ip:3003` (admin/your-grafana-password)
- **Prometheus**: `http://your-dev-vm-ip:9091`
- **Loki**: `http://your-dev-vm-ip:3101`
- **Redis**: Internal only (port 6380)

## Security Best Practices

1. **Never commit secrets** to your repository
2. **Rotate secrets regularly** (every 90 days recommended)
3. **Use strong passwords** for all secrets
4. **Limit secret access** to necessary team members only
5. **Use different secrets** for dev, staging, and production environments
6. **Monitor secret usage** in GitHub Actions logs
7. **Revoke and regenerate** secrets if they are compromised

## Troubleshooting

### SSH Connection Issues
- Ensure your SSH key is in the correct format (RSA, ED25519)
- Verify the SSH key has proper line breaks
- Check that the VM allows SSH connections from GitHub Actions IPs

### Database Connection Issues
- Verify Supabase URL includes `https://`
- Ensure API keys are copied completely without truncation
- Check that your Supabase project is active and accessible

### OAuth Issues
- Ensure redirect URL matches exactly (including protocol and port)
- Verify Google OAuth consent screen is configured
- Check that authorized redirect URIs are added in Google Cloud Console

## Support

For issues with:
- **GitHub Secrets**: Check GitHub documentation on [Encrypted Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- **Supabase**: Visit [Supabase Documentation](https://supabase.com/docs)
- **Google OAuth**: Check [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)





