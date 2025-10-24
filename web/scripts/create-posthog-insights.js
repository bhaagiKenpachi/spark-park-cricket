#!/usr/bin/env node

/**
 * Script to create PostHog insights and dashboards programmatically
 * Run this after verifying events are being captured
 */

const POSTHOG_API_KEY = 'phc_91Oi5cxrr7HPQXzP8hrYJ8oOxvdOQ5XmmGxAsbQ7OaV';
const POSTHOG_HOST = 'https://us.i.posthog.com';

// Function to create an insight
async function createInsight(name, description, query) {
    const response = await fetch(`${POSTHOG_HOST}/api/projects/1/insights/`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${POSTHOG_API_KEY}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            name,
            description,
            query: {
                kind: 'DataVisualizationNode',
                source: {
                    kind: 'EventsQuery',
                    select: query.select,
                    where: query.where,
                    orderBy: query.orderBy || ['timestamp'],
                    limit: query.limit || 100
                }
            }
        })
    });

    if (!response.ok) {
        console.error(`Failed to create insight "${name}":`, await response.text());
        return null;
    }

    const result = await response.json();
    console.log(`✅ Created insight: ${name} (ID: ${result.id})`);
    return result;
}

// Function to create a dashboard
async function createDashboard(name, description, insightIds) {
    const response = await fetch(`${POSTHOG_HOST}/api/projects/1/dashboards/`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${POSTHOG_API_KEY}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            name,
            description,
            tiles: insightIds.map((id, index) => ({
                insight: id,
                layouts: {
                    sm: { x: (index % 2) * 6, y: Math.floor(index / 2) * 4, w: 6, h: 4 },
                    lg: { x: (index % 3) * 4, y: Math.floor(index / 3) * 4, w: 4, h: 4 }
                }
            }))
        })
    });

    if (!response.ok) {
        console.error(`Failed to create dashboard "${name}":`, await response.text());
        return null;
    }

    const result = await response.json();
    console.log(`✅ Created dashboard: ${name} (ID: ${result.id})`);
    return result;
}

async function main() {
    console.log('🚀 Creating PostHog insights and dashboards...\n');

    try {
        // Create insights for different analytics
        const insights = [];

        // 1. Series Analytics
        const seriesViewed = await createInsight(
            'Series Views Over Time',
            'Track how many times series are viewed',
            {
                select: ['timestamp', 'properties.series_id', 'properties.series_name'],
                where: [['event', '=', 'series_viewed']]
            }
        );
        if (seriesViewed) insights.push(seriesViewed.id);

        const seriesCreated = await createInsight(
            'Series Creation Events',
            'Track when new series are created',
            {
                select: ['timestamp', 'properties.series_id', 'properties.series_name', 'properties.user_id'],
                where: [['event', '=', 'series_created']]
            }
        );
        if (seriesCreated) insights.push(seriesCreated.id);

        // 2. Match Analytics
        const matchViews = await createInsight(
            'Match Views',
            'Track match viewing patterns',
            {
                select: ['timestamp', 'properties.match_id', 'properties.series_id'],
                where: [['event', '=', 'match_viewed']]
            }
        );
        if (matchViews) insights.push(matchViews.id);

        const matchCreated = await createInsight(
            'Match Creation',
            'Track new match creation',
            {
                select: ['timestamp', 'properties.match_id', 'properties.series_id', 'properties.total_overs'],
                where: [['event', '=', 'match_created']]
            }
        );
        if (matchCreated) insights.push(matchCreated.id);

        // 3. Scoring Analytics
        const ballAdded = await createInsight(
            'Ball Scoring Activity',
            'Track ball-by-ball scoring',
            {
                select: ['timestamp', 'properties.match_id', 'properties.ball_type', 'properties.runs'],
                where: [['event', '=', 'ball_added']]
            }
        );
        if (ballAdded) insights.push(ballAdded.id);

        const overCompleted = await createInsight(
            'Over Completion',
            'Track when overs are completed',
            {
                select: ['timestamp', 'properties.match_id', 'properties.over_number', 'properties.total_runs'],
                where: [['event', '=', 'over_completed']]
            }
        );
        if (overCompleted) insights.push(overCompleted.id);

        // 4. Performance Analytics
        const pageLoadTimes = await createInsight(
            'Page Load Performance',
            'Track page load times across the app',
            {
                select: ['timestamp', 'properties.metric_value', 'properties.page'],
                where: [['event', '=', 'performance_metric'], ['properties.metric_name', '=', 'page_load_time']]
            }
        );
        if (pageLoadTimes) insights.push(pageLoadTimes.id);

        const apiResponseTimes = await createInsight(
            'API Response Times',
            'Track API performance',
            {
                select: ['timestamp', 'properties.metric_value', 'properties.component'],
                where: [['event', '=', 'performance_metric'], ['properties.metric_name', '=', 'api_response_time']]
            }
        );
        if (apiResponseTimes) insights.push(apiResponseTimes.id);

        // 5. Error Tracking
        const errors = await createInsight(
            'Application Errors',
            'Track errors and exceptions',
            {
                select: ['timestamp', 'properties.error_type', 'properties.error_message', 'properties.component'],
                where: [['event', '=', 'error_occurred']]
            }
        );
        if (errors) insights.push(errors.id);

        // Create main dashboard
        if (insights.length > 0) {
            await createDashboard(
                'Cricket Analytics Dashboard',
                'Comprehensive analytics dashboard for cricket tournament management',
                insights
            );
        }

        console.log('\n🎉 All insights and dashboards created successfully!');
        console.log('\n📊 Next steps:');
        console.log('1. Go to your PostHog dashboard');
        console.log('2. Navigate to "Dashboards" to see your new dashboard');
        console.log('3. Check "Insights" to see individual analytics');
        console.log('4. Use "Events" to see raw event data');

    } catch (error) {
        console.error('❌ Error creating insights:', error);
    }
}

// Run the script
main();
