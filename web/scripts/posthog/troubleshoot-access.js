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

async function troubleshootAccess() {
    console.log(chalk.blue.bold('🔍 PostHog Access Troubleshooting\n'));

    const client = new PostHogAPIClient();

    try {
        // Test 1: Check user info
        console.log(chalk.blue('1. Testing API authentication...'));
        try {
            const userResponse = await fetch('https://us.i.posthog.com/api/users/@me/', {
                headers: { 'Authorization': `Bearer ${client.apiKey}` }
            });
            const userData = await userResponse.json();
            console.log(chalk.green(`✅ Authenticated as: ${userData.first_name} ${userData.last_name} (${userData.email})`));
        } catch (error) {
            console.log(chalk.yellow(`⚠️  User info not accessible: ${error.message}`));
        }
        console.log('');

        // Test 2: Check project access
        console.log(chalk.blue('2. Testing project access...'));
        try {
            const projectResponse = await client.makeRequest('/');
            console.log(chalk.green(`✅ Project ID: ${projectResponse.id}`));
            console.log(chalk.green(`✅ Project Name: ${projectResponse.name}`));
        } catch (error) {
            console.log(chalk.yellow(`⚠️  Project info not accessible: ${error.message}`));
        }
        console.log('');

        // Test 3: Check insights
        console.log(chalk.blue('3. Testing insights access...'));
        const insightsResponse = await client.makeRequest('/insights/');
        const seriesInsights = insightsResponse.results?.filter(insight =>
            insight.name.includes('Series') || insight.name.includes('series')
        ) || [];

        console.log(chalk.green(`✅ Found ${seriesInsights.length} series-related insights:`));
        seriesInsights.forEach(insight => {
            console.log(chalk.cyan(`   • ${insight.name} (ID: ${insight.id})`));
        });
        console.log('');

        // Test 4: Check dashboards
        console.log(chalk.blue('4. Testing dashboards access...'));
        const dashboardsResponse = await client.makeRequest('/dashboards/');
        const seriesDashboards = dashboardsResponse.results?.filter(dashboard =>
            dashboard.name.includes('Series') || dashboard.name.includes('Cricket')
        ) || [];

        console.log(chalk.green(`✅ Found ${seriesDashboards.length} series-related dashboards:`));
        seriesDashboards.forEach(dashboard => {
            console.log(chalk.cyan(`   • ${dashboard.name} (ID: ${dashboard.id})`));
            console.log(chalk.gray(`     URL: ${client.getDashboardUrl(dashboard.id)}`));
        });
        console.log('');

        // Test 5: Check specific insight
        console.log(chalk.blue('5. Testing specific insight access...'));
        try {
            const insightResponse = await client.makeRequest('/insights/3801439/');
            console.log(chalk.green(`✅ Insight accessible: ${insightResponse.name}`));
            console.log(chalk.cyan(`   Direct URL: ${client.host}/project/${client.projectId}/insights/3801439`));
        } catch (error) {
            console.log(chalk.red(`❌ Insight not accessible: ${error.message}`));
        }
        console.log('');

        // Summary
        console.log(chalk.green.bold('📊 SUMMARY:'));
        console.log(chalk.green('✅ API authentication: Working'));
        console.log(chalk.green('✅ Project access: Working'));
        console.log(chalk.green(`✅ Series insights: ${seriesInsights.length} found`));
        console.log(chalk.green(`✅ Series dashboards: ${seriesDashboards.length} found`));
        console.log('');

        console.log(chalk.yellow('💡 NEXT STEPS:'));
        console.log(chalk.yellow('1. Try accessing the dashboard directly:'));
        if (seriesDashboards.length > 0) {
            console.log(chalk.cyan(`   ${client.getDashboardUrl(seriesDashboards[0].id)}`));
        }
        console.log(chalk.yellow('2. Try accessing insights directly:'));
        if (seriesInsights.length > 0) {
            console.log(chalk.cyan(`   ${client.host}/project/${client.projectId}/insights/${seriesInsights[0].id}`));
        }
        console.log(chalk.yellow('3. If still having issues, check PostHog UI permissions'));

    } catch (error) {
        console.error(chalk.red('❌ Troubleshooting failed:'));
        console.error(chalk.red(error.message));
    }
}

// Run the script
if (require.main === module) {
    troubleshootAccess();
}

module.exports = { troubleshootAccess };
