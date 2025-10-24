#!/usr/bin/env node

const { Command } = require('commander');
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
const { createSeriesDashboard } = require('./create-series-dashboard');

const program = new Command();

program
    .name('posthog-dashboard')
    .description('Manage PostHog dashboards programmatically')
    .version('1.0.0');

program
    .command('create-series')
    .description('Create Series Analytics Dashboard with predefined insights')
    .action(async () => {
        try {
            await createSeriesDashboard();
        } catch (error) {
            console.error(chalk.red('Failed to create series dashboard:', error.message));
            process.exit(1);
        }
    });

program
    .command('list')
    .description('List all dashboards in the project')
    .action(async () => {
        try {
            const client = new PostHogAPIClient();
            await client.listDashboards();
        } catch (error) {
            console.error(chalk.red('Failed to list dashboards:', error.message));
            process.exit(1);
        }
    });

program
    .command('delete <id>')
    .description('Delete a dashboard by ID')
    .action(async (id) => {
        try {
            const client = new PostHogAPIClient();
            await client.deleteDashboard(id);
        } catch (error) {
            console.error(chalk.red('Failed to delete dashboard:', error.message));
            process.exit(1);
        }
    });

program
    .command('info')
    .description('Show configuration information')
    .action(() => {
        console.log(chalk.blue.bold('📊 PostHog Dashboard Manager'));
        console.log(chalk.gray('Configuration:'));
        console.log(chalk.gray(`  Host: ${process.env.POSTHOG_HOST || 'https://us.i.posthog.com'}`));
        console.log(chalk.gray(`  Project ID: ${process.env.POSTHOG_PROJECT_ID || 'Not set'}`));
        console.log(chalk.gray(`  API Key: ${process.env.POSTHOG_PERSONAL_API_KEY ? 'Set' : 'Not set'}`));

        if (!process.env.POSTHOG_PERSONAL_API_KEY || !process.env.POSTHOG_PROJECT_ID) {
            console.log(chalk.yellow('\n⚠️  Missing configuration. Please check your .env file.'));
        }
    });

// Add help for common issues
program.on('--help', () => {
    console.log(chalk.yellow('\n💡 Common Issues:'));
    console.log(chalk.yellow('  • Authentication failed: Check your Personal API Key'));
    console.log(chalk.yellow('  • Project not found: Verify your Project ID'));
    console.log(chalk.yellow('  • Missing .env file: Create scripts/posthog/.env with your credentials'));
    console.log(chalk.yellow('\n📖 For detailed setup instructions, see scripts/posthog/README.md'));
});

program.parse();
