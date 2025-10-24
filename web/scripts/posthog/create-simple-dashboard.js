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

async function createSimpleDashboard() {
    console.log(chalk.blue.bold('📊 Creating Simple Series Dashboard\n'));

    const client = new PostHogAPIClient();

    try {
        // Create a simple dashboard
        const dashboardConfig = {
            name: 'Cricket Series Analytics',
            description: 'Analytics dashboard for cricket series management',
            pinned: true
        };

        console.log(chalk.blue('📋 Creating dashboard...'));
        const dashboard = await client.createDashboard(dashboardConfig);

        console.log(chalk.green(`✅ Dashboard created: ${dashboard.name} (ID: ${dashboard.id})`));
        console.log(chalk.cyan(`🔗 Dashboard URL: ${client.getDashboardUrl(dashboard.id)}`));

        return dashboard;

    } catch (error) {
        console.error(chalk.red('❌ Failed to create dashboard:'));
        console.error(chalk.red(error.message));
        return null;
    }
}

// Run the script
if (require.main === module) {
    createSimpleDashboard();
}

module.exports = { createSimpleDashboard };
