#!/usr/bin/env node

const chalk = require('chalk');
const path = require('path');
const fs = require('fs');

// Try to load .env.local first, then .env
const envLocalPath = path.join(__dirname, '../../.env.local');
const envPath = path.join(__dirname, '../../.env');

if (fs.existsSync(envLocalPath)) {
    require('dotenv').config({ path: envLocalPath });
} else if (fs.existsSync(envPath)) {
    require('dotenv').config({ path: envPath });
}

const PostHogAPIClient = require('./api-client');

// Series insights that were created
const SERIES_INSIGHT_IDS = [
    { id: 3801439, name: 'Series Views Over Time' },
    { id: 3801440, name: 'Series Created' },
    { id: 3801441, name: 'Series by Name' },
    { id: 3801442, name: 'Series Edits' },
    { id: 3801443, name: 'Users Creating Series' },
    { id: 3801444, name: 'Series Engagement' }
];

async function addInsightsToDashboard() {
    console.log(chalk.blue.bold('🔗 Adding Series Insights to Dashboard\n'));

    const client = new PostHogAPIClient();
    const dashboardId = 615755; // The dashboard we just created

    try {
        console.log(chalk.blue(`📋 Working with dashboard ID: ${dashboardId}`));
        console.log(chalk.cyan(`🔗 Dashboard URL: ${client.getDashboardUrl(dashboardId)}\n`));

        // Get current dashboard
        const dashboard = await client.makeRequest(`/dashboards/${dashboardId}/`);
        const tiles = dashboard.tiles || [];

        console.log(chalk.blue(`📊 Current dashboard has ${tiles.length} tiles`));

        // Add each insight to the dashboard
        let addedCount = 0;
        for (const insight of SERIES_INSIGHT_IDS) {
            try {
                console.log(chalk.blue(`🔗 Adding insight: ${insight.name} (ID: ${insight.id})`));

                // Add insight to tiles
                tiles.push({
                    insight: insight.id,
                    layouts: {
                        sm: { x: (tiles.length % 2) * 6, y: Math.floor(tiles.length / 2) * 4, w: 6, h: 4 },
                        lg: { x: (tiles.length % 3) * 4, y: Math.floor(tiles.length / 3) * 4, w: 4, h: 4 }
                    }
                });

                // Update dashboard with new tiles
                await client.makeRequest(`/dashboards/${dashboardId}/`, {
                    method: 'PATCH',
                    body: JSON.stringify({ tiles }),
                });

                console.log(chalk.green(`✅ Added ${insight.name} to dashboard`));
                addedCount++;

            } catch (error) {
                console.log(chalk.yellow(`⚠️  Could not add ${insight.name}: ${error.message}`));
            }
        }

        console.log(chalk.green.bold(`\n🎉 Successfully added ${addedCount} insights to the dashboard!`));
        console.log(chalk.cyan(`📊 Dashboard URL: ${client.getDashboardUrl(dashboardId)}`));
        console.log(chalk.yellow('\n💡 The dashboard will now show all series analytics insights.'));

    } catch (error) {
        console.error(chalk.red('❌ Failed to add insights to dashboard:'));
        console.error(chalk.red(error.message));
    }
}

// Run the script
if (require.main === module) {
    addInsightsToDashboard();
}

module.exports = { addInsightsToDashboard };
