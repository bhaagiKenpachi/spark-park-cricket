# CI Integration Test Fix Guide

## Issue
The GitHub Actions CI is failing with exit code 1 on Integration Tests. Based on the [CI run logs](https://github.com/bhaagiKenpachi/spark-park-cricket/actions/runs/18154509923/job/51671390348), the tests are failing because the `testing_db` schema is not properly set up in Supabase.

## Root Cause
The integration tests require a `testing_db` schema in Supabase with specific tables and permissions. The CI environment doesn't have this schema set up, causing all database operations to fail.

## Solution

### Step 1: Set Up the Testing Database Schema

1. **Go to Supabase Dashboard**
   - Navigate to: `https://your-project.supabase.co/project/default/sql`
   - Replace `your-project` with your actual Supabase project identifier

2. **Execute the Schema SQL**
   - Copy the contents of `internal/database/migrations/complete_schema_testing_db.sql`
   - Paste it into the SQL Editor in Supabase Dashboard
   - Execute the SQL

3. **Verify Schema Creation**
   - The script should create a `testing_db` schema with the following tables:
     - `users`
     - `user_sessions`
     - `oauth_states`
     - `series`
     - `matches`
     - `innings`
     - `overs`
     - `balls`

### Step 2: Verify the Setup

You can verify the schema is set up correctly by running:

```bash
cd backend
./scripts/verify-test-db.sh
```

### Step 3: Test Locally

Run the integration tests locally to ensure they work:

```bash
cd backend
go test ./tests/integration/... -v -race -timeout=15m -count=1
```

### Step 4: Trigger CI

After setting up the schema, push your changes or trigger a new CI run. The integration tests should now pass.

## Alternative: Manual Schema Setup

If the automated setup doesn't work, you can manually set up the schema:

### 1. Create the Schema

```sql
CREATE SCHEMA IF NOT EXISTS testing_db;
```

### 2. Set Up Tables

Execute the complete SQL from `internal/database/migrations/complete_schema_testing_db.sql` in the Supabase SQL Editor.

### 3. Verify Permissions

Ensure the service role has access to the `testing_db` schema:

```sql
GRANT ALL ON SCHEMA testing_db TO service_role;
GRANT ALL ON ALL TABLES IN SCHEMA testing_db TO service_role;
GRANT ALL ON ALL SEQUENCES IN SCHEMA testing_db TO service_role;
```

## CI Improvements Applied

The following improvements have been made to the CI workflow:

1. **Added Race Detection**: Tests now run with `-race` flag to detect race conditions
2. **Improved Error Handling**: Better error messages and debugging output
3. **Database Schema Verification**: Automated verification of test database setup
4. **Enhanced Logging**: More detailed environment and test information
5. **Auto-Setup Scripts**: Attempts to automatically set up the test database

## Troubleshooting

### If Tests Still Fail

1. **Check CI Logs**: Look for the "Environment Details" section to see configuration
2. **Verify Database Access**: Ensure `SUPABASE_URL` and `SUPABASE_API_KEY` secrets are set
3. **Check Schema**: Verify the `testing_db` schema exists and has all required tables
4. **Test Locally**: Run tests locally with the same configuration as CI

### Common Issues

1. **Missing Schema**: The `testing_db` schema doesn't exist
2. **Missing Tables**: Required tables are missing from the schema
3. **Permission Issues**: Service role doesn't have access to the schema
4. **Environment Variables**: CI secrets are not properly configured

## Files Modified

- `.github/workflows/backend-ci.yml` - Enhanced CI workflow with better error handling
- `scripts/verify-test-db.sh` - Script to verify database schema setup
- `scripts/auto-setup-test-db.sh` - Script to attempt automatic schema setup
- `CI_FIX_GUIDE.md` - This troubleshooting guide

## Next Steps

1. Set up the `testing_db` schema in Supabase
2. Verify the setup using the verification script
3. Test locally to ensure everything works
4. Push changes to trigger CI
5. Monitor CI logs for successful test execution

The integration tests should now pass in the CI environment! 🎉
