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

async function createWebAnalyticsDashboard() {
    console.log(chalk.blue.bold('🌐 Creating Web Analytics Dashboard\n'));

    const client = new PostHogAPIClient();

    try {
        // Step 1: Create the web analytics dashboard
        console.log(chalk.blue('📊 Creating Web Analytics Dashboard...'));

        const dashboardConfig = {
            name: 'Web Analytics Dashboard',
            description: 'Comprehensive web analytics for application performance, user behavior, and technical metrics',
            pinned: true,
            show_tile_badges: true,
            tags: ['web', 'analytics', 'performance', 'user-behavior', 'technical']
        };

        const dashboard = await client.createDashboard(dashboardConfig);
        console.log(chalk.green(`✅ Dashboard created: ${dashboard.name}`));
        console.log(chalk.cyan(`   ID: ${dashboard.id}`));
        console.log(chalk.cyan(`   URL: ${client.getDashboardUrl(dashboard.id)}`));
        console.log('');

        // Step 2: Create comprehensive web analytics insights
        console.log(chalk.blue('📈 Creating Web Analytics Insights...'));

        const webAnalyticsInsights = [
            {
                name: 'Page Views Over Time',
                description: 'Track overall page views and traffic patterns',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview',
                            name: 'Page Views',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['web', 'traffic', 'page-views']
            },
            {
                name: 'Daily Active Users',
                description: 'Track daily active users (DAU) for user engagement',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview',
                            name: 'Daily Active Users',
                            math: 'dau'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['web', 'users', 'engagement', 'dau']
            },
            {
                name: 'Page Load Performance',
                description: 'Monitor page load times and performance metrics',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'performance_metric',
                            name: 'Page Load Time',
                            math: 'avg',
                            properties: [{
                                key: 'metric_name',
                                operator: 'exact',
                                value: 'page_load_time',
                                type: 'event'
                            }]
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['web', 'performance', 'page-speed']
            },
            {
                name: 'API Response Times',
                description: 'Track API performance and response times',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'performance_metric',
                            name: 'API Response Time',
                            math: 'avg',
                            properties: [{
                                key: 'metric_name',
                                operator: 'exact',
                                value: 'api_response_time',
                                type: 'event'
                            }]
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['web', 'api', 'performance', 'backend']
            },
            {
                name: 'Error Rate Tracking',
                description: 'Monitor application errors and error rates',
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
                tags: ['web', 'errors', 'monitoring', 'reliability']
            },
            {
                name: 'User Session Duration',
                description: 'Track how long users stay on the application',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'time_on_page',
                            name: 'Session Duration',
                            math: 'avg',
                            properties: [{
                                key: 'duration',
                                operator: 'gt',
                                value: 0,
                                type: 'event'
                            }]
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['web', 'engagement', 'sessions', 'user-behavior']
            },
            {
                name: 'Feature Usage',
                description: 'Track which features users interact with most',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: 'feature_used',
                            name: 'Feature Usage',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day'
                    }
                },
                tags: ['web', 'features', 'usage', 'user-behavior']
            },
            {
                name: 'Browser & Device Analytics',
                description: 'Track user devices and browsers for technical insights',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview',
                            name: 'Page Views',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day',
                        breakdownFilter: {
                            breakdown_type: 'event',
                            breakdown: '$browser'
                        }
                    }
                },
                tags: ['web', 'devices', 'browsers', 'technical']
            },
            {
                name: 'User Journey Funnel',
                description: 'Track user journey through key application steps',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'FunnelsQuery',
                        series: [
                            {
                                kind: 'EventsNode',
                                event: '$pageview',
                                name: 'Landing Page'
                            },
                            {
                                kind: 'EventsNode',
                                event: 'user_logged_in',
                                name: 'User Login'
                            },
                            {
                                kind: 'EventsNode',
                                event: 'series_viewed',
                                name: 'Series Interaction'
                            },
                            {
                                kind: 'EventsNode',
                                event: 'match_viewed',
                                name: 'Match Interaction'
                            }
                        ],
                        dateRange: {
                            date_from: '-30d'
                        },
                        funnelsFilter: {
                            layout: 'horizontal',
                            funnelVizType: 'steps',
                            funnelOrderType: 'ordered'
                        }
                    }
                },
                tags: ['web', 'funnel', 'conversion', 'user-journey']
            },
            {
                name: 'Geographic Distribution',
                description: 'Track user locations and geographic distribution',
                query: {
                    kind: 'InsightVizNode',
                    source: {
                        kind: 'TrendsQuery',
                        series: [{
                            kind: 'EventsNode',
                            event: '$pageview',
                            name: 'Page Views',
                            math: 'total'
                        }],
                        dateRange: {
                            date_from: '-30d'
                        },
                        interval: 'day',
                        breakdownFilter: {
                            breakdown_type: 'person',
                            breakdown: '$geoip_country_code'
                        }
                    }
                },
                tags: ['web', 'geography', 'location', 'global']
            }
        ];

        const createdInsights = [];
        let addedToDashboard = 0;

        for (const insightConfig of webAnalyticsInsights) {
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

        console.log(chalk.green.bold('\n🎉 Web Analytics Dashboard Created!'));
        console.log(chalk.cyan(`📊 Dashboard: ${client.getDashboardUrl(dashboard.id)}`));
        console.log(chalk.cyan(`📈 Created ${createdInsights.length} insights`));
        console.log(chalk.cyan(`🔗 Added ${addedToDashboard} insights to dashboard`));
        console.log('');

        console.log(chalk.yellow('💡 Web Analytics Coverage:'));
        console.log(chalk.yellow('   • Traffic & Page Views'));
        console.log(chalk.yellow('   • User Engagement (DAU, Sessions)'));
        console.log(chalk.yellow('   • Performance Metrics (Load Times, API)'));
        console.log(chalk.yellow('   • Error Monitoring'));
        console.log(chalk.yellow('   • Feature Usage'));
        console.log(chalk.yellow('   • User Journey & Funnels'));
        console.log(chalk.yellow('   • Geographic & Device Analytics'));

        console.log(chalk.magenta('\n🌐 Your Web Analytics are Ready!'));
        console.log(chalk.magenta('   • Complete web performance monitoring'));
        console.log(chalk.magenta('   • User behavior tracking'));
        console.log(chalk.magenta('   • Technical metrics and reliability'));
        console.log(chalk.magenta('   • Business intelligence insights'));

        return { dashboard, insights: createdInsights };

    } catch (error) {
        console.error(chalk.red('❌ Failed to create web analytics dashboard:'));
        console.error(chalk.red(error.message));
        return null;
    }
}

// Run the script
if (require.main === module) {
    createWebAnalyticsDashboard();
}

module.exports = { createWebAnalyticsDashboard };
