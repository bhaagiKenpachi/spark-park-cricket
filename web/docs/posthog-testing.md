# PostHog Testing Guide

This guide explains how to test PostHog analytics locally using Docker Compose and verify that events are being tracked correctly.

## Table of Contents

1. [Local Setup with Docker Compose](#local-setup-with-docker-compose)
2. [Event Testing Checklist](#event-testing-checklist)
3. [Feature Flags Testing](#feature-flags-testing)
4. [Session Recording Verification](#session-recording-verification)
5. [Troubleshooting](#troubleshooting)

## Local Setup with Docker Compose

### Step 1: Start PostHog Services

```bash
# From the project root directory
cd /path/to/spark-park-cricket

# Start all services including PostHog
docker compose up -d

# Check if PostHog is running
docker compose ps
```

You should see these PostHog containers running:
- `cricket-posthog` (web service on port 8001)
- `cricket-posthog-worker` (background worker)
- `cricket-posthog-db` (PostgreSQL)
- `cricket-posthog-clickhouse` (analytics database)
- `cricket-posthog-redis` (queue)

### Step 2: Access PostHog Dashboard

1. Open your browser and go to `http://localhost:8001`
2. Create a new account (this is local, no email required)
3. Set up your organization and project
4. Copy your **Project API Key** from Settings > Project Settings

### Step 3: Update Environment Variables

Update the PostHog key in `docker-compose.yml`:

```yaml
web:
  environment:
    - NEXT_PUBLIC_POSTHOG_KEY=<your-project-api-key>  # Replace with actual key
    - NEXT_PUBLIC_POSTHOG_HOST=http://localhost:8001
```

Or create a `.env.local` file in the `web` directory:

```env
NEXT_PUBLIC_POSTHOG_KEY=phc_your_actual_project_key
NEXT_PUBLIC_POSTHOG_HOST=http://localhost:8001
NEXT_PUBLIC_POSTHOG_DEBUG=true
```

### Step 4: Restart Frontend

```bash
docker compose restart web

# Or if running locally
cd web
npm run dev
```

## Event Testing Checklist

### Series Events

- [ ] **Series Viewed**: Navigate to the homepage
  - Go to PostHog > Events
  - Look for `series_viewed` events
  - Verify properties: `series_id`, `series_name`, `start_date`, `end_date`

- [ ] **Series Created**: Click "Create Series" and submit form
  - Look for `series_created` event
  - Verify series details are captured

- [ ] **Series Edited**: Edit an existing series
  - Look for `series_edited` event

- [ ] **Series Deleted**: Delete a series
  - Look for `series_deleted` event

- [ ] **Pagination Changed**: Change page or page size
  - Look for `series_pagination_changed` event
  - Verify properties: `page`, `page_size`, `total_items`

### Match Events

- [ ] **Match Viewed**: Click on a match to view scorecard
  - Look for `match_viewed` event
  - Verify properties: `match_id`, `series_id`, `match_number`, `match_status`

- [ ] **Match Created**: Create a new match in a series
  - Look for `match_created` event

- [ ] **Match Deleted**: Delete a match
  - Look for `match_deleted` event

### Scorecard Events

- [ ] **Scorecard Viewed**: View any match scorecard
  - Look for `scorecard_viewed` event
  - Verify properties: `match_id`, `innings_number`, `match_status`

- [ ] **Live Scoring Started**: Start live scoring on a match
  - Look for `live_scoring_started` event

- [ ] **Ball Added**: Score a ball (run, wide, wicket, etc.)
  - Look for `ball_added` event
  - Verify properties: `match_id`, `innings_number`, `over_number`, `ball_number`, `ball_type`, `runs`, `is_wicket`

- [ ] **Over Completed**: Complete 6 legal deliveries
  - Look for `over_completed` event (if implemented)

- [ ] **Innings Completed**: Complete an innings
  - Look for `innings_completed` event

### Authentication Events

- [ ] **User Logged In**: Sign in with Google
  - Look for `user_logged_in` event
  - Verify user properties: `user_id`, `email`, `name`, `provider`
  - Check that user is identified in PostHog (see Persons)

- [ ] **User Logged Out**: Sign out
  - Look for `user_logged_out` event
  - Verify user session is reset

### Page View Events

- [ ] **Automatic Page Views**: Navigate between pages
  - Look for `$pageview` events
  - Verify `$current_url` property is correct

## Feature Flags Testing

### Step 1: Create Feature Flags in PostHog

1. Go to PostHog Dashboard > Feature Flags
2. Click "New Feature Flag"
3. Create test flags:
   - `new-scorecard-ui` (boolean)
   - `enable-ads` (boolean)
   - `beta-features` (boolean)

### Step 2: Test Flag Behavior

```typescript
// Example usage in a component
import { useFeatureFlag } from '@/hooks/useFeatureFlag';

const { enabled, loading } = useFeatureFlag('new-scorecard-ui', false);

if (loading) return <div>Loading...</div>;
if (enabled) return <NewUI />;
return <OldUI />;
```

### Step 3: Verify in PostHog

1. Enable/disable flags in PostHog dashboard
2. Refresh your app
3. Verify the UI changes according to flag state
4. Check "Feature Flag Calls" in PostHog Insights

## Session Recording Verification

### Step 1: Enable Session Recording

Session recording is enabled by default in the PostHog configuration. Verify in `src/lib/posthog.ts`:

```typescript
session_recording: {
  recordCrossOriginIframes: true,
  maskAllInputs: false,
  maskInputOptions: {
    password: true,
  },
}
```

### Step 2: Generate Test Sessions

1. Navigate through your app
2. Create a series
3. Create a match
4. Start live scoring
5. Add some balls

### Step 3: View Recordings

1. Go to PostHog Dashboard > Session Recordings
2. Find your recent sessions
3. Click to view playback
4. Verify you can see:
   - Mouse movements
   - Clicks
   - Page navigations
   - Form inputs (non-passwords)

### Tips for Testing Recordings

- **Console Logs**: Check browser console for PostHog messages
- **Network Tab**: Verify recording data is being sent to PostHog
- **Recording Quality**: Ensure clicks and interactions are captured
- **Privacy**: Verify sensitive data (passwords) are masked

## Troubleshooting

### Events Not Showing Up

1. **Check PostHog Initialization**
   ```javascript
   // Open browser console
   console.log(window.posthog);
   // Should show PostHog object, not undefined
   ```

2. **Verify API Key**
   ```javascript
   console.log(process.env.NEXT_PUBLIC_POSTHOG_KEY);
   // Should show your project key starting with 'phc_'
   ```

3. **Check Network Requests**
   - Open DevTools > Network tab
   - Filter by "posthog" or "capture"
   - Look for POST requests to `/e/` or `/decide/`
   - Status should be 200

4. **Enable Debug Mode**
   ```env
   NEXT_PUBLIC_POSTHOG_DEBUG=true
   ```
   Then check browser console for PostHog debug messages

### PostHog Container Issues

```bash
# Check container logs
docker compose logs posthog
docker compose logs posthog-worker

# Restart PostHog services
docker compose restart posthog posthog-worker

# Check health status
docker compose ps
```

### Data Not Appearing in Dashboard

- **Wait Time**: Events may take 10-30 seconds to appear
- **Refresh**: Try refreshing the PostHog dashboard
- **Date Range**: Check the date filter in PostHog
- **Project**: Ensure you're viewing the correct project

### Session Recordings Not Working

1. **Check Browser Compatibility**: Works best on Chrome/Firefox
2. **Check CSP Headers**: Verify `next.config.ts` allows PostHog domains
3. **Storage**: Ensure localStorage is enabled
4. **Network**: Check for blocked requests in DevTools

### Common Errors

**Error: "PostHog not initialized"**
- Solution: Check that PostHogProvider is wrapping your app in layout.tsx

**Error: "Invalid API key"**
- Solution: Verify your API key in environment variables

**Error: "CORS policy"**
- Solution: Check `next.config.ts` CSP headers include PostHog domains

**Error: "Cannot read property 'capture'"**
- Solution: Ensure PostHog is initialized before tracking events

## Testing Best Practices

1. **Use Debug Mode**: Always test with debug mode enabled locally
2. **Check Multiple Event Types**: Test at least one event from each category
3. **Verify User Identification**: Ensure logged-in users are properly identified
4. **Test Feature Flags**: Toggle flags and verify UI changes
5. **Review Session Recordings**: Watch at least one full session
6. **Test on Multiple Browsers**: Chrome, Firefox, Safari
7. **Test Mobile**: Use Chrome DevTools mobile emulation

## Advanced Testing

### Test Event Properties

```javascript
// In browser console
posthog.capture('test_event', {
  test_property: 'test_value',
  timestamp: new Date().toISOString()
});
```

### Test User Properties

```javascript
posthog.identify('test_user_123', {
  name: 'Test User',
  email: 'test@example.com'
});
```

### Test Groups

```javascript
posthog.group('team', 'test_team_123', {
  name: 'Test Team'
});
```

## Automated Testing

For CI/CD, you can mock PostHog in tests:

```typescript
// jest.setup.js
jest.mock('@/lib/posthog', () => ({
  posthog: {
    capture: jest.fn(),
    identify: jest.fn(),
    reset: jest.fn(),
    isFeatureEnabled: jest.fn(() => false),
  },
  initPostHog: jest.fn(),
}));
```

## Next Steps

After successful local testing:
1. Set up PostHog on development VM instance
2. Test with dev environment
3. Set up PostHog on production VM instance
4. Deploy to production
5. Monitor events in production dashboard


