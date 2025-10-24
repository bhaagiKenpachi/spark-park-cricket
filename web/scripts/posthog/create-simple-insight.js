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

async function createSimpleInsight() {
    console.log(chalk.blue.bold('📊 Creating Simple Test Insight\n'));

    const client = new PostHogAPIClient();

    try {
        // Create a very simple insight
        const insightConfig = {
            name: 'Test Insight - Page Views',
            description: 'Simple test insight to verify PostHog access',
            filters: {
                events: [{ id: '$pageview' }], // Use built-in pageview event
                date_from: '-7d', // Last 7 days
                interval: 'day',
                insight: 'TRENDS',
                display: 'ActionsLineGraph'
            }
        };

        console.log(chalk.blue('📈 Creating test insight...'));
        const insight = await client.createInsight(insightConfig);

        console.log(chalk.green(`✅ Insight created successfully!`));
        console.log(chalk.cyan(`   Name: ${insight.name}`));
        console.log(chalk.cyan(`   ID: ${insight.id}`));
        console.log(chalk.cyan(`   URL: ${client.host}/project/${client.projectId}/insights/${insight.id}`));
        console.log('');

        console.log(chalk.yellow('💡 Next steps:'));
        console.log(chalk.yellow('1. Try accessing the insight URL above'));
        console.log(chalk.yellow('2. If it works, we can create more insights'));
        console.log(chalk.yellow('3. If it doesn\'t work, we know the issue is with web access'));

        return insight;

    } catch (error) {
        console.error(chalk.red('❌ Failed to create test insight:'));
        console.error(chalk.red(error.message));

        if (error.message.includes('authentication')) {
            console.error(chalk.yellow('\n💡 Authentication issue. Check your Personal API Key.'));
        } else if (error.message.includes('project')) {
            console.error(chalk.yellow('\n💡 Project access issue. Check your Project ID.'));
        } else {
            console.error(chalk.yellow('\n💡 Unknown error. Check your PostHog configuration.'));
        }

        return null;
    }
}

// Run the script
if (require.main === module) {
    createSimpleInsight();
}

module.exports = { createSimpleInsight };
