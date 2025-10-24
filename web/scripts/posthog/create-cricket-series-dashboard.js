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

async function createCricketSeriesDashboard() {
    console.log(chalk.blue.bold('🏏 Creating Cricket Series Analytics Dashboard\n'));

    const client = new PostHogAPIClient();

    try {
        // Step 1: Create the main dashboard
        console.log(chalk.blue('📊 Creating Cricket Series Dashboard...'));

        const dashboardConfig = {
            name: 'Cricket Series Analytics',
            description: 'Comprehensive analytics dashboard for cricket series management, tracking series views, match interactions, and user engagement',
            pinned: true,
            show_tile_badges: true,
            tags: ['cricket', 'series', 'analytics', 'tournament-management']
        };

        const dashboard = await client.createDashboard(dashboardConfig);
        console.log(chalk.green(`✅ Dashboard created: ${dashboard.name}`));
        console.log(chalk.cyan(`   ID: ${dashboard.id}`));
        console.log(chalk.cyan(`   URL: ${client.getDashboardUrl(dashboard.id)}`));
        console.log('');

        // Step 2: Create meaningful series insights
        console.log(chalk.blue('📈 Creating Series Analytics Insights...'));

        const seriesInsights = [
            {
                name: 'Series Page Views',
                description: 'Track how many users view series pages over time',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview',
                            name: 'Series Views',
                            math: 'total',
                            properties: [{
                                key: '$current_url',
                                operator: 'icontains',
                                value: '/series/',
                                type: 'event'
                            }]
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['cricket', 'series', 'page-views']
            },
            {
                name: 'Series Created Events',
                description: 'Track when new cricket series are created',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'series_created',
                            name: 'Series Created',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['cricket', 'series', 'creation']
            },
            {
                name: 'Match Interactions',
                description: 'Track user interactions with matches (views, scoring, etc.)',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'match_viewed',
                            name: 'Match Views',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['cricket', 'matches', 'interactions']
            },
            {
                name: 'User Engagement',
                description: 'Track active users engaging with cricket features',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview',
                            name: 'Active Users',
                            math: 'dau'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['cricket', 'users', 'engagement']
            },
            {
                name: 'Scorecard Interactions',
                description: 'Track how users interact with live scorecards',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'scorecard_viewed',
                            name: 'Scorecard Views',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['cricket', 'scorecard', 'live-scoring']
            },
            {
                name: 'Error Tracking',
                description: 'Monitor errors and issues in the cricket application',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'error_occurred',
                            name: 'Errors',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['cricket', 'errors', 'monitoring']
            }
        ];

        const createdInsights = [];
        let addedToDashboard = 0;

        for (const insightConfig of seriesInsights) {
            console.log(chalk.blue(`   Creating: ${insightConfig.name}`));

            try {
                const insight = await client.createInsight(insightConfig);
                createdInsights.push(insight);
                console.log(chalk.green(`   ✅ Created: ${insight.name} (ID: ${insight.id})`));

                // Try to add to dashboard
                try {
                    await client.addInsightToDashboard(dashboard.id, insight.id);
                    addedToDashboard++;
                    console.log(chalk.cyan(`   🔗 Added to dashboard`));
                } catch (addError) {
                    console.log(chalk.yellow(`   ⚠️  Could not add to dashboard automatically`));
                }

            } catch (error) {
                console.log(chalk.red(`   ❌ Failed to create: ${insightConfig.name}`));
                console.log(chalk.gray(`      Error: ${error.message}`));
            }
        }

        console.log(chalk.green.bold('\n🎉 Cricket Series Dashboard Created!'));
        console.log(chalk.cyan(`📊 Dashboard: ${client.getDashboardUrl(dashboard.id)}`));
        console.log(chalk.cyan(`📈 Created ${createdInsights.length} insights`));
        console.log(chalk.cyan(`🔗 Added ${addedToDashboard} insights to dashboard`));
        console.log('');

        console.log(chalk.yellow('💡 Next Steps:'));
        console.log(chalk.yellow('   1. Go to your dashboard URL above'));
        console.log(chalk.yellow('   2. Start using your cricket app to generate events'));
        console.log(chalk.yellow('   3. Watch the analytics populate in real-time'));
        console.log(chalk.yellow('   4. Add more insights as needed'));

        console.log(chalk.magenta('\n🏏 Your Cricket Analytics are Ready!'));
        console.log(chalk.magenta('   • Series tracking'));
        console.log(chalk.magenta('   • Match analytics'));
        console.log(chalk.magenta('   • User engagement'));
        console.log(chalk.magenta('   • Error monitoring'));

        return { dashboard, insights: createdInsights };

    } catch (error) {
        console.error(chalk.red('❌ Failed to create cricket series dashboard:'));
        console.error(chalk.red(error.message));
        return null;
    }
}

// Run the script
if (require.main === module) {
    createCricketSeriesDashboard();
}

module.exports = { createCricketSeriesDashboard };
