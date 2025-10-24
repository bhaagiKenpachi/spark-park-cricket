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

async function testAPIConnection() {
    console.log(chalk.blue.bold('🔍 Testing PostHog API Connection\n'));

    const client = new PostHogAPIClient();

    try {
        // Test 1: Check if we can make a basic request
        console.log(chalk.blue('1. Testing basic API connection...'));

        // Try to get project info (this should work)
        const projectInfo = await client.makeRequest('/');
        console.log(chalk.green('✅ Basic API connection works'));
        console.log(chalk.cyan(`   Project: ${projectInfo.name || 'Unknown'}`));
        console.log(chalk.cyan(`   ID: ${projectInfo.id}`));
        console.log('');

        // Test 2: Try to create a very simple insight
        console.log(chalk.blue('2. Testing insight creation...'));

        const simpleInsight = {
            name: 'Test Insight - API Connection',
            description: 'Testing if we can create insights',
            filters: {
                events: [{ id: '$pageview' }],
                date_from: '-1d',
                insight: 'TRENDS'
            }
        };

        const insight = await client.createInsight(simpleInsight);
        console.log(chalk.green('✅ Insight creation works'));
        console.log(chalk.cyan(`   Insight ID: ${insight.id}`));
        console.log(chalk.cyan(`   URL: ${client.host}/project/${client.projectId}/insights/${insight.id}`));
        console.log('');

        // Test 3: Try to fetch the insight we just created
        console.log(chalk.blue('3. Testing insight retrieval...'));

        try {
            const fetchedInsight = await client.makeRequest(`/insights/${insight.id}/`);
            console.log(chalk.green('✅ Insight retrieval works'));
            console.log(chalk.cyan(`   Fetched: ${fetchedInsight.name}`));
        } catch (fetchError) {
            console.log(chalk.yellow('⚠️  Insight retrieval failed, but creation worked'));
            console.log(chalk.gray(`   Error: ${fetchError.message}`));
        }

        console.log(chalk.green.bold('\n🎉 API Connection Test Complete!'));
        console.log(chalk.yellow('💡 Try accessing the insight URL above in your browser'));

        return insight;

    } catch (error) {
        console.error(chalk.red('❌ API Connection Test Failed:'));
        console.error(chalk.red(error.message));

        // Try to get more details about the error
        if (error.message.includes('401')) {
            console.log(chalk.yellow('\n💡 Authentication Error:'));
            console.log(chalk.yellow('   - Check your Personal API Key'));
            console.log(chalk.yellow('   - Make sure it starts with "phx_"'));
            console.log(chalk.yellow('   - Verify it\'s not expired'));
        } else if (error.message.includes('404')) {
            console.log(chalk.yellow('\n💡 Project Not Found:'));
            console.log(chalk.yellow('   - Check your Project ID'));
            console.log(chalk.yellow('   - Verify you have access to this project'));
        } else if (error.message.includes('500')) {
            console.log(chalk.yellow('\n💡 Server Error:'));
            console.log(chalk.yellow('   - PostHog servers might be experiencing issues'));
            console.log(chalk.yellow('   - Try again in a few minutes'));
            console.log(chalk.yellow('   - Check PostHog status page'));
        }

        return null;
    }
}

// Run the script
if (require.main === module) {
    testAPIConnection();
}

module.exports = { testAPIConnection };
