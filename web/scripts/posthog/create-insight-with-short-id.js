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

async function createInsightWithShortId() {
    console.log(chalk.blue.bold('📊 Creating Insight with Proper Short ID Handling\n'));

    const client = new PostHogAPIClient();

    try {
        // Create insight with the correct structure
        const insightConfig = {
            name: 'Cricket Series Analytics - Short ID Test',
            description: 'Test insight with proper short ID handling',
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
            tags: ['cricket', 'series', 'analytics', 'short-id-test']
        };

        console.log(chalk.blue('📈 Creating insight with proper short ID handling...'));
        const insight = await client.createInsight(insightConfig);

        console.log(chalk.green.bold('\n🎉 Insight Created Successfully!'));
        console.log(chalk.cyan(`   Name: ${insight.name}`));
        console.log(chalk.cyan(`   Database ID: ${insight.id}`));
        console.log(chalk.cyan(`   Short ID: ${insight.short_id}`));
        console.log(chalk.cyan(`   Correct URL: ${client.getInsightUrl(insight.short_id)}`));
        console.log('');

        console.log(chalk.yellow('💡 Key Points:'));
        console.log(chalk.yellow('   • Use short_id for URLs, not the numeric ID'));
        console.log(chalk.yellow('   • Short ID is what appears in the PostHog UI'));
        console.log(chalk.yellow('   • This insight should now be visible in your saved insights'));
        console.log('');

        console.log(chalk.green('🔗 Test Links:'));
        console.log(chalk.cyan(`   Direct insight: ${client.getInsightUrl(insight.short_id)}`));
        console.log(chalk.cyan(`   Dashboard: ${client.getDashboardUrl(615947)}`));
        console.log('');

        console.log(chalk.magenta('📋 Next Steps:'));
        console.log(chalk.magenta('   1. Go to the insight URL above'));
        console.log(chalk.magenta('   2. Click "Add to dashboard"'));
        console.log(chalk.magenta('   3. Select "Cricket Analytics Dashboard"'));
        console.log(chalk.magenta('   4. The insight will appear on your dashboard'));

        return insight;

    } catch (error) {
        console.error(chalk.red('❌ Failed to create insight with short ID:'));
        console.error(chalk.red(error.message));
        return null;
    }
}

// Run the script
if (require.main === module) {
    createInsightWithShortId();
}

module.exports = { createInsightWithShortId };
