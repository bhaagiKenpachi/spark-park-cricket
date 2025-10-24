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

async function createProperAPIInsight() {
    console.log(chalk.blue.bold('📊 Creating Proper PostHog API Insight (Correct Format)\n'));

    const client = new PostHogAPIClient();

    try {
        // Create insight using the exact PostHog API format from documentation
        const insightConfig = {
            name: 'Series Page Views - API Format',
            description: 'Track series page views using proper PostHog API format',
            query: {
                kind: 'TrendsQuery',
                series: [{
                    kind: 'EventsNode',
                    event: '$pageview',
                    name: '$pageview',
                    properties: [
                        {
                            key: '$current_url',
                            operator: 'icontains',
                            value: '/series/',
                            type: 'string'
                        }
                    ]
                }],
                dateRange: {
                    date_from: '-7d'
                },
                interval: 'day',
                breakdownFilter: {
                    breakdown_type: 'event',
                    breakdown: '$current_url'
                }
            },
            tags: ['cricket', 'series', 'analytics']
        };

        console.log(chalk.blue('📈 Creating insight with proper API format...'));
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
        console.error(chalk.red('❌ Failed to create proper API insight:'));
        console.error(chalk.red(error.message));

        // Try the simplest possible approach
        console.log(chalk.yellow('\n🔄 Trying simplest approach...'));

        try {
            const simpleConfig = {
                name: 'Basic Page Views - Simple',
                description: 'Basic page view tracking',
                query: {
                    kind: 'TrendsQuery',
                    series: [{
                        kind: 'EventsNode',
                        event: '$pageview'
                    }],
                    dateRange: {
                        date_from: '-7d'
                    }
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
    createProperAPIInsight();
}

module.exports = { createProperAPIInsight };
