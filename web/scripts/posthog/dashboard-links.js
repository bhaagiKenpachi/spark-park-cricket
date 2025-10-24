#!/usr/bin/env node

const chalk = require('chalk');

// Dashboard and insights information
const DASHBOARD_ID = 615755;
const PROJECT_ID = 238221;
const HOST = 'https://us.i.posthog.com';

const INSIGHTS = [
    { id: 3801439, name: 'Series Views Over Time' },
    { id: 3801440, name: 'Series Created' },
    { id: 3801441, name: 'Series by Name' },
    { id: 3801442, name: 'Series Edits' },
    { id: 3801443, name: 'Users Creating Series' },
    { id: 3801444, name: 'Series Engagement' }
];

function showDashboardLinks() {
    console.log(chalk.blue.bold('📊 PostHog Dashboard & Insights Links\n'));

    console.log(chalk.green('🎯 MAIN DASHBOARD:'));
    console.log(chalk.cyan(`   ${HOST}/project/${PROJECT_ID}/dashboard/${DASHBOARD_ID}`));
    console.log('');

    console.log(chalk.green('📈 INDIVIDUAL INSIGHTS:'));
    INSIGHTS.forEach((insight, index) => {
        console.log(chalk.cyan(`   ${index + 1}. ${insight.name}`));
        console.log(chalk.gray(`      ID: ${insight.id}`));
        console.log(chalk.gray(`      URL: ${HOST}/project/${PROJECT_ID}/insights/${insight.id}`));
        console.log('');
    });

    console.log(chalk.yellow('💡 MANUAL STEPS:'));
    console.log(chalk.yellow('1. Go to the main dashboard URL above'));
    console.log(chalk.yellow('2. Click "Add insight" or "Edit dashboard"'));
    console.log(chalk.yellow('3. Search for each insight by name or ID'));
    console.log(chalk.yellow('4. Add them to the dashboard layout'));
    console.log('');

    console.log(chalk.green('🎉 Once added, your dashboard will show:'));
    console.log(chalk.gray('   • Real-time series analytics'));
    console.log(chalk.gray('   • Auto-updating charts'));
    console.log(chalk.gray('   • Comprehensive insights'));
    console.log(chalk.gray('   • Grafana-like live data'));
}

// Run the script
if (require.main === module) {
    showDashboardLinks();
}

module.exports = { showDashboardLinks };
