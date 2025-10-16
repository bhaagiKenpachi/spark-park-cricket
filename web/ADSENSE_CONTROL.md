# AdSense Enable/Disable Control

## Overview

You can now enable or disable Google AdSense ads across your entire application using a single environment variable. This is useful for:

- **Development/Testing**: Disable ads while developing or testing
- **A/B Testing**: Test user engagement with/without ads
- **Temporary Disabling**: Quickly disable ads for maintenance or special events
- **Regional Control**: Different settings for dev/staging/production environments

## Environment Variable

### `NEXT_PUBLIC_ENABLE_ADS`

- **Type**: String (evaluated as boolean)
- **Default**: `true` (ads enabled)
- **Values**:
  - `"true"` or omitted - Ads **enabled** ✅
  - `"false"` - Ads **disabled** ❌

## Usage

### 1. Local Development (`.env.local`)

```bash
# Enable ads (default)
NEXT_PUBLIC_ENABLE_ADS=true

# Disable ads
NEXT_PUBLIC_ENABLE_ADS=false
```

**Example `.env.local` file:**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/api/v1/graphql
NEXT_PUBLIC_ADSENSE_CLIENT_ID=ca-pub-5474524579770573
NEXT_PUBLIC_ENABLE_ADS=false  # Disable ads in development
```

### 2. Docker Build

```bash
# Build with ads enabled (default)
docker build \
  --build-arg NEXT_PUBLIC_ENABLE_ADS=true \
  -t my-cricket-app .

# Build with ads disabled
docker build \
  --build-arg NEXT_PUBLIC_ENABLE_ADS=false \
  -t my-cricket-app .
```

### 3. GitHub Actions (Secrets)

Add the secret to your GitHub repository:

```bash
# Using GitHub CLI
gh secret set NEXT_PUBLIC_ENABLE_ADS --repo bhaagiKenpachi/spark-park-cricket --body "true"

# For specific environments
gh secret set NEXT_PUBLIC_ENABLE_ADS --repo bhaagiKenpachi/spark-park-cricket --env dev --body "false"
gh secret set NEXT_PUBLIC_ENABLE_ADS --repo bhaagiKenpachi/spark-park-cricket --env prod --body "true"
```

**Via GitHub Dashboard:**
1. Go to repository **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret** (or environment-specific secret)
3. Name: `NEXT_PUBLIC_ENABLE_ADS`
4. Value: `true` or `false`
5. Click **Add secret**

### 4. Production Deployment

The deployment workflows automatically use the secret:

```yaml
# In deploy-dev.yml or deploy-prod.yml
--build-arg NEXT_PUBLIC_ENABLE_ADS=${{ secrets.NEXT_PUBLIC_ENABLE_ADS || 'true' }}
```

If the secret is not set, it defaults to `true` (ads enabled).

## How It Works

When `NEXT_PUBLIC_ENABLE_ADS=false`:

1. **AdSenseScript Component**
   - Script tag is **not loaded**
   - Logs: `"Ads are disabled via NEXT_PUBLIC_ENABLE_ADS environment variable."`

2. **ResponsiveAd Component**
   - Returns empty fragment `<></>`
   - No ad slot rendered

3. **InFeedAd Component**
   - Returns empty fragment `<></>`
   - No ads between series items

4. **OverAdModal Component**
   - Returns empty fragment `<></>`
   - No modal ads after overs

## Testing

### Test Ads Disabled

1. **Set environment variable:**
   ```bash
   # In .env.local
   NEXT_PUBLIC_ENABLE_ADS=false
   ```

2. **Restart development server:**
   ```bash
   npm run dev
   ```

3. **Verify in browser:**
   - Open DevTools → Console
   - Should see: `"Ads are disabled via NEXT_PUBLIC_ENABLE_ADS environment variable."`
   - No AdSense script loaded in Network tab
   - No ad slots visible on page

### Test Ads Enabled

1. **Set environment variable:**
   ```bash
   # In .env.local
   NEXT_PUBLIC_ENABLE_ADS=true
   ```

2. **Restart development server:**
   ```bash
   npm run dev
   ```

3. **Verify in browser:**
   - AdSense script loads in Network tab
   - Ad slots rendered on page (may be blank until approval)
   - No console messages about disabled ads

## Production Scenarios

### Scenario 1: Disable Ads During Maintenance

```bash
# Set GitHub secret to false
gh secret set NEXT_PUBLIC_ENABLE_ADS --repo bhaagiKenpachi/spark-park-cricket --env prod --body "false"

