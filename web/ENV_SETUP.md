# Environment Variables Setup for PostHog

## Required Environment Variables

Create a `.env.local` file in the `web` directory with the following variables:

```env
# Backend API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/api/v1/graphql
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws

# PostHog Analytics Configuration
# Get your API key from PostHog dashboard after creating a project
NEXT_PUBLIC_POSTHOG_KEY=phc_your_project_api_key_here
NEXT_PUBLIC_POSTHOG_HOST=http://localhost:8001
NEXT_PUBLIC_POSTHOG_DEBUG=false

# Ads Configuration (optional)
NEXT_PUBLIC_ENABLE_ADS=false
```

## For Docker Compose

PostHog environment variables are already configured in `docker-compose.yml` under the `web` service:

```yaml
# PostHog Configuration
- NEXT_PUBLIC_POSTHOG_KEY=phc_test_key_placeholder
- NEXT_PUBLIC_POSTHOG_HOST=http://localhost:8001
- NEXT_PUBLIC_POSTHOG_DEBUG=false
```

After starting PostHog with Docker Compose, you'll need to:
1. Access PostHog at `http://localhost:8001`
2. Create an account and project
3. Get your project API key from Settings
4. Update the `NEXT_PUBLIC_POSTHOG_KEY` in docker-compose.yml or your .env.local

## For Development Environment

```env
NEXT_PUBLIC_POSTHOG_KEY=your_dev_project_key
NEXT_PUBLIC_POSTHOG_HOST=https://posthog-dev.yourdomain.com
NEXT_PUBLIC_POSTHOG_DEBUG=true
```

## For Production Environment

```env
NEXT_PUBLIC_POSTHOG_KEY=your_prod_project_key
NEXT_PUBLIC_POSTHOG_HOST=https://posthog.yourdomain.com
NEXT_PUBLIC_POSTHOG_DEBUG=false
```

## Notes

- The `NEXT_PUBLIC_` prefix is required for Next.js to expose these variables to the browser
- PostHog KEY should start with `phc_`
- Debug mode should only be enabled in development
- For self-hosted PostHog, the HOST should point to your PostHog instance URL


