# Spark Park Cricket API - Postman Collection

## Overview

This Postman collection provides a simplified cricket tournament management system with Team A vs Team B matches and toss functionality.

## Features

- **Simplified Team Structure**: Team A and Team B with configurable player counts
- **Toss System**: Heads/Tails toss with winner batting first
- **Match Management**: Create, update, and manage cricket matches
- **Series Management**: Organize matches into tournaments/series
- **Real-time Updates**: WebSocket support for live match updates
- **Live by Default**: All matches start with "live" status

## Setup

1. **Import Collection**: Import `Spark_Park_Cricket_API.postman_collection.json`
2. **Import Environment**: Import `Spark_Park_Cricket_Environment.postman_environment.json`
3. **Start Server**: Ensure the backend server is running on `http://localhost:8080`
4. **Select Environment**: Choose "Spark Park Cricket Environment" in Postman

## API Structure

### Health Checks
- `GET /` - Home endpoint
- `GET /health` - Basic health check
- `GET /health/database` - Database health check

### Series Management
- `POST /api/v1/series` - Create series
- `GET /api/v1/series` - List all series
- `GET /api/v1/series/{id}` - Get specific series
- `PUT /api/v1/series/{id}` - Update series
- `DELETE /api/v1/series/{id}` - Delete series

### Match Management
- `POST /api/v1/matches` - Create match with toss
- `GET /api/v1/matches` - List all matches
- `GET /api/v1/matches/{id}` - Get specific match
- `PUT /api/v1/matches/{id}` - Update match
- `DELETE /api/v1/matches/{id}` - Delete match
- `GET /api/v1/matches/series/{series_id}` - Get matches by series

### Scorecard Management
- `POST /api/v1/scorecard/start` - Start scoring for a match
- `POST /api/v1/scorecard/ball` - Add ball events
- `GET /api/v1/scorecard/{match_id}` - Get complete scorecard
- `GET /api/v1/scorecard/{match_id}/current-over` - Get current over
- `GET /api/v1/scorecard/{match_id}/innings/{innings_number}` - Get innings details
- `GET /api/v1/scorecard/{match_id}/innings/{innings_number}/over/{over_number}` - Get over details

### WebSocket
- `GET /api/v1/ws/match/{match_id}` - WebSocket connection
- `GET /api/v1/ws/stats` - Connection statistics
- `GET /api/v1/ws/stats/{match_id}` - Room statistics
- `POST /api/v1/ws/test/{match_id}` - Test broadcast

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `base_url` | API base URL | `http://localhost:8080` |
| `series_id` | Series UUID | `d577f3b7-c8aa-413e-8c43-021f233aaa33` |
| `match_id` | Match UUID | `99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7` |
| `team_a_player_count` | Team A players | `11` |
| `team_b_player_count` | Team B players | `11` |
| `total_overs` | Match overs | `20` |
| `toss_winner` | Toss winner | `A` or `B` |
| `toss_type` | Toss result | `H` (Heads) or `T` (Tails) |
| `batting_team` | Current batting team | `A` or `B` |

## Example Workflow

### 1. Create Series
```json
POST /api/v1/series
{
  "name": "Vijay vs Venkat",
  "start_date": "2024-03-22T00:00:00Z",
  "end_date": "2024-03-23T23:59:59Z"
}
```

### 2. Create Match with Toss
```json
POST /api/v1/matches
{
  "series_id": "{{series_id}}",
  "match_number": 1,
  "date": "2024-03-22T20:00:00Z",
  "team_a_player_count": 11,
  "team_b_player_count": 11,
  "total_overs": 20,
  "toss_winner": "A",
  "toss_type": "H"
}
```

### 3. Start Scoring
```json
POST /api/v1/scorecard/start
{
  "match_id": "{{match_id}}"
}
```

### 4. Add Ball Events
```json
POST /api/v1/scorecard/ball
{
  "match_id": "{{match_id}}",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "1",
  "is_wicket": false,
  "byes": 0
}
```

### 5. Add Wicket
```json
POST /api/v1/scorecard/ball
{
  "match_id": "{{match_id}}",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "WC",
  "is_wicket": true,
  "wicket_type": "bowled",
  "byes": 0
}
```

### 6. Add Wide Ball with Byes
```json
POST /api/v1/scorecard/ball
{
  "match_id": "{{match_id}}",
  "innings_number": 1,
  "ball_type": "wide",
  "run_type": "WD",
  "is_wicket": false,
  "byes": 2
}
```

