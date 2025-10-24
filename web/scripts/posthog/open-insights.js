#!/usr/bin/env node

const chalk = require('chalk');

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

function showNavigationLinks() {
    console.log(chalk.blue.bold('🔗 PostHog Navigation Links\n'));

    console.log(chalk.green('📊 MAIN SECTIONS:'));
    console.log(chalk.cyan(`   Dashboard: ${HOST}/project/${PROJECT_ID}/dashboard/615755`));
    console.log(chalk.cyan(`   All Insights: ${HOST}/project/${PROJECT_ID}/insights`));
    console.log('');

    console.log(chalk.green('📈 YOUR SERIES INSIGHTS:'));
    INSIGHTS.forEach((insight, index) => {
        console.log(chalk.cyan(`   ${index + 1}. ${insight.name}`));
        console.log(chalk.gray(`      Direct Link: ${HOST}/project/${PROJECT_ID}/insights/${insight.id}`));
        console.log('');
    });

    console.log(chalk.yellow('💡 RECOMMENDED STEPS:'));
    console.log(chalk.yellow('1. First, go to "All Insights" to see if your insights are there'));
    console.log(chalk.yellow('2. If you see them, go back to the Dashboard'));
    console.log(chalk.yellow('3. Click "Add insight" in the dashboard'));
    console.log(chalk.yellow('4. Select your insights from the list'));
    console.log('');

    console.log(chalk.green('🎯 QUICK ACCESS:'));
    console.log(chalk.cyan(`   Copy this URL to open insights: ${HOST}/project/${PROJECT_ID}/insights`));
}

// Run the script
if (require.main === module) {
    showNavigationLinks();
}

module.exports = { showNavigationLinks };
