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

async function cleanupAll() {
    console.log(chalk.red.bold('🧹 Cleaning up all PostHog resources\n'));

    const client = new PostHogAPIClient();

    try {
        // Step 1: List all dashboards
        console.log(chalk.blue('1. Finding dashboards to delete...'));
        const dashboardsResponse = await client.makeRequest('/dashboards/');
        const seriesDashboards = dashboardsResponse.results?.filter(dashboard =>
            dashboard.name.includes('Series') || dashboard.name.includes('Cricket') || dashboard.name.includes('Analytics')
        ) || [];

        console.log(chalk.yellow(`Found ${seriesDashboards.length} dashboards to delete:`));
        seriesDashboards.forEach(dashboard => {
            console.log(chalk.gray(`   • ${dashboard.name} (ID: ${dashboard.id})`));
        });
        console.log('');

        // Step 2: Delete dashboards
        for (const dashboard of seriesDashboards) {
            try {
                console.log(chalk.blue(`🗑️  Deleting dashboard: ${dashboard.name} (ID: ${dashboard.id})`));
                await client.deleteDashboard(dashboard.id);
                console.log(chalk.green(`✅ Deleted dashboard: ${dashboard.name}`));
            } catch (error) {
                console.log(chalk.red(`❌ Failed to delete dashboard ${dashboard.name}: ${error.message}`));
            }
        }
        console.log('');

        // Step 3: List all insights
        console.log(chalk.blue('2. Finding insights to delete...'));
        const insightsResponse = await client.makeRequest('/insights/');
        const seriesInsights = insightsResponse.results?.filter(insight =>
            insight.name.includes('Series') || insight.name.includes('series') ||
            insight.name.includes('Cricket') || insight.name.includes('Analytics')
        ) || [];

        console.log(chalk.yellow(`Found ${seriesInsights.length} insights to delete:`));
        seriesInsights.forEach(insight => {
            console.log(chalk.gray(`   • ${insight.name} (ID: ${insight.id})`));
        });
        console.log('');

        // Step 4: Delete insights
        for (const insight of seriesInsights) {
            try {
                console.log(chalk.blue(`🗑️  Deleting insight: ${insight.name} (ID: ${insight.id})`));
                await client.makeRequest(`/insights/${insight.id}/`, { method: 'DELETE' });
                console.log(chalk.green(`✅ Deleted insight: ${insight.name}`));
            } catch (error) {
                console.log(chalk.red(`❌ Failed to delete insight ${insight.name}: ${error.message}`));
            }
        }
        console.log('');

        console.log(chalk.green.bold('🎉 Cleanup completed!'));
        console.log(chalk.yellow('💡 All series-related dashboards and insights have been deleted.'));
        console.log(chalk.yellow('   You can now start fresh with a clean PostHog project.'));

    } catch (error) {
        console.error(chalk.red('❌ Cleanup failed:'));
        console.error(chalk.red(error.message));
    }
}

// Run the script
if (require.main === module) {
    cleanupAll();
}

module.exports = { cleanupAll };