### 7. Get Scorecard
```json
GET /api/v1/scorecard/{{match_id}}
```

### 8. Update Match (Change Batting Team)
```json
PUT /api/v1/matches/{{match_id}}
{
  "batting_team": "B"
}
```

## Data Models

### Series
```json
{
  "id": "uuid",
  "name": "Series Name",
  "start_date": "2024-03-22T00:00:00Z",
  "end_date": "2024-03-23T23:59:59Z",
  "created_at": "2024-03-22T00:00:00Z",
  "updated_at": "2024-03-22T00:00:00Z"
}
```

### Match
```json
{
  "id": "uuid",
  "series_id": "uuid",
  "match_number": 1,
  "date": "2024-03-22T20:00:00Z",
  "status": "live",
  "team_a_player_count": 11,
  "team_b_player_count": 11,
  "total_overs": 20,
  "toss_winner": "A",
  "toss_type": "H",
  "batting_team": "A",
  "created_at": "2024-03-22T00:00:00Z",
  "updated_at": "2024-03-22T00:00:00Z"
}
```

## Ball Types and Run Types

### Ball Types
- `good` - Legal delivery
- `wide` - Wide ball (extra)
- `no_ball` - No ball (extra)
- `dead_ball` - Dead ball

### Run Types
- `0` - Dot ball (0 runs)
- `1-9` - Regular runs (1-9)
- `NB` - No ball (1 run + extra)
- `WD` - Wide (1 run + extra)
- `LB` - Leg byes
- `WC` - Wicket

### Ball Event Request Structure
```json
{
    "match_id": "uuid",
    "innings_number": 1,
    "ball_type": "good",
    "run_type": "1",
    "is_wicket": false,
    "wicket_type": "bowled", // Required if is_wicket is true
    "byes": 0 // Additional runs from byes (0-6)
}
```

### Wicket Types
- `bowled` - Bowled out
- `caught` - Caught out
- `lbw` - Leg before wicket
- `run_out` - Run out
- `stumped` - Stumped
- `hit_wicket` - Hit wicket

## Field Validation

### Match Fields
- **team_a_player_count**: 1-11 players
- **team_b_player_count**: 1-11 players
- **total_overs**: 1-20 overs
- **toss_winner**: "A" or "B"
- **toss_type**: "H" (Heads) or "T" (Tails)
- **batting_team**: "A" or "B"
- **status**: "live", "completed", "cancelled"

## Auto-Generated Features

The collection includes automatic features:

1. **Auto UUID Generation**: Generates UUIDs for series_id and match_id if not set
2. **Auto ID Extraction**: Extracts IDs from responses and sets environment variables
3. **Pre-request Scripts**: Automatically populate missing variables
4. **Test Scripts**: Validate responses and extract data

## Error Handling

All endpoints return standardized error responses:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

Common error codes:
- `VALIDATION_ERROR`: Invalid input data
- `NOT_FOUND`: Resource not found
- `INTERNAL_ERROR`: Server error

## Testing

1. **Health Check**: Start with health endpoints to verify server
2. **Create Series**: Create a test series
3. **Create Match**: Create a match with toss
4. **Start Scoring**: Begin scoring for the match
5. **Add Balls**: Test different ball types and run combinations
6. **Add Wickets**: Test wicket scenarios
7. **Test Extras**: Try wide balls, no balls, and byes
8. **Get Scorecard**: Verify complete scorecard data
9. **WebSocket**: Test real-time connections
10. **Update Match**: Test match updates

## Notes

- **Simplified Architecture**: No complex team/player management
- **Toss Integration**: Toss winner automatically bats first
- **Live by Default**: Matches start in "live" status
- **Configurable**: Adjust player counts and overs per match
- **Ball-by-Ball Scoring**: Complete tracking of all ball types and run types
- **Extras Support**: Wide balls, no balls, byes, leg byes
- **Wicket Tracking**: Multiple wicket types with proper validation
- **Automatic Over Management**: Automatic over and ball number detection
- **Match Completion Logic**: Automatic innings and match completion handling
- **Real-time Ready**: WebSocket infrastructure for live scoring

## Support

For issues or questions:
1. Check server logs for detailed error messages
2. Verify environment variables are set correctly
3. Ensure server is running on the correct port
4. Check database connection health