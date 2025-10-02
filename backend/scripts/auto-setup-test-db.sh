#!/bin/bash

# Auto Setup Test Database Schema
# This script attempts to automatically set up the testing_db schema using Supabase API

set -e

echo "🚀 Attempting to auto-setup test database schema..."

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

# Test basic connectivity first
echo "🔍 Testing database connectivity..."
if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" "$SUPABASE_URL/rest/v1/" > /dev/null; then
    echo "✅ Database connectivity test successful"
else
    echo "❌ Database connectivity test failed"
    echo "Please check your SUPABASE_URL and SUPABASE_API_KEY"
    exit 1
fi

# Try to create the testing_db schema using SQL execution
echo "🔧 Attempting to create testing_db schema..."

# Extract the schema content
SCHEMA_CONTENT=$(cat "$SCHEMA_FILE")

# Try to execute the schema using Supabase's SQL execution endpoint
# Note: This requires the SQL execution API to be enabled
echo "📝 Executing schema setup..."

# Use curl to execute SQL via Supabase's REST API
# This is a simplified approach - in production, you might need to use Supabase's SQL execution API
RESPONSE=$(curl -s -w "%{http_code}" -X POST \
  -H "apikey: $SUPABASE_API_KEY" \
  -H "Authorization: Bearer $SUPABASE_API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"query\": \"$SCHEMA_CONTENT\"}" \
  "$SUPABASE_URL/rest/v1/rpc/exec_sql" 2>/dev/null || echo "000")

HTTP_CODE="${RESPONSE: -3}"

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Schema setup completed successfully"
else
    echo "⚠️ Automated schema setup failed (HTTP: $HTTP_CODE)"
    echo "This is expected if SQL execution API is not enabled"
    echo "Manual setup required - see instructions below"
fi

# Verify the schema was created
echo "🔍 Verifying schema setup..."
if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" \
    "$SUPABASE_URL/rest/v1/users?select=id&limit=1" \
    -H "Accept-Profile: testing_db" > /dev/null 2>&1; then
    echo "✅ testing_db schema is accessible"
    echo "🎉 Test database schema setup completed successfully!"
    exit 0
else
    echo "❌ testing_db schema is not accessible"
    echo ""
    echo "🚨 Manual schema setup required!"
    echo ""
    echo "📋 Required Actions:"
    echo "1. Go to Supabase Dashboard: $SUPABASE_URL/project/default/sql"
    echo "2. Execute the SQL from: $SCHEMA_FILE"
    echo "3. Ensure the 'testing_db' schema is created with all required tables"
    echo ""
    echo "📄 Schema file preview:"
    echo "=========================================="
    head -20 "$SCHEMA_FILE"
    echo "=========================================="
    echo ""
    echo "⚠️ Integration tests will fail without proper schema setup!"
    exit 1
fi
