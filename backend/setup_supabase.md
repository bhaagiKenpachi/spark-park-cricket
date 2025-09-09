# Supabase Setup Guide for Spark Park Cricket Backend

This guide will help you set up a test database in Supabase and start the backend server for API testing.

## Prerequisites

- Go 1.23+ installed
- Supabase account (free tier is sufficient)
- Postman installed

## Step 1: Create Supabase Project

### 1.1 Go to Supabase Dashboard
1. Visit [https://supabase.com](https://supabase.com)
2. Sign in or create a free account
3. Click "New Project"

### 1.2 Create New Project
1. **Organization**: Select your organization (or create one)
2. **Name**: `spark-park-cricket-test`
3. **Database Password**: Generate a strong password (save it!)
4. **Region**: Choose closest to your location
5. **Pricing Plan**: Free tier is sufficient for testing
6. Click "Create new project"

### 1.3 Wait for Setup
- Project creation takes 1-2 minutes
- You'll see a progress indicator

## Step 2: Get Supabase Credentials

### 2.1 Get Project URL and API Key
1. In your Supabase dashboard, go to **Settings** → **API**
2. Copy the following values:
   - **Project URL** (looks like: `https://your-project-ref.supabase.co`)
   - **anon public** key (starts with `eyJ...`)

### 2.2 Get Database Connection Details
1. Go to **Settings** → **Database**
2. Copy the **Connection string** (you'll need the password you set earlier)

## Step 3: Configure Environment Variables

### 3.1 Create .env File
Create a `.env` file in the backend directory:

1. Navigate to the backend directory: `/Users/luffybhaagi/dojima/spark-park-cricket/backend`
2. Create a new file named `.env`
3. Copy the contents from `env.example` and edit with your credentials

### 3.2 Add Environment Variables
Edit the `.env` file with your Supabase credentials:

```env
# Supabase Configuration
SUPABASE_URL=https://your-project-ref.supabase.co
SUPABASE_API_KEY=your_anon_public_key_here

# Database Configuration (for migrations)
DATABASE_PASSWORD=your_database_password

# Server Configuration
PORT=8080
```

**Replace the following:**
- `your-project-ref` with your actual project reference
- `your_anon_public_key_here` with your actual anon public key
- `your_database_password` with the database password you set

## Step 4: Run Database Migrations

### 4.1 Using Go Migration Tool (Recommended)
1. Install dependencies: `go mod tidy`
2. Run migrations: `go run cmd/migrate/main.go`

This will automatically:
- Connect to your Supabase database
- Create the `schema_migrations` table
- Apply all pending migrations from `internal/database/migrations/`
- Track applied migrations to prevent duplicates

### 4.2 Alternative: Manual Migration via Supabase Dashboard
1. Go to **SQL Editor** in your Supabase dashboard
2. Copy the contents of `internal/database/migrations/001_initial_schema.sql`
3. Paste it into the SQL editor
4. Click "Run" to execute the migration

### 4.3 Verify Tables Created
1. Go to **Table Editor** in Supabase dashboard
2. You should see the following tables:
   - `series`
   - `teams`
   - `players`
   - `matches`
   - `live_scoreboard`
   - `overs`
   - `balls`

## Step 5: Start Backend Server

### 5.1 Install Dependencies
1. Navigate to the backend directory: `/Users/luffybhaagi/dojima/spark-park-cricket/backend`
2. Run: `go mod tidy` to download Go dependencies

### 5.2 Start Server
Run: `go run cmd/server/main.go` to start the server

You should see output like:
```
Starting Spark Park Cricket Backend Server...
Server running on port 8080
Database connected successfully
WebSocket hub started
```

## Step 6: Test API Connection

### 6.1 Test Health Endpoint
Test the basic health check by running:
`curl http://localhost:8080/health`

Expected response:
```json
{"success":true,"data":{"status":"OK","service":"spark-park-cricket-backend"}}
```

### 6.2 Test with Postman
1. Open Postman
2. Import the collection: `postman/Spark_Park_Cricket_API.postman_collection.json`
3. Import the environment: `postman/Spark_Park_Cricket_Environment.postman_environment.json`
4. Set the environment to use `http://localhost:8080`
5. Test the "Home" endpoint first

## Step 7: Create Test Data

### 7.1 Create a Series
Run this command to create a test series:
```
curl -X POST http://localhost:8080/api/v1/series \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Series 2024",
    "start_date": "2024-03-22T00:00:00Z",
    "end_date": "2024-05-26T23:59:59Z"
  }'
```

### 7.2 Create Teams
Create Team 1:
```
curl -X POST http://localhost:8080/api/v1/teams \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mumbai Indians",
    "players_count": 11
  }'
```

Create Team 2:
```
curl -X POST http://localhost:8080/api/v1/teams \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chennai Super Kings",
    "players_count": 11
  }'
```

## Troubleshooting

### Common Issues

#### 1. Connection Refused
- Check if server is running by looking for the Go process
- Check if port 8080 is in use by another application

#### 2. Database Connection Error
- Verify your `.env` file has correct credentials
- Check Supabase project is active
- Ensure API key has correct permissions

#### 3. Migration Errors
- Check if UUID extension is enabled in Supabase
- Verify table names don't conflict
- Check for syntax errors in SQL

#### 4. CORS Issues
- The server has CORS enabled for all origins
- If testing from browser, ensure proper headers

### Debug Steps

1. **Check environment variables** - Verify your `.env` file has the correct values
2. **Test database connection** - Run: `curl -X GET http://localhost:8080/health/database`
3. **Check server logs** - Look at the console output when starting the server

## Next Steps

1. **Test All Endpoints**: Use the Postman collection to test all API endpoints
2. **Create Test Data**: Set up series, teams, players, and matches
3. **Test Live Scoring**: Try the scoreboard and ball tracking features
4. **Test WebSocket**: Connect to WebSocket endpoints for real-time updates

## Migration Tool

The project includes a Go-based migration tool for easy database management:

### Features
- ✅ **Automatic Migration Tracking**: Tracks applied migrations
- ✅ **Transaction Safety**: Each migration runs in a transaction
- ✅ **Version Control**: Migrations ordered by version number
- ✅ **Idempotent**: Safe to run multiple times

### Usage
```bash
# Run all pending migrations
go run cmd/migrate/main.go
```

### Documentation
- [Migration Tool README](./cmd/migrate/README.md)

## Useful Links

- [Supabase Dashboard](https://supabase.com/dashboard)
- [Supabase Documentation](https://supabase.com/docs)
- [Postman Collection Documentation](./postman/README.md)
- [API cURL Examples](./postman/curl_examples.md)

## Support

If you encounter issues:
1. Check the server logs for error messages
2. Verify your Supabase credentials
3. Ensure all dependencies are installed
4. Check the troubleshooting section above