# Trigger deployment
git push origin main
```

### Scenario 2: Enable Ads After Testing

```bash
# Set GitHub secret to true
gh secret set NEXT_PUBLIC_ENABLE_ADS --repo bhaagiKenpachi/spark-park-cricket --env prod --body "true"

# Trigger deployment
git push origin main
```

### Scenario 3: Different Settings per Environment

```bash
# Dev: Ads disabled (for development)
gh secret set NEXT_PUBLIC_ENABLE_ADS --env dev --body "false"

# Prod: Ads enabled (for revenue)
gh secret set NEXT_PUBLIC_ENABLE_ADS --env prod --body "true"
```

## Environment Configuration Matrix

| Environment | Variable Value | Result | Use Case |
|-------------|---------------|--------|----------|
| **Local Dev** | `false` | ❌ Ads disabled | Development without distractions |
| **Local Dev** | `true` | ✅ Ads enabled | Test ad integration locally |
| **Dev/Staging** | `false` | ❌ Ads disabled | Testing without ads |
| **Dev/Staging** | `true` | ✅ Ads enabled | Test ad functionality in staging |
| **Production** | `true` | ✅ Ads enabled | Generate revenue |
| **Production** | `false` | ❌ Ads disabled | Temporary maintenance mode |

## Troubleshooting

### Issue: Ads still showing after setting to `false`

**Solution:**
1. Verify the environment variable is set correctly
2. Restart the development server (if local)
3. Clear browser cache
4. Check the console for the "Ads are disabled" message

### Issue: Ads not showing after setting to `true`

**Solution:**
1. Verify `NEXT_PUBLIC_ADSENSE_CLIENT_ID` is also set
2. Check that the domain is approved in Google AdSense
3. Wait for AdSense approval (24-48 hours for new domains)
4. Check browser console for errors

### Issue: Environment variable not being picked up in Docker

**Solution:**
1. Verify the build arg is passed in the Docker build command
2. Rebuild the Docker image (don't use cached layers):
   ```bash
   docker build --no-cache --build-arg NEXT_PUBLIC_ENABLE_ADS=true -t my-app .
   ```

### Issue: GitHub Actions deployment not respecting the setting

**Solution:**
1. Verify the secret is set in the correct environment (dev/prod)
2. Check the workflow logs to see what value was used
3. Ensure the secret name matches exactly: `NEXT_PUBLIC_ENABLE_ADS`

## Best Practices

1. **Default to Enabled**: Keep ads enabled in production unless there's a specific reason to disable
2. **Disable in Development**: Set to `false` in `.env.local` for cleaner development experience
3. **Use Environment Secrets**: Store in GitHub secrets for production control
4. **Document Changes**: Comment in code or commits when changing this setting
5. **Monitor Revenue**: Track AdSense dashboard when toggling ads
6. **Test Before Production**: Always test ad changes in staging/dev first

## Related Files

- `web/src/components/ads/AdSenseScript.tsx` - Main script loader
- `web/src/components/ads/ResponsiveAd.tsx` - Responsive display ads
- `web/src/components/ads/InFeedAd.tsx` - In-feed ads
- `web/src/components/ads/OverAdModal.tsx` - Modal ads after overs
- `web/Dockerfile` - Docker build configuration
- `.github/workflows/deploy-dev.yml` - Dev deployment workflow
- `.github/workflows/deploy-prod.yml` - Prod deployment workflow

## Additional Resources

- [Next.js Environment Variables](https://nextjs.org/docs/pages/building-your-application/configuring/environment-variables)
- [Docker Build Arguments](https://docs.docker.com/engine/reference/builder/#arg)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)

