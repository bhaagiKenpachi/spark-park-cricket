# PostHog Dashboard Management

This directory contains tools to programmatically create and manage PostHog dashboards for the Cricket Analytics application.

## 🚀 Quick Start

### 1. Get Your PostHog Credentials

#### Personal API Key
1. Go to [PostHog Settings > Personal API Keys](https://us.i.posthog.com/project/settings#api-keys)
2. Click "Create Personal API Key"
3. Give it a name like "Dashboard Management"
4. Copy the generated key (should start with `phx_` or similar, NOT `phc_`)

**Important**: You need a Personal API Key (starts with `phx_`), not a Project API Key (starts with `phc_`). The Project API Key is used for sending events, but the Personal API Key is required for dashboard management.

#### Project ID
1. Go to your PostHog project dashboard
2. Look at the URL: `https://us.i.posthog.com/project/12345/dashboard/...`
3. The number after `/project/` is your Project ID (e.g., `12345`)

### 2. Setup Environment Variables

Add your PostHog credentials to the main `.env` file in the web directory:

```bash
# Edit the main .env file
nano .env
```

Add these variables to your existing `.env` file:
```bash
# PostHog Dashboard Management
POSTHOG_PERSONAL_API_KEY=phc_your_personal_api_key_here
POSTHOG_PROJECT_ID=12345
POSTHOG_HOST=https://us.i.posthog.com
```

### 3. Install Dependencies

```bash
# From the web directory
npm install
```

### 4. Create Series Analytics Dashboard

```bash
# Create the series dashboard
npm run posthog:dashboard:create

# Or run directly
node scripts/posthog/create-series-dashboard.js
```

## 📊 Available Commands

### Dashboard Management

```bash
# Create Series Analytics Dashboard
npm run posthog:dashboard:create

# List all dashboards
npm run posthog:dashboard:manage list

# Delete a dashboard
npm run posthog:dashboard:manage delete <dashboard_id>

# Show configuration info
npm run posthog:dashboard:manage info
```

### Direct CLI Usage

```bash
# Create series dashboard
node scripts/posthog/manage-dashboard.js create-series

# List dashboards
node scripts/posthog/manage-dashboard.js list

# Delete dashboard
node scripts/posthog/manage-dashboard.js delete 12345

# Show info
node scripts/posthog/manage-dashboard.js info
```

## 🏏 Series Analytics Dashboard

The Series Analytics Dashboard includes 6 predefined insights:

1. **Series Views Over Time** - Line chart showing series views over 30 days
2. **Series Created** - Bar chart tracking new series creation
3. **Series by Name** - Table showing most viewed series
4. **Series Edits** - Line chart tracking series editing activity
5. **Users Creating Series** - Unique count of users creating series
6. **Series Engagement** - Combined view of all series activities

## 🔄 Auto-Updating Dashboards

Once created, dashboards automatically update as new events are captured:
- No manual refresh needed
- Similar to Grafana's live queries
- Real-time data visualization
- Historical trend analysis

## 🛠️ Troubleshooting

### Authentication Errors
```
❌ API Request failed: API Error: 401 - Authentication failed
```
**Solution:** Check your Personal API Key in the `.env` file

### Project Not Found
```
❌ API Request failed: API Error: 404 - Project not found
```
**Solution:** Verify your Project ID in the `.env` file

### Missing Environment Variables
```
❌ Missing required environment variables
```
**Solution:** Add the required variables to your main `.env` file in the web directory:
```bash
POSTHOG_PERSONAL_API_KEY=phc_your_personal_api_key_here
POSTHOG_PROJECT_ID=12345
POSTHOG_HOST=https://us.i.posthog.com
```

### Network Errors
```
❌ API Request failed: fetch failed
```
**Solution:** Check your internet connection and PostHog host URL

## 📁 File Structure

```
scripts/posthog/
├── api-client.js         # PostHog API client
├── create-series-dashboard.js  # Series dashboard creator
├── manage-dashboard.js   # CLI tool
└── README.md            # This file

web/
├── .env                  # Main environment file (add your credentials here)
└── package.json         # NPM scripts for dashboard management
```

## 🔗 Useful Links

- [PostHog Personal API Keys](https://us.i.posthog.com/project/settings#api-keys)
- [PostHog API Documentation](https://posthog.com/docs/api)
- [PostHog Dashboard Documentation](https://posthog.com/docs/user-guides/dashboards)

## 💡 Tips

- Dashboards are created in your PostHog project and accessible via the web UI
- All insights auto-update with new data - no manual refresh needed
- You can customize insights after creation through the PostHog UI
- Use the `list` command to see all your dashboards
- Use the `delete` command to clean up test dashboards
