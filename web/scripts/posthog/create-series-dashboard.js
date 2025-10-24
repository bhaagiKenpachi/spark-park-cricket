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

// Series Dashboard Insights Configuration
const SERIES_INSIGHTS = [
    {
        name: 'Series Views Over Time',
        description: 'Track series views over the last 30 days',
        filters: {
            events: [{ id: 'series_viewed' }],
            date_from: '-30d',
            interval: 'day',
            insight: 'TRENDS',
            display: 'ActionsLineGraph'
        }
    },
    {
        name: 'Series Created',
        description: 'Track when new series are created',
        filters: {
            events: [{ id: 'series_created' }],
            date_from: '-30d',
            interval: 'day',
            insight: 'TRENDS',
            display: 'ActionsBar'
        }
    },
    {
        name: 'Series by Name',
        description: 'Most viewed series by name',
        filters: {
            events: [{ id: 'series_viewed' }],
            date_from: '-30d',
            breakdown: 'properties.series_name',
            insight: 'TRENDS',
            display: 'ActionsTable'
        }
    },
    {
        name: 'Series Edits',
        description: 'Track series editing activity',
        filters: {
            events: [{ id: 'series_edited' }],
            date_from: '-30d',
            interval: 'day',
            insight: 'TRENDS',
            display: 'ActionsLineGraph'
        }
    },
    {
        name: 'Users Creating Series',
        description: 'Unique users who created series',
        filters: {
            events: [{ id: 'series_created' }],
            date_from: '-30d',
            insight: 'TRENDS',
            display: 'ActionsLineGraph',
            math: 'dau' // Daily Active Users
        }
    },
    {
        name: 'Series Engagement',
        description: 'Combined series activity (views, creates, edits)',
        filters: {
            events: [
                { id: 'series_viewed', name: 'Series Views' },
                { id: 'series_created', name: 'Series Created' },
                { id: 'series_edited', name: 'Series Edited' }
            ],
            date_from: '-30d',
            interval: 'day',
            insight: 'TRENDS',
            display: 'ActionsLineGraph'
        }
    }
];

async function createSeriesDashboard() {
    console.log(chalk.blue.bold('🏏 Creating Series Analytics Dashboard for PostHog\n'));

    const client = new PostHogAPIClient();

    try {
        // Create dashboard
        const dashboardConfig = {
            name: 'Series Analytics',
            description: 'Comprehensive analytics dashboard for cricket series management',
            pinned: true
        };

        const dashboard = await client.createDashboard(dashboardConfig);
        console.log(chalk.green(`\n📋 Dashboard created successfully!`));
        console.log(chalk.cyan(`🔗 View dashboard: ${client.getDashboardUrl(dashboard.id)}\n`));

        // Create insights and add to dashboard
        const insightIds = [];

        for (const insightConfig of SERIES_INSIGHTS) {
            const insight = await client.createInsight(insightConfig);
            insightIds.push(insight.id);

            // Add insight to dashboard
            await client.addInsightToDashboard(dashboard.id, insight.id);
        }

        console.log(chalk.green.bold('\n🎉 Series Analytics Dashboard created successfully!'));
        console.log(chalk.cyan(`📊 Dashboard URL: ${client.getDashboardUrl(dashboard.id)}`));
        console.log(chalk.gray(`\nThe dashboard includes ${SERIES_INSIGHTS.length} insights:`));

        SERIES_INSIGHTS.forEach((insight, index) => {
            console.log(chalk.gray(`  ${index + 1}. ${insight.name}`));
        });

        console.log(chalk.yellow('\n💡 The dashboard will automatically update as new series events are captured.'));
        console.log(chalk.yellow('   No manual refresh needed - similar to Grafana live queries!'));

    } catch (error) {
        console.error(chalk.red.bold('\n❌ Failed to create Series Analytics Dashboard:'));
        console.error(chalk.red(error.message));

        if (error.message.includes('authentication')) {
            console.error(chalk.yellow('\n💡 Make sure you have:'));
            console.error(chalk.yellow('   1. Generated a Personal API Key from PostHog Settings'));
            console.error(chalk.yellow('   2. Set POSTHOG_PERSONAL_API_KEY in your .env file'));
            console.error(chalk.yellow('   3. Set POSTHOG_PROJECT_ID in your .env file'));
        }

        process.exit(1);
    }
}

// Run the script
if (require.main === module) {
    createSeriesDashboard();
}

module.exports = { createSeriesDashboard, SERIES_INSIGHTS };
