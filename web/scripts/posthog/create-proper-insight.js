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

async function createProperInsight() {
    console.log(chalk.blue.bold('📊 Creating Proper PostHog Insight\n'));

    const client = new PostHogAPIClient();

    try {
        // Create a proper insight using the new format
        const insightConfig = {
            name: 'Test Insight - Page Views (Proper Format)',
            description: 'Test insight using the correct PostHog format',
            query: {
                kind: 'InsightVizNode',
                source: {
                    kind: 'TrendsQuery',
                    series: [
                        {
                            kind: 'EventsNode',
                            name: '$pageview',
                            event: '$pageview',
                            custom_name: 'Page Views'
                        }
                    ],
                    interval: 'day',
                    dateRange: {
                        date_from: '-7d',
                        explicitDate: false
                    },
                    properties: [],
                    filterTestAccounts: false
                },
                full: true
            }
        };

        console.log(chalk.blue('📈 Creating proper insight...'));
        const insight = await client.createInsight(insightConfig);

        console.log(chalk.green(`✅ Insight created successfully!`));
        console.log(chalk.cyan(`   Name: ${insight.name}`));
        console.log(chalk.cyan(`   ID: ${insight.id}`));
        console.log(chalk.cyan(`   URL: ${client.host}/project/${client.projectId}/insights/${insight.id}`));
        console.log('');

        console.log(chalk.yellow('💡 Test this insight:'));
        console.log(chalk.cyan(`   ${client.host}/project/${client.projectId}/insights/${insight.id}`));

        return insight;

    } catch (error) {
        console.error(chalk.red('❌ Failed to create proper insight:'));
        console.error(chalk.red(error.message));

        // Try a simpler approach
        console.log(chalk.yellow('\n🔄 Trying simpler approach...'));

        try {
            const simpleConfig = {
                name: 'Simple Test Insight',
                description: 'Very simple test insight',
                filters: {
                    events: [{ id: '$pageview' }],
                    date_from: '-7d',
                    interval: 'day',
                    insight: 'TRENDS'
                }
            };

            const simpleInsight = await client.createInsight(simpleConfig);
            console.log(chalk.green(`✅ Simple insight created: ${simpleInsight.name} (ID: ${simpleInsight.id})`));
            console.log(chalk.cyan(`   URL: ${client.host}/project/${client.projectId}/insights/${simpleInsight.id}`));

            return simpleInsight;
        } catch (simpleError) {
            console.error(chalk.red('❌ Simple approach also failed:'));
            console.error(chalk.red(simpleError.message));
            return null;
        }
    }
}

// Run the script
if (require.main === module) {
    createProperInsight();
}

module.exports = { createProperInsight };
