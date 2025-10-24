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

async function createCorrectInsight() {
    console.log(chalk.blue.bold('📊 Creating Correct PostHog Insight (API Format)\n'));

    const client = new PostHogAPIClient();

    try {
        // Create insight using the correct PostHog API format
        const insightConfig = {
            name: 'Series Views - Correct Format',
            description: 'Track series page views using proper PostHog API format',
            filters: {
                events: [
                    {
                        id: '$pageview',
                        name: '$pageview',
                        type: 'events',
                        order: 0,
                        properties: [
                            {
                                key: '$current_url',
                                operator: 'icontains',
                                value: '/series/',
                                type: 'string'
                            }
                        ]
                    }
                ],
                date_from: '-7d',
                date_to: null,
                interval: 'day',
                insight: 'TRENDS',
                display: 'ActionsLineGraph',
                breakdown_type: 'event',
                breakdown: '$current_url'
            },
            tags: ['cricket', 'series', 'analytics']
        };

        console.log(chalk.blue('📈 Creating insight with correct API format...'));
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
        console.error(chalk.red('❌ Failed to create correct insight:'));
        console.error(chalk.red(error.message));

        // Try the simplest possible approach
        console.log(chalk.yellow('\n🔄 Trying simplest approach...'));

        try {
            const simpleConfig = {
                name: 'Simple Page Views',
                description: 'Basic page view tracking',
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
    createCorrectInsight();
}

module.exports = { createCorrectInsight };
