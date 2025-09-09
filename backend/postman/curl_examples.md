# Spark Park Cricket API - cURL Examples

This document provides cURL command examples for testing the Spark Park Cricket API endpoints.

## Base Configuration

```bash
# Set base URL
BASE_URL="http://localhost:8080"
API_VERSION="v1"

# Common headers
HEADERS="-H 'Content-Type: application/json'"
```

## Health Checks

### Home Endpoint
```bash
curl -X GET "$BASE_URL/"
```

### Health Check
```bash
curl -X GET "$BASE_URL/health"
```

### Database Health
```bash
curl -X GET "$BASE_URL/health/database"
```

### WebSocket Health
```bash
curl -X GET "$BASE_URL/health/websocket"
```

### System Health
```bash
curl -X GET "$BASE_URL/health/system"
```

### Metrics
```bash
curl -X GET "$BASE_URL/health/metrics"
```

## Series Management

### List Series
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/series?limit=10&offset=0"
```

### Create Series
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/series" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "IPL 2024",
    "start_date": "2024-03-22T00:00:00Z",
    "end_date": "2024-05-26T23:59:59Z"
  }'
```

### Get Series
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/series/{series_id}"
```

### Update Series
```bash
curl -X PUT "$BASE_URL/api/$API_VERSION/series/{series_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "IPL 2024 - Updated",
    "end_date": "2024-05-30T23:59:59Z"
  }'
```

### Delete Series
```bash
curl -X DELETE "$BASE_URL/api/$API_VERSION/series/{series_id}"
```

## Team Management

### List Teams
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/teams?limit=20&offset=0"
```

### Create Team
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/teams" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mumbai Indians",
    "players_count": 11
  }'
```

### Get Team
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/teams/{team_id}"
```

### Update Team
```bash
curl -X PUT "$BASE_URL/api/$API_VERSION/teams/{team_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mumbai Indians - Updated",
    "players_count": 12
  }'
```

### List Team Players
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/teams/{team_id}/players"
```

### Add Player to Team
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/teams/{team_id}/players" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rohit Sharma",
    "team_id": "{team_id}"
  }'
```

## Player Management

### List Players
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/players?limit=50&offset=0"
```

### List Players by Team
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/players?team_id={team_id}&limit=50&offset=0"
```

### Create Player
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/players" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Virat Kohli",
    "team_id": "{team_id}"
  }'
```

### Get Player
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/players/{player_id}"
```

### Update Player
```bash
curl -X PUT "$BASE_URL/api/$API_VERSION/players/{player_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Virat Kohli - Captain",
    "team_id": "{team_id}"
  }'
```

### Delete Player
```bash
curl -X DELETE "$BASE_URL/api/$API_VERSION/players/{player_id}"
```

## Match Management

### List Matches
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/matches?limit=20&offset=0"
```

### List Matches by Series
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/matches?series_id={series_id}&limit=20&offset=0"
```

### List Live Matches
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/matches?status=live&limit=20&offset=0"
```

### Create Match
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/matches" \
  -H "Content-Type: application/json" \
  -d '{
    "series_id": "{series_id}",
    "match_number": 1,
    "date": "2024-03-22T19:30:00Z",
    "team1_id": "{team1_id}",
    "team2_id": "{team2_id}"
  }'
```

### Get Match
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/matches/{match_id}"
```

### Update Match
```bash
curl -X PUT "$BASE_URL/api/$API_VERSION/matches/{match_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "live",
    "date": "2024-03-22T20:00:00Z"
  }'
```

### Delete Match
```bash
curl -X DELETE "$BASE_URL/api/$API_VERSION/matches/{match_id}"
```

## Live Scoreboard

### Get Scoreboard
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}"
```

### Add Ball - Good Ball (4 runs)
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "good",
    "runs": 4,
    "is_wicket": false,
    "batsman_id": "{batsman_id}",
    "bowler_id": "{bowler_id}"
  }'
```

### Add Ball - Wide Ball
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "wide",
    "runs": 1,
    "is_wicket": false,
    "batsman_id": "{batsman_id}",
    "bowler_id": "{bowler_id}"
  }'
```

### Add Ball - No Ball
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "no_ball",
    "runs": 2,
    "is_wicket": false,
    "batsman_id": "{batsman_id}",
    "bowler_id": "{bowler_id}"
  }'
```

### Add Ball - Wicket
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "good",
    "runs": 0,
    "is_wicket": true,
    "batsman_id": "{batsman_id}",
    "bowler_id": "{bowler_id}"
  }'
```

### Add Ball - Dead Ball
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "dead_ball",
    "runs": 0,
    "is_wicket": false,
    "batsman_id": "{batsman_id}",
    "bowler_id": "{bowler_id}"
  }'
```

### Update Score
```bash
curl -X PUT "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/score" \
  -H "Content-Type: application/json" \
  -d '{
    "score": 150
  }'
```

### Update Wickets
```bash
curl -X PUT "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/wicket" \
  -H "Content-Type: application/json" \
  -d '{
    "wickets": 3
  }'
