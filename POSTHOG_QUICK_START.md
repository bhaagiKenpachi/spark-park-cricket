# PostHog Quick Start Guide

Get PostHog analytics running in your Spark Park Cricket app in 5 minutes!

## Prerequisites

- Docker and Docker Compose installed
- Spark Park Cricket project cloned

## Step 1: Start PostHog (1 minute)

```bash
# Navigate to project root
cd /path/to/spark-park-cricket

# Start all services including PostHog
docker compose up -d

# Verify PostHog is running (wait ~30 seconds for full startup)
docker compose ps | grep posthog
```

You should see 5 PostHog containers running:
- cricket-posthog (web)
- cricket-posthog-worker
- cricket-posthog-db
- cricket-posthog-clickhouse
- cricket-posthog-redis

## Step 2: Access PostHog Dashboard (1 minute)

1. Open browser: http://localhost:8001
2. Create account (local, no email required)
3. Set up organization (any name)
4. Create project: "Spark Park Cricket - Local"

## Step 3: Get Your API Key (1 minute)

1. In PostHog dashboard, click Settings (gear icon)
2. Go to "Project" tab
3. Copy the "Project API Key" (starts with `phc_`)

Example: `phc_abc123xyz789...`

## Step 4: Configure Frontend (1 minute)

**Option A: Using Docker Compose (Recommended)**

Edit `docker-compose.yml`:

```yaml
web:
  environment:
    - NEXT_PUBLIC_POSTHOG_KEY=phc_your_actual_key_here  # Replace this
    - NEXT_PUBLIC_POSTHOG_HOST=http://localhost:8001
    - NEXT_PUBLIC_POSTHOG_DEBUG=true
```

Restart frontend:
```bash
docker compose restart web
```

**Option B: Running Locally**

Create `web/.env.local`:

```env
NEXT_PUBLIC_POSTHOG_KEY=phc_your_actual_key_here
NEXT_PUBLIC_POSTHOG_HOST=http://localhost:8001
NEXT_PUBLIC_POSTHOG_DEBUG=true
```

Start dev server:
```bash
cd web
npm run dev
```

## Step 5: Verify It's Working (1 minute)

1. Open your app: http://localhost:3000
2. Navigate around (view series, matches, etc.)
3. Go back to PostHog dashboard
4. Click "Events" in sidebar
5. You should see events appearing!

Look for:
- `series_viewed`
- `$pageview`
- `match_viewed` (if you clicked a match)

### Browser Console Check

Open browser console (F12), type:
```javascript
window.posthog
```

Should return PostHog object (not undefined).

Enable debug to see all events:
```javascript
posthog.debug()
```

## Troubleshooting

### PostHog Dashboard Not Loading

```bash
# Check logs
docker compose logs posthog

# Restart PostHog
docker compose restart posthog
```

### Events Not Showing

1. **Check API Key**: Verify it's correct in environment variables
2. **Check Network**: Open DevTools > Network, filter by "posthog"
3. **Wait**: Events may take 10-30 seconds to appear
4. **Restart Frontend**: `docker compose restart web`

### "PostHog not initialized" Error

- Make sure you restarted the frontend after adding the API key
- Check browser console for errors
- Verify PostHogProvider is in layout.tsx (it is!)

## What's Being Tracked?

Your app is now tracking:

✅ **Page Views**: Every page you visit
✅ **Series Events**: View, create, edit, delete
✅ **Match Events**: View, create, delete
✅ **Scorecard Events**: View, live scoring, ball-by-ball
✅ **Auth Events**: Login, logout, user identification
✅ **Session Recordings**: Full user sessions (video replay)

## Next Steps

### 1. Explore PostHog Dashboard

- **Events**: See all tracked events
- **Insights**: Create charts and analytics
- **Session Recordings**: Watch user sessions
- **Feature Flags**: Create A/B tests
- **Persons**: See identified users

### 2. Create Your First Insight

1. Go to "Insights" > "New Insight"
2. Select event: `series_viewed`
3. Choose visualization: "Line chart"
4. Set date range: "Last 7 days"
5. Click "Save"

### 3. Test Session Recording

1. Go to "Session Recordings"
2. Your recent session should appear
3. Click to watch replay
4. See all your interactions!

### 4. Create Feature Flag

1. Go to "Feature Flags" > "New Feature Flag"
2. Name: `new-scorecard-ui`
3. Enable for: "100% of users"
4. Save

Use in code:
```typescript
import { useFeatureFlag } from '@/hooks/useFeatureFlag';

const { enabled } = useFeatureFlag('new-scorecard-ui');
if (enabled) {
  // Show new UI
}
```

## Testing Checklist

Test these to verify everything works:

- [ ] PostHog dashboard accessible at localhost:8001
- [ ] Events appearing in PostHog Events tab
- [ ] Create a series (check for `series_created` event)
- [ ] View matches (check for `match_viewed` event)
- [ ] Login (check for `user_logged_in` event)
- [ ] Session recording captured
- [ ] No errors in browser console

## Production Deployment

Once tested locally, deploy to production:

1. **Read**: [PostHog Deployment Guide](web/docs/posthog-deployment.md)
2. **Provision VM**: Recommended Hetzner CX21 (€4.90/month)
3. **Deploy PostHog**: Follow deployment guide
4. **Update App**: Use production PostHog URL
5. **Monitor**: Check analytics in production

## Need Help?

- **Testing Guide**: [web/docs/posthog-testing.md](web/docs/posthog-testing.md)
- **Deployment Guide**: [web/docs/posthog-deployment.md](web/docs/posthog-deployment.md)
- **Environment Setup**: [web/ENV_SETUP.md](web/ENV_SETUP.md)
- **Full Summary**: [POSTHOG_INTEGRATION_SUMMARY.md](POSTHOG_INTEGRATION_SUMMARY.md)

## Common Commands

```bash
# View PostHog logs
docker compose logs -f posthog

# Restart PostHog
docker compose restart posthog

# Stop PostHog
docker compose stop posthog posthog-worker posthog-db posthog-clickhouse posthog-redis

# Start PostHog
docker compose up -d posthog posthog-worker posthog-db posthog-clickhouse posthog-redis

# Full restart (if issues)
docker compose down
docker compose up -d
```

## Success! 🎉

If you can see events in PostHog dashboard, you're all set!

Your cricket app now has:
- ✅ Complete user analytics
- ✅ Session recordings
- ✅ Feature flags
- ✅ User identification
- ✅ Real-time insights

Start making data-driven decisions to improve your app!


