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

async function createProperSavedInsight() {
    console.log(chalk.blue.bold('📊 Creating Proper Saved Insight\n'));

    const client = new PostHogAPIClient();

    try {
        // Create insight with the correct structure for PostHog API
        const insightConfig = {
            name: 'Series Page Views - Properly Saved',
            description: 'Track series page views with correct API structure',
            query: {
                kind: 'InsightVizNode',
                source: {
                    kind: 'TrendsQuery',
                    series: [{
                        kind: 'EventsNode',
                        event: '$pageview',
                        name: '$pageview',
                        math: 'total'
                    }],
                    dateRange: {
                        date_from: '-7d'
                    },
                    interval: 'day'
                }
            },
            tags: ['cricket', 'series', 'analytics', 'saved']
        };

        console.log(chalk.blue('📈 Creating properly structured insight...'));
        const insight = await client.createInsight(insightConfig);

        console.log(chalk.green(`✅ Insight created and saved successfully!`));
        console.log(chalk.cyan(`   Name: ${insight.name}`));
        console.log(chalk.cyan(`   ID: ${insight.id}`));
        console.log(chalk.cyan(`   URL: ${client.host}/project/${client.projectId}/insights/${insight.id}`));
        console.log('');

        // Test if we can retrieve the insight
        console.log(chalk.blue('🔍 Testing if insight is properly saved...'));

        try {
            const fetchedInsight = await client.makeRequest(`/insights/${insight.id}/`);
            console.log(chalk.green('✅ Insight is properly saved and retrievable'));
            console.log(chalk.cyan(`   Retrieved: ${fetchedInsight.name}`));
        } catch (fetchError) {
            console.log(chalk.yellow('⚠️  Insight created but may not be fully saved'));
            console.log(chalk.gray(`   Error: ${fetchError.message}`));
        }

        console.log(chalk.green.bold('\n🎉 Proper Insight Created!'));
        console.log(chalk.yellow('💡 This insight should now appear in your saved insights list'));
        console.log(chalk.cyan(`📈 Insight URL: ${client.host}/project/${client.projectId}/insights/${insight.id}`));

        return insight;

    } catch (error) {
        console.error(chalk.red('❌ Failed to create proper saved insight:'));
        console.error(chalk.red(error.message));

        // Try the simplest possible structure
        console.log(chalk.yellow('\n🔄 Trying minimal structure...'));

        try {
            const minimalConfig = {
                name: 'Minimal Test Insight',
                description: 'Minimal insight structure test',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview'
                        }]
                    }
                }
            };

            const minimalInsight = await client.createInsight(minimalConfig);
            console.log(chalk.green(`✅ Minimal insight created: ${minimalInsight.name} (ID: ${minimalInsight.id})`));
            console.log(chalk.cyan(`   URL: ${client.host}/project/${client.projectId}/insights/${minimalInsight.id}`));

            return minimalInsight;
        } catch (minimalError) {
            console.error(chalk.red('❌ Minimal approach also failed:'));
            console.error(chalk.red(minimalError.message));
            return null;
        }
    }
}

// Run the script
if (require.main === module) {
    createProperSavedInsight();
}

module.exports = { createProperSavedInsight };