```

## WebSocket

### WebSocket Stats
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/ws/stats"
```

### WebSocket Room Stats
```bash
curl -X GET "$BASE_URL/api/$API_VERSION/ws/stats/{match_id}"
```

### Test WebSocket Broadcast
```bash
curl -X POST "$BASE_URL/api/$API_VERSION/ws/test/{match_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Test broadcast message",
    "type": "test"
  }'
```

## Complete Test Workflow

Here's a step-by-step workflow to test the API from scratch:

### 1. Health Check
```
curl -X GET "http://localhost:8080/health"
```

### 2. Create Series
```
curl -X POST "http://localhost:8080/api/v1/series" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Series 2024",
    "start_date": "2024-03-22T00:00:00Z",
    "end_date": "2024-05-26T23:59:59Z"
  }'
```

### 3. Create Teams
Create Team A:
```
curl -X POST "http://localhost:8080/api/v1/teams" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Team A",
    "players_count": 11
  }'
```

Create Team B:
```
curl -X POST "http://localhost:8080/api/v1/teams" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Team B",
    "players_count": 11
  }'
```

### 4. Create Players
Create Batsman:
```
curl -X POST "http://localhost:8080/api/v1/players" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Batsman",
    "team_id": "TEAM_A_ID_FROM_PREVIOUS_RESPONSE"
  }'
```

Create Bowler:
```
curl -X POST "http://localhost:8080/api/v1/players" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Bowler",
    "team_id": "TEAM_B_ID_FROM_PREVIOUS_RESPONSE"
  }'
```

### 5. Create Match
```
curl -X POST "http://localhost:8080/api/v1/matches" \
  -H "Content-Type: application/json" \
  -d '{
    "series_id": "SERIES_ID_FROM_PREVIOUS_RESPONSE",
    "match_number": 1,
    "date": "2024-03-22T19:30:00Z",
    "team1_id": "TEAM_A_ID_FROM_PREVIOUS_RESPONSE",
    "team2_id": "TEAM_B_ID_FROM_PREVIOUS_RESPONSE"
  }'
```

### 6. Test Scoreboard
```
curl -X GET "http://localhost:8080/api/v1/scoreboard/MATCH_ID_FROM_PREVIOUS_RESPONSE"
```

### 7. Add Balls
Add a good ball (4 runs):
```
curl -X POST "http://localhost:8080/api/v1/scoreboard/MATCH_ID/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "good",
    "runs": 4,
    "is_wicket": false,
    "batsman_id": "BATSMAN_ID_FROM_PREVIOUS_RESPONSE",
    "bowler_id": "BOWLER_ID_FROM_PREVIOUS_RESPONSE"
  }'
```

Add a wide ball:
```
curl -X POST "http://localhost:8080/api/v1/scoreboard/MATCH_ID/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "wide",
    "runs": 1,
    "is_wicket": false,
    "batsman_id": "BATSMAN_ID_FROM_PREVIOUS_RESPONSE",
    "bowler_id": "BOWLER_ID_FROM_PREVIOUS_RESPONSE"
  }'
```

### 8. Check Final Scoreboard
```
curl -X GET "http://localhost:8080/api/v1/scoreboard/MATCH_ID"
```

**Note**: Replace the placeholder IDs (like `TEAM_A_ID_FROM_PREVIOUS_RESPONSE`) with the actual IDs returned from the previous API calls.

## Error Handling Examples

### Validation Error
```bash
# This should return a validation error
curl -X POST "$BASE_URL/api/$API_VERSION/series" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "AB",  // Too short
    "start_date": "2024-03-22T00:00:00Z",
    "end_date": "2024-03-21T23:59:59Z"  // End before start
  }'
```

### Not Found Error
```bash
# This should return a 404 error
curl -X GET "$BASE_URL/api/$API_VERSION/series/non-existent-id"
```

### Invalid Ball Event
```bash
# This should return a validation error
curl -X POST "$BASE_URL/api/$API_VERSION/scoreboard/{match_id}/ball" \
  -H "Content-Type: application/json" \
  -d '{
    "ball_type": "good",
    "runs": 7,  // Invalid: more than 6 runs
    "is_wicket": false,
    "batsman_id": "{batsman_id}",
    "bowler_id": "{bowler_id}"
  }'
```

## Tips

1. **Use jq for JSON parsing**: Install jq to parse JSON responses easily
2. **Store IDs**: Save response IDs to environment variables for subsequent requests
3. **Check status codes**: Always check HTTP status codes for error handling
4. **Use verbose mode**: Add `-v` flag to see request/response headers
5. **Pretty print JSON**: Pipe responses through `jq .` for formatted output

## Troubleshooting

### Connection Refused
```bash
# Check if server is running
curl -X GET "$BASE_URL/health"
```

### CORS Issues
```bash
# Add CORS headers if needed
curl -X OPTIONS "$BASE_URL/api/$API_VERSION/series" \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type"
```

### Rate Limiting
```bash
# Check rate limit headers
curl -I "$BASE_URL/api/$API_VERSION/series"
```
