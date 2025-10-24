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

async function cleanupAndRestart() {
    console.log(chalk.blue.bold('🧹 Cleaning Up Test Insights and Starting Fresh\n'));

    const client = new PostHogAPIClient();

    try {
        // List of test insight IDs to delete
        const testInsightIds = [
            3802021, // Test Insight - Page Views (Proper Format)
            3802089, // Simple Page Views
            3802102, // Basic Page Views - Simple
            3802134, // Test Insight - API Connection
            3802211, // Series Page Views - Properly Saved
            3802256  // Cricket Series Analytics - Short ID Test
        ];

        console.log(chalk.blue('🗑️  Deleting test insights...'));
        let deletedCount = 0;

        for (const insightId of testInsightIds) {
            try {
                await client.deleteInsight(insightId);
                console.log(chalk.green(`   ✅ Deleted insight ${insightId}`));
                deletedCount++;
            } catch (error) {
                console.log(chalk.yellow(`   ⚠️  Could not delete insight ${insightId}: ${error.message}`));
            }
        }

        console.log(chalk.green(`\n✅ Cleaned up ${deletedCount} test insights`));

        // Delete test dashboards
        console.log(chalk.blue('\n🗑️  Deleting test dashboards...'));
        const testDashboardIds = [615946, 615947]; // Cricket Analytics Dashboard (duplicates)

        for (const dashboardId of testDashboardIds) {
            try {
                await client.deleteDashboard(dashboardId);
                console.log(chalk.green(`   ✅ Deleted dashboard ${dashboardId}`));
            } catch (error) {
                console.log(chalk.yellow(`   ⚠️  Could not delete dashboard ${dashboardId}: ${error.message}`));
            }
        }

        console.log(chalk.green.bold('\n🎉 Cleanup Complete!'));
        console.log(chalk.yellow('💡 Ready to create fresh, meaningful cricket analytics'));

        return true;

    } catch (error) {
        console.error(chalk.red('❌ Cleanup failed:'));
        console.error(chalk.red(error.message));
        return false;
    }
}

// Run the script
if (require.main === module) {
    cleanupAndRestart();
}

module.exports = { cleanupAndRestart };
