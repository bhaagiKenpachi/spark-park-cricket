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

async function createDashboardWithInsight() {
    console.log(chalk.blue.bold('📊 Creating Dashboard with Insight\n'));

    const client = new PostHogAPIClient();

    try {
        // Step 1: Create a dashboard
        console.log(chalk.blue('1. Creating dashboard...'));

        const dashboardConfig = {
            name: 'Cricket Analytics Dashboard',
            description: 'Dashboard for cricket series analytics',
            pinned: true,
            show_tile_badges: true,
            tags: ['cricket', 'analytics', 'test']
        };

        const dashboard = await client.createDashboard(dashboardConfig);
        console.log(chalk.green('✅ Dashboard created successfully'));
        console.log(chalk.cyan(`   Dashboard ID: ${dashboard.id}`));
        console.log(chalk.cyan(`   Dashboard URL: ${client.getDashboardUrl(dashboard.id)}`));
        console.log('');

        // Step 2: Add the existing insight to the dashboard
        console.log(chalk.blue('2. Adding insight to dashboard...'));

        const insightId = 3802134; // The insight we created earlier
        console.log(chalk.blue(`   Adding insight ${insightId} to dashboard ${dashboard.id}`));

        const result = await client.addInsightToDashboard(dashboard.id, insightId);

        if (result) {
            console.log(chalk.green('✅ Insight added to dashboard successfully'));
        } else {
            console.log(chalk.yellow('⚠️  Insight could not be added automatically'));
            console.log(chalk.yellow('   You can add it manually in the PostHog UI'));
        }

        console.log(chalk.green.bold('\n🎉 Dashboard Setup Complete!'));
        console.log(chalk.cyan(`📊 Dashboard URL: ${client.getDashboardUrl(dashboard.id)}`));
        console.log(chalk.cyan(`📈 Insight URL: ${client.getInsightUrl(insightId)}`));
        console.log('');
        console.log(chalk.yellow('💡 Next Steps:'));
        console.log(chalk.yellow('   1. Go to the dashboard URL above'));
        console.log(chalk.yellow('   2. You should see your insight there'));
        console.log(chalk.yellow('   3. If not, click "Add insight" and search for your insight'));

        return { dashboard, insightId };

    } catch (error) {
        console.error(chalk.red('❌ Failed to create dashboard with insight:'));
        console.error(chalk.red(error.message));

        // Try to create just the dashboard without adding insight
        console.log(chalk.yellow('\n🔄 Trying to create dashboard only...'));

        try {
            const dashboardConfig = {
                name: 'Cricket Analytics Dashboard',
                description: 'Dashboard for cricket series analytics',
                pinned: true,
                tags: ['cricket', 'analytics', 'test']
            };

            const dashboard = await client.createDashboard(dashboardConfig);
            console.log(chalk.green('✅ Dashboard created successfully'));
            console.log(chalk.cyan(`   Dashboard ID: ${dashboard.id}`));
            console.log(chalk.cyan(`   Dashboard URL: ${client.getDashboardUrl(dashboard.id)}`));
            console.log('');
            console.log(chalk.yellow('💡 Manual Steps:'));
            console.log(chalk.yellow('   1. Go to the dashboard URL above'));
            console.log(chalk.yellow('   2. Click "Add insight" or "Edit dashboard"'));
            console.log(chalk.yellow('   3. Search for "Test Insight - API Connection"'));
            console.log(chalk.yellow('   4. Add it to the dashboard'));

            return { dashboard, insightId: 3802134 };
        } catch (dashboardError) {
            console.error(chalk.red('❌ Dashboard creation also failed:'));
            console.error(chalk.red(dashboardError.message));
            return null;
        }
    }
}

// Run the script
if (require.main === module) {
    createDashboardWithInsight();
}

module.exports = { createDashboardWithInsight };
