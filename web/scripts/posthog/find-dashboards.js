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

async function findDashboards() {
    console.log(chalk.blue('📋 Fetching all dashboards...'));

    const client = new PostHogAPIClient();

    try {
        const response = await client.makeRequest('/dashboards/');

        if (response.results && response.results.length > 0) {
            console.log(chalk.green('\n📊 Available Dashboards:'));
            response.results.forEach((dashboard, index) => {
                console.log(chalk.cyan(`  ${index + 1}. ${dashboard.name} (ID: ${dashboard.id})`));
                console.log(chalk.gray(`     Description: ${dashboard.description || 'No description'}`));
                console.log(chalk.gray(`     Created: ${new Date(dashboard.created_at).toLocaleDateString()}`));
                console.log(chalk.gray(`     URL: ${client.getDashboardUrl(dashboard.id)}`));
                console.log(chalk.gray(`     Tiles: ${dashboard.tiles?.length || 0}`));
                console.log('');
            });
        } else {
            console.log(chalk.yellow('No dashboards found.'));
        }

        return response.results || [];
    } catch (error) {
        console.error(chalk.red('Failed to fetch dashboards:', error.message));
        return [];
    }
}

// Run the script
if (require.main === module) {
    findDashboards();
}

module.exports = { findDashboards };
