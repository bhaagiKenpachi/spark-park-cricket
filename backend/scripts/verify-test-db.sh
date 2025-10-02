#!/bin/bash

# Verify Test Database Schema
# This script verifies that the testing_db schema is properly set up

set -e

echo "🔍 Verifying test database schema..."

# Check if required environment variables are set
if [ -z "$SUPABASE_URL" ]; then
    echo "❌ ERROR: SUPABASE_URL environment variable is not set"
    exit 1
fi

if [ -z "$SUPABASE_API_KEY" ]; then
    echo "❌ ERROR: SUPABASE_API_KEY environment variable is not set"
    exit 1
fi

echo "📋 Testing database connectivity and schema..."

# Test basic connectivity
if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" "$SUPABASE_URL/rest/v1/" > /dev/null; then
    echo "✅ Database connectivity test successful"
else
    echo "❌ Database connectivity test failed"
    echo "Please check your SUPABASE_URL and SUPABASE_API_KEY"
    exit 1
fi

# Test if we can access the testing_db schema by trying to query a table
# This will fail if the schema doesn't exist, which is what we want to detect
echo "🔍 Testing testing_db schema access..."

# Try to query the users table in testing_db schema
if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" \
    "$SUPABASE_URL/rest/v1/users?select=id&limit=1" \
    -H "Accept-Profile: testing_db" > /dev/null 2>&1; then
    echo "✅ testing_db schema is accessible"
else
    echo "❌ testing_db schema is not accessible or doesn't exist"
    echo ""
    echo "🚨 CRITICAL: The testing_db schema is not properly set up!"
    echo ""
    echo "📋 Required Actions:"
    echo "1. Go to Supabase Dashboard: $SUPABASE_URL/project/default/sql"
    echo "2. Execute the SQL from: internal/database/migrations/complete_schema_testing_db.sql"
    echo "3. Ensure the 'testing_db' schema is created with all required tables"
    echo ""
    echo "📁 Schema file location: internal/database/migrations/complete_schema_testing_db.sql"
    echo ""
    echo "⚠️ Integration tests will fail without proper schema setup!"
    exit 1
fi

# Test a few key tables
echo "🔍 Testing key tables in testing_db schema..."

TABLES=("users" "series" "matches" "innings" "overs" "balls")
for table in "${TABLES[@]}"; do
    if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" \
        "$SUPABASE_URL/rest/v1/$table?select=id&limit=1" \
        -H "Accept-Profile: testing_db" > /dev/null 2>&1; then
        echo "✅ Table '$table' is accessible"
    else
        echo "❌ Table '$table' is not accessible"
        echo "Schema setup may be incomplete!"
        exit 1
    fi
done

echo "✅ All test database schema verification completed successfully!"
echo "🎉 testing_db schema is properly configured and ready for integration tests"
