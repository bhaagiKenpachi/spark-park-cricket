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

async function listInsights() {
    console.log(chalk.blue('📊 Fetching insights...'));

    const client = new PostHogAPIClient();

    try {
        const response = await client.makeRequest('/insights/');

        if (response.results && response.results.length > 0) {
            console.log(chalk.green('\n📈 Available Insights:'));
            response.results.forEach(insight => {
                console.log(chalk.cyan(`  • ${insight.name} (ID: ${insight.id})`));
                console.log(chalk.gray(`    ${insight.description || 'No description'}`));
                console.log(chalk.gray(`    Created: ${new Date(insight.created_at).toLocaleDateString()}`));
                console.log('');
            });
        } else {
            console.log(chalk.yellow('No insights found.'));
        }

        return response.results || [];
    } catch (error) {
        console.error(chalk.red('Failed to fetch insights:', error.message));
        return [];
    }
}

// Run the script
if (require.main === module) {
    listInsights();
}

module.exports = { listInsights };
