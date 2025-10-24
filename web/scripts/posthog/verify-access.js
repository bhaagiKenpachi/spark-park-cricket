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

async function verifyAccess() {
    console.log(chalk.blue.bold('🔍 Verifying PostHog Access\n'));

    const projectId = process.env.POSTHOG_PROJECT_ID;
    const apiKey = process.env.POSTHOG_PERSONAL_API_KEY;
    const host = process.env.POSTHOG_HOST || 'https://us.i.posthog.com';

    console.log(chalk.blue('📊 Configuration:'));
    console.log(chalk.gray(`   Project ID: ${projectId}`));
    console.log(chalk.gray(`   Host: ${host}`));
    console.log(chalk.gray(`   API Key: ${apiKey ? 'Set' : 'Not set'}`));
    console.log('');

    if (!projectId || !apiKey) {
        console.log(chalk.red('❌ Missing configuration. Please check your .env.local file.'));
        return;
    }

    try {
        // Test 1: Check if we can access the project
        console.log(chalk.blue('1. Testing project access...'));
        const projectResponse = await fetch(`${host}/api/projects/${projectId}/`, {
            headers: { 'Authorization': `Bearer ${apiKey}` }
        });

        if (projectResponse.ok) {
            const projectData = await projectResponse.json();
            console.log(chalk.green(`✅ Project accessible: ${projectData.name}`));
        } else {
            console.log(chalk.red(`❌ Project not accessible: ${projectResponse.status}`));
            console.log(chalk.yellow('💡 This might be why you\'re getting 404 errors'));
        }
        console.log('');

        // Test 2: Check insights
        console.log(chalk.blue('2. Testing insights access...'));
        const insightsResponse = await fetch(`${host}/api/projects/${projectId}/insights/`, {
            headers: { 'Authorization': `Bearer ${apiKey}` }
        });

        if (insightsResponse.ok) {
            const insightsData = await insightsResponse.json();
            const seriesInsights = insightsData.results?.filter(insight =>
                insight.name.includes('Series') || insight.name.includes('series')
            ) || [];
            console.log(chalk.green(`✅ Insights accessible: ${seriesInsights.length} series insights found`));
        } else {
            console.log(chalk.red(`❌ Insights not accessible: ${insightsResponse.status}`));
        }
        console.log('');

        // Test 3: Check dashboards
        console.log(chalk.blue('3. Testing dashboards access...'));
        const dashboardsResponse = await fetch(`${host}/api/projects/${projectId}/dashboards/`, {
            headers: { 'Authorization': `Bearer ${apiKey}` }
        });

        if (dashboardsResponse.ok) {
            const dashboardsData = await dashboardsResponse.json();
            const seriesDashboards = dashboardsData.results?.filter(dashboard =>
                dashboard.name.includes('Series') || dashboard.name.includes('Cricket')
            ) || [];
            console.log(chalk.green(`✅ Dashboards accessible: ${seriesDashboards.length} series dashboards found`));
        } else {
            console.log(chalk.red(`❌ Dashboards not accessible: ${dashboardsResponse.status}`));
        }
        console.log('');

        // Generate correct URLs
        console.log(chalk.green.bold('🔗 Correct URLs to try:'));
        console.log(chalk.cyan(`   Main Dashboard: ${host}/project/${projectId}/dashboard/615755`));
        console.log(chalk.cyan(`   All Insights: ${host}/project/${projectId}/insights`));
        console.log(chalk.cyan(`   All Dashboards: ${host}/project/${projectId}/dashboards`));
        console.log('');

        console.log(chalk.yellow('💡 If you\'re still getting 404 errors:'));
        console.log(chalk.yellow('1. Make sure you\'re logged into the correct PostHog account'));
        console.log(chalk.yellow('2. Check if you have access to project 238221'));
        console.log(chalk.yellow('3. Try logging out and back into PostHog'));
        console.log(chalk.yellow('4. Check if the project ID is correct'));

    } catch (error) {
        console.error(chalk.red('❌ Verification failed:'));
        console.error(chalk.red(error.message));
    }
}

// Run the script
if (require.main === module) {
    verifyAccess();
}

module.exports = { verifyAccess };
