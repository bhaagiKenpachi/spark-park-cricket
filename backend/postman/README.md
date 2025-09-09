# Spark Park Cricket API - Postman Collection

This directory contains Postman collections and environments for testing the Spark Park Cricket Backend API.

## Files

- `Spark_Park_Cricket_API.postman_collection.json` - Main API collection with all endpoints
- `Spark_Park_Cricket_Environment.postman_environment.json` - Environment variables for different stages
- `README.md` - This documentation file

## Import Instructions

1. **Import Collection**: Import `Spark_Park_Cricket_API.postman_collection.json` into Postman
2. **Import Environment**: Import `Spark_Park_Cricket_Environment.postman_environment.json` into Postman
3. **Select Environment**: Choose the appropriate environment in Postman (Development/Staging/Production)

## API Endpoints Overview

### Health Checks
- `GET /` - Welcome message and API version
- `GET /health` - Basic health check
- `GET /health/database` - Database connectivity check
- `GET /health/websocket` - WebSocket service health
- `GET /health/system` - System resources and performance
- `GET /health/ready` - Kubernetes readiness probe
- `GET /health/live` - Kubernetes liveness probe
- `GET /health/metrics` - Application metrics

### Series Management
- `GET /api/v1/series` - List all series/tournaments
- `POST /api/v1/series` - Create new series
- `GET /api/v1/series/{id}` - Get series details
- `PUT /api/v1/series/{id}` - Update series
- `DELETE /api/v1/series/{id}` - Delete series

### Team Management
- `GET /api/v1/teams` - List all teams
- `POST /api/v1/teams` - Create new team
- `GET /api/v1/teams/{id}` - Get team details
- `PUT /api/v1/teams/{id}` - Update team
- `GET /api/v1/teams/{id}/players` - List team players
- `POST /api/v1/teams/{id}/players` - Add player to team

### Player Management
- `GET /api/v1/players` - List all players
- `POST /api/v1/players` - Create new player
- `GET /api/v1/players/{id}` - Get player details
- `PUT /api/v1/players/{id}` - Update player
- `DELETE /api/v1/players/{id}` - Delete player

### Match Management
- `GET /api/v1/matches` - List all matches
- `POST /api/v1/matches` - Create new match
- `GET /api/v1/matches/{id}` - Get match details
- `PUT /api/v1/matches/{id}` - Update match
- `DELETE /api/v1/matches/{id}` - Delete match

### Live Scoreboard
- `GET /api/v1/scoreboard/{match_id}` - Get live scoreboard
- `POST /api/v1/scoreboard/{match_id}/ball` - Add ball event
- `PUT /api/v1/scoreboard/{match_id}/score` - Update score
- `PUT /api/v1/scoreboard/{match_id}/wicket` - Update wickets

### WebSocket
- `GET /api/v1/ws/match/{match_id}` - WebSocket connection
- `GET /api/v1/ws/stats` - WebSocket statistics
- `GET /api/v1/ws/stats/{match_id}` - Match-specific WebSocket stats
- `POST /api/v1/ws/test/{match_id}` - Test WebSocket broadcast

## Environment Variables

The collection uses the following environment variables:

- `base_url` - API base URL (default: http://localhost:8080)
- `api_version` - API version (default: v1)
- `series_id` - Series ID for testing (auto-generated)
- `team_id` - Team ID for testing (auto-generated)
- `team1_id` - First team ID for matches (auto-generated)
- `team2_id` - Second team ID for matches (auto-generated)
- `match_id` - Match ID for testing (auto-generated)
- `player_id` - Player ID for testing (auto-generated)
- `batsman_id` - Batsman ID for scoreboard (auto-generated)
- `bowler_id` - Bowler ID for scoreboard (auto-generated)
- `auth_token` - Authentication token (if required)

## Testing Workflow

### 1. Setup Test Data
1. **Create Series**: Use "Create Series" request to create a tournament
2. **Create Teams**: Use "Create Team" requests to create two teams
3. **Create Players**: Use "Create Player" requests to add players to teams
4. **Create Match**: Use "Create Match" request to create a match between teams

### 2. Test Live Scoring
1. **Get Scoreboard**: Retrieve initial scoreboard for the match
2. **Add Balls**: Use various "Add Ball" requests to simulate match events:
   - Good balls (0-6 runs)
   - Wide balls (extra runs)
   - No balls (extra runs)
   - Wickets
   - Dead balls
3. **Update Score/Wickets**: Test manual score updates

### 3. Test WebSocket
1. **Connect WebSocket**: Use WebSocket connection for real-time updates
2. **Check Stats**: Monitor WebSocket connection statistics
3. **Test Broadcast**: Send test messages to WebSocket clients

## Ball Types and Validation

### Ball Types
- `good` - Regular ball (0-6 runs)
- `wide` - Wide ball (extra runs)
- `no_ball` - No ball (extra runs)
- `dead_ball` - Dead ball (no runs, not counted)

### Validation Rules
- **Runs**: 0-6 for good balls, 0+ for extras
- **Wickets**: 0-10 maximum
- **Ball Types**: Must be one of the valid types
- **Player IDs**: Must be valid player IDs
- **Match Status**: scheduled, live, completed, cancelled

## Error Handling

The API returns standardized error responses:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "field": "runs",
      "reason": "must be between 0 and 6"
    }
  }
}
```

## Rate Limiting

- **Rate Limit**: 100 requests per minute per IP
- **Timeout**: 60 seconds for requests
- **CORS**: Enabled for all origins

## WebSocket Testing

For WebSocket testing in Postman:

1. Use the "WebSocket Connection" request
2. Connect to the WebSocket endpoint
3. Send messages and receive real-time updates
4. Monitor connection statistics

## Sample Test Data

### Series
```json
{
  "name": "IPL 2024",
  "start_date": "2024-03-22T00:00:00Z",
  "end_date": "2024-05-26T23:59:59Z"
}
```

### Team
```json
{
  "name": "Mumbai Indians",
  "players_count": 11
}
```

### Player
```json
{
  "name": "Rohit Sharma",
  "team_id": "{{team_id}}"
}
```

### Match
```json
{
  "series_id": "{{series_id}}",
  "match_number": 1,
  "date": "2024-03-22T19:30:00Z",
  "team1_id": "{{team1_id}}",
  "team2_id": "{{team2_id}}"
}
```

### Ball Event
```json
{
  "ball_type": "good",
  "runs": 4,
  "is_wicket": false,
  "batsman_id": "{{batsman_id}}",
  "bowler_id": "{{bowler_id}}"
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure the backend server is running on the correct port
2. **CORS Errors**: Check CORS configuration in the backend
3. **Validation Errors**: Verify request body format and required fields
4. **WebSocket Connection Failed**: Ensure WebSocket endpoint is accessible

### Debug Tips

1. Check server logs for detailed error messages
2. Use Postman Console to view request/response details
3. Verify environment variables are set correctly
4. Test with curl commands for comparison

## Development Setup

1. Start the backend server: `go run cmd/server/main.go`
2. Server runs on `http://localhost:8080` by default
3. Import the Postman collection
4. Select the appropriate environment
5. Start testing the API endpoints

## Contributing

When adding new endpoints:

1. Update the Postman collection with new requests
2. Add appropriate test scripts
3. Update this README with new endpoint documentation
4. Ensure environment variables are properly configured
