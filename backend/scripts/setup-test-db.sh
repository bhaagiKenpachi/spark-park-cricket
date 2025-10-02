#!/bin/bash

# Setup Test Database Schema
# This script sets up the testing_db schema required for integration tests

set -e

echo "🔧 Setting up test database schema..."

# Check if required environment variables are set
if [ -z "$SUPABASE_URL" ]; then
    echo "❌ ERROR: SUPABASE_URL environment variable is not set"
    exit 1
fi

if [ -z "$SUPABASE_API_KEY" ]; then
    echo "❌ ERROR: SUPABASE_API_KEY environment variable is not set"
    exit 1
fi

# Check if the schema file exists
SCHEMA_FILE="internal/database/migrations/complete_schema_testing_db.sql"
if [ ! -f "$SCHEMA_FILE" ]; then
    echo "❌ ERROR: Schema file not found: $SCHEMA_FILE"
    exit 1
fi

echo "📁 Found schema file: $SCHEMA_FILE"
echo "🗄️ Setting up testing_db schema..."

# For now, we'll just log what needs to be done
# In a production environment, you would execute the SQL using Supabase's API
echo "⚠️ NOTE: The testing_db schema needs to be set up manually in Supabase Dashboard"
echo "📋 Required steps:"
echo "   1. Go to Supabase Dashboard: $SUPABASE_URL/project/default/sql"
echo "   2. Execute the SQL from: $SCHEMA_FILE"
echo "   3. Ensure the 'testing_db' schema is created with all required tables"

# Check if we can connect to the database (basic connectivity test)
echo "🔍 Testing database connectivity..."
if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" "$SUPABASE_URL/rest/v1/" > /dev/null; then
    echo "✅ Database connectivity test successful"
else
    echo "❌ Database connectivity test failed"
    echo "Please check your SUPABASE_URL and SUPABASE_API_KEY"
    exit 1
fi

echo "✅ Test database setup script completed"
echo "📝 Manual schema setup required - see instructions above"
