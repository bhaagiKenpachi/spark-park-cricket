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

class PostHogAPIClient {
    constructor() {
        this.apiKey = process.env.POSTHOG_PERSONAL_API_KEY;
        this.projectId = process.env.POSTHOG_PROJECT_ID;
        this.host = process.env.POSTHOG_HOST || 'https://us.i.posthog.com';

        if (!this.apiKey || !this.projectId) {
            console.error(chalk.red('❌ Missing required environment variables:'));
            if (!this.apiKey) console.error(chalk.red('   - POSTHOG_PERSONAL_API_KEY'));
            if (!this.projectId) console.error(chalk.red('   - POSTHOG_PROJECT_ID'));
            console.error(chalk.yellow('\n💡 Please add these variables to your main .env file in the web directory.'));
            console.error(chalk.yellow('   See scripts/posthog/README.md for instructions on getting these values.'));
            process.exit(1);
        }

        // Check if API key looks like a Project API Key instead of Personal API Key
        if (this.apiKey && this.apiKey.startsWith('phc_')) {
            console.error(chalk.red('❌ Invalid API Key Type:'));
            console.error(chalk.red('   You are using a Project API Key (phc_...) but need a Personal API Key (phx_...)'));
            console.error(chalk.yellow('\n💡 Please generate a Personal API Key from PostHog Settings > Personal API Keys'));
            console.error(chalk.yellow('   The Project API Key is used for sending events, Personal API Key is for dashboard management.'));
            process.exit(1);
        }
    }

    async makeRequest(endpoint, options = {}) {
        const url = `${this.host}/api/projects/${this.projectId}${endpoint}`;

        const defaultOptions = {
            headers: {
                'Authorization': `Bearer ${this.apiKey}`,
                'Content-Type': 'application/json',
            },
        };

        const requestOptions = {
            ...defaultOptions,
            ...options,
            headers: {
                ...defaultOptions.headers,
                ...options.headers,
            },
        };

        try {
            const response = await fetch(url, requestOptions);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(`API Error: ${response.status} - ${data.detail || data.message || 'Unknown error'}`);
            }

            return data;
        } catch (error) {
            console.error(chalk.red(`❌ API Request failed: ${error.message}`));
            throw error;
        }
    }

    async createInsight(insightConfig) {
        console.log(chalk.blue(`📊 Creating insight: ${insightConfig.name}`));

        const response = await this.makeRequest('/insights/', {
            method: 'POST',
            body: JSON.stringify(insightConfig),
        });

        console.log(chalk.green(`✅ Created insight: ${insightConfig.name} (ID: ${response.id})`));
        if (response.short_id) {
            console.log(chalk.cyan(`   Short ID: ${response.short_id}`));
            console.log(chalk.cyan(`   URL: ${this.getInsightUrl(response.short_id)}`));
        }
        return response;
    }

    async createDashboard(dashboardConfig) {
        console.log(chalk.blue(`📋 Creating dashboard: ${dashboardConfig.name}`));

        const response = await this.makeRequest('/dashboards/', {
            method: 'POST',
            body: JSON.stringify(dashboardConfig),
        });

        console.log(chalk.green(`✅ Created dashboard: ${dashboardConfig.name} (ID: ${response.id})`));
        return response;
    }

    async addInsightToDashboard(dashboardId, insightId) {
        console.log(chalk.blue(`🔗 Adding insight ${insightId} to dashboard ${dashboardId}`));

        try {
            // Get current dashboard
            const dashboard = await this.makeRequest(`/dashboards/${dashboardId}/`);

            // Add insight to tiles
            const tiles = dashboard.tiles || [];
            tiles.push({
                insight: insightId,
                layouts: {
                    sm: { x: (tiles.length % 2) * 6, y: Math.floor(tiles.length / 2) * 4, w: 6, h: 4 },
                    lg: { x: (tiles.length % 3) * 4, y: Math.floor(tiles.length / 3) * 4, w: 4, h: 4 }
                }
            });

            const response = await this.makeRequest(`/dashboards/${dashboardId}/`, {
                method: 'PATCH',
                body: JSON.stringify({ tiles }),
            });

            console.log(chalk.green(`✅ Added insight to dashboard`));
            return response;
        } catch (error) {
            console.log(chalk.yellow(`⚠️  Could not add insight to dashboard automatically. You can add it manually in PostHog UI.`));
            console.log(chalk.gray(`   Dashboard URL: ${this.getDashboardUrl(dashboardId)}`));
            console.log(chalk.gray(`   Insight ID: ${insightId}`));
            // Don't throw error, just log warning and continue
            return null;
        }
    }

    async listDashboards() {
        console.log(chalk.blue('📋 Fetching dashboards...'));

        const response = await this.makeRequest('/dashboards/');

        if (response.results && response.results.length > 0) {
            console.log(chalk.green('\n📊 Available Dashboards:'));
            response.results.forEach(dashboard => {
                console.log(chalk.cyan(`  • ${dashboard.name} (ID: ${dashboard.id})`));
                console.log(chalk.gray(`    ${dashboard.description || 'No description'}`));
            });
        } else {
            console.log(chalk.yellow('No dashboards found.'));
        }

        return response.results || [];
    }

    async deleteDashboard(dashboardId) {
        console.log(chalk.blue(`🗑️  Deleting dashboard: ${dashboardId}`));

        await this.makeRequest(`/dashboards/${dashboardId}/`, {
            method: 'DELETE',
        });

        console.log(chalk.green(`✅ Deleted dashboard: ${dashboardId}`));
    }

    getDashboardUrl(dashboardId) {
        return `${this.host}/project/${this.projectId}/dashboard/${dashboardId}`;
    }

    getInsightUrl(shortId) {
        return `${this.host}/project/${this.projectId}/insights/${shortId}`;
    }

    async deleteInsight(insightId) {
        console.log(chalk.blue(`🗑️  Deleting insight: ${insightId}`));

        await this.makeRequest(`/insights/${insightId}/`, {
            method: 'DELETE',
        });

        console.log(chalk.green(`✅ Deleted insight: ${insightId}`));
    }
}

module.exports = PostHogAPIClient;
