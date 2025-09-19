# Spark Park Cricket API - cURL Examples

## Simplified Cricket Tournament Management System

This document provides cURL examples for the simplified Spark Park Cricket API with Team A vs Team B matches and toss functionality.

## Base URL
```
http://localhost:8080
```

## Health Checks

### Home
```bash
curl --location 'http://localhost:8080/'
```

### Health Check
```bash
curl --location 'http://localhost:8080/health'
```

### Database Health
```bash
curl --location 'http://localhost:8080/health/database'
```

## Series Management

### Create Series
```bash
curl --location 'http://localhost:8080/api/v1/series' \
--header 'Content-Type: application/json' \
--data '{
  "name": "Vijay vs Venkat",
  "start_date": "2024-03-22T00:00:00Z",
  "end_date": "2024-03-23T23:59:59Z"
}'
```

### List Series
```bash
curl --location 'http://localhost:8080/api/v1/series'
```

### Get Series
```bash
curl --location 'http://localhost:8080/api/v1/series/d577f3b7-c8aa-413e-8c43-021f233aaa33'
```

### Update Series
```bash
curl --location --request PUT 'http://localhost:8080/api/v1/series/d577f3b7-c8aa-413e-8c43-021f233aaa33' \
--header 'Content-Type: application/json' \
--data '{
  "name": "Updated Series Name",
  "start_date": "2024-03-22T00:00:00Z",
  "end_date": "2024-03-25T23:59:59Z"
}'
```

### Delete Series
```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/series/d577f3b7-c8aa-413e-8c43-021f233aaa33'
```

## Match Management

### Create Match (with Toss)
```bash
curl --location 'http://localhost:8080/api/v1/matches' \
--header 'Content-Type: application/json' \
--data '{
  "series_id": "d577f3b7-c8aa-413e-8c43-021f233aaa33",
  "match_number": 1,
  "date": "2024-03-22T20:00:00Z",
  "team_a_player_count": 11,
  "team_b_player_count": 11,
  "total_overs": 20,
  "toss_winner": "A",
  "toss_type": "H"
}'
```

### List Matches
```bash
curl --location 'http://localhost:8080/api/v1/matches'
```

### Get Match
```bash
curl --location 'http://localhost:8080/api/v1/matches/99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7'
```

### Update Match
```bash
curl --location --request PUT 'http://localhost:8080/api/v1/matches/99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7' \
--header 'Content-Type: application/json' \
--data '{
  "team_a_player_count": 10,
  "team_b_player_count": 10,
  "total_overs": 15,
  "batting_team": "B"
}'
```

### Delete Match
```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/matches/99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7'
```

### Get Matches by Series
```bash
curl --location 'http://localhost:8080/api/v1/matches/series/d577f3b7-c8aa-413e-8c43-021f233aaa33'
```

## WebSocket

### Match WebSocket Connection
```bash
curl --location 'http://localhost:8080/api/v1/ws/match/99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7'
```

### WebSocket Stats
```bash
curl --location 'http://localhost:8080/api/v1/ws/stats'
```

### Room Stats
```bash
curl --location 'http://localhost:8080/api/v1/ws/stats/99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7'
```

### Test Broadcast
```bash
curl --location --request POST 'http://localhost:8080/api/v1/ws/test/99ba81c4-5a4e-43c7-ac3d-0a0c24e792a7'
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

## Field Descriptions

### Match Fields
- **team_a_player_count**: Number of players in Team A (1-11)
- **team_b_player_count**: Number of players in Team B (1-11)
- **total_overs**: Total overs for the match (1-20)
- **toss_winner**: Team that won the toss ("A" or "B")
- **toss_type**: Toss result ("H" for Heads, "T" for Tails)
- **batting_team**: Team currently batting ("A" or "B")
- **status**: Match status ("live", "completed", "cancelled")

### Run Types (for future scoring)
- **Numbers**: "1", "2", "3", "4", "5", "6", "7", "8", "9"
- **Special**: "NB" (No Ball), "WD" (Wide), "LB" (Leg Byes)

## Notes

1. **Simplified Structure**: The API now uses Team A and Team B instead of complex team management
2. **Toss System**: Matches include toss functionality with Heads/Tails
3. **Live by Default**: All matches start with "live" status
4. **Configurable**: Player counts and overs can be adjusted per match
5. **Real-time**: WebSocket support for live updates (placeholder implementation)

## Error Responses

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

## Scorecard Management

### Start Scoring
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/start' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id"
}'
```

### Add Ball (1 run)
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "1",
  "is_wicket": false,
  "byes": 0
}'
```

### Add Dot Ball (0 runs)
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "0",
  "is_wicket": false,
  "byes": 0
}'
```

### Add Dot Ball with Byes
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "0",
  "is_wicket": false,
  "byes": 1
}'
```

### Add Wide Ball (Extras)
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "wide",
  "run_type": "WD",
  "is_wicket": false,
  "byes": 0
}'
```

### Add No Ball (Extras)
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "no_ball",
  "run_type": "NB",
  "is_wicket": false,
  "byes": 0
}'
```

### Add Wicket
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "WC",
  "is_wicket": true,
  "wicket_type": "bowled",
  "byes": 0
}'
```

### Add Boundary (4 runs)
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 3,
  "ball_type": "good",
  "run_type": "4",
  "is_wicket": false
}'
```

### Add Six (6 runs)
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 4,
  "ball_type": "good",
  "run_type": "6",
  "is_wicket": false
}'
```

### Add Wide Ball
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 5,
  "ball_type": "wide",
  "run_type": "WD",
  "is_wicket": false
}'
```

### Add No Ball
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "no_ball",
  "run_type": "NB",
  "is_wicket": false,
  "byes": 0
}'
```

### Add Ball with Byes
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "ball_type": "good",
  "run_type": "1",
  "is_wicket": false,
  "byes": 2
}'
```

### Get Complete Scorecard
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/your-match-id'
```

### Get Current Over
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/your-match-id/current-over?innings=1'
```

### Get Innings Details
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/your-match-id/innings/1'
```

### Get Over Details
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/your-match-id/innings/1/over/1'
```

### Undo Last Ball
```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/scorecard/your-match-id/ball?innings=1'
```

### Undo Last Ball - Innings 2
```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/scorecard/your-match-id/ball?innings=2'
```

## Scorecard Workflow Example

### Complete Cricket Match Scoring Workflow

1. **Create a Series**
```bash
curl --location 'http://localhost:8080/api/v1/series' \
--header 'Content-Type: application/json' \
--data '{
  "name": "Test Scorecard Series",
  "start_date": "2024-03-22T00:00:00Z",
  "end_date": "2024-03-23T23:59:59Z"
}'
```

2. **Create a Match with Toss**
```bash
curl --location 'http://localhost:8080/api/v1/matches' \
--header 'Content-Type: application/json' \
--data '{
  "series_id": "your-series-id",
  "match_number": 1,
  "date": "2024-03-22T20:00:00Z",
  "team_a_player_count": 11,
  "team_b_player_count": 11,
  "total_overs": 20,
  "toss_winner": "A",
  "toss_type": "H"
}'
```

3. **Start Scoring**
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/start' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id"
}'
```

4. **Add Balls to Complete First Over**
```bash
# Ball 1: 1 run
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 1,
  "ball_type": "good",
  "run_type": "1",
  "is_wicket": false
}'

# Ball 2: 2 runs
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 2,
  "ball_type": "good",
  "run_type": "2",
  "is_wicket": false
}'

# Ball 3: 4 runs (boundary)
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 3,
  "ball_type": "good",
  "run_type": "4",
  "is_wicket": false
}'

# Ball 4: 6 runs (six)
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 4,
  "ball_type": "good",
  "run_type": "6",
  "is_wicket": false
}'

# Ball 5: Wide ball
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 5,
  "ball_type": "wide",
  "run_type": "WD",
  "is_wicket": false
}'

# Ball 6: Wicket
curl --location 'http://localhost:8080/api/v1/scorecard/ball' \
--header 'Content-Type: application/json' \
--data '{
  "match_id": "your-match-id",
  "innings_number": 1,
  "over_number": 1,
  "ball_number": 6,
  "ball_type": "good",
  "run_type": "WC",
  "is_wicket": true,
  "wicket_type": "bowled"
}'
```

5. **View Complete Scorecard**
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/your-match-id'
```

6. **Undo Last Ball (if needed)**
```bash
# Undo the wicket ball from the previous step
curl --location --request DELETE 'http://localhost:8080/api/v1/scorecard/your-match-id/ball?innings=1'
```

7. **View Updated Scorecard**
```bash
curl --location 'http://localhost:8080/api/v1/scorecard/your-match-id'
```

## Run Types and Ball Types

### Run Types
- `"0"`: Dot Ball (0 runs)
- `"1"` to `"9"`: Regular runs (1-9)
- `"NB"`: No Ball (1 run + extra ball)
- `"WD"`: Wide (1 run + extra ball)
- `"LB"`: Leg Byes (count of runs)
- `"WC"`: Wicket (0 runs, wicket taken)

### Ball Types
- `"good"`: Regular ball
- `"wide"`: Wide ball
- `"no_ball"`: No ball
- `"dead_ball"`: Dead ball

### Wicket Types
- `"bowled"`: Bowled
- `"caught"`: Caught
- `"lbw"`: Leg Before Wicket
- `"run_out"`: Run Out
- `"stumped"`: Stumped
- `"hit_wicket"`: Hit Wicket

## Undo Ball Functionality

The undo ball feature allows you to remove the last ball from the current over and automatically recalculate all statistics.

### Key Features:
- **Removes Last Ball**: Deletes the most recently added ball from the current over
- **Statistics Recalculation**: Automatically updates runs, balls, wickets, and overs
- **Status Management**: Properly handles over and innings status changes
- **Validation**: Ensures match is live and innings is in progress

### Usage:
```bash
# Undo last ball from innings 1
curl --location --request DELETE 'http://localhost:8080/api/v1/scorecard/your-match-id/ball?innings=1'

# Undo last ball from innings 2
curl --location --request DELETE 'http://localhost:8080/api/v1/scorecard/your-match-id/ball?innings=2'
```

### What Gets Undone:
- **Runs**: Removes runs from the ball (including byes)
- **Balls**: Decrements legal ball count (good balls only)
- **Wickets**: Removes wicket if the ball was a wicket
- **Overs**: Recalculates over progress and completion status
- **Innings**: Updates innings statistics and completion status

### Error Cases:
- Match not found
- Match not live
- Innings not found or not in progress
- No balls to undo
- Invalid innings number (must be 1 or 2)

## GraphQL Queries

GraphQL allows you to fetch only the fields you need, reducing payload size and improving performance.

### Get Live Scorecard - Minimal
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query GetLiveScorecard($matchId: String!) { liveScorecard(match_id: $matchId) { match_id current_score { runs wickets overs run_rate } } }",
    "variables": {
        "matchId": "d577f3b7-c8aa-413e-8c43-021f233aaa33"
    }
}'
```

### Get Live Scorecard - Detailed
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query GetDetailedScorecard($matchId: String!) { liveScorecard(match_id: $matchId) { match_id match_number series_name team_a team_b total_overs toss_winner toss_type current_innings match_status current_score { runs wickets overs balls run_rate } current_over { over_number total_runs total_balls total_wickets status balls { ball_number ball_type run_type runs byes is_wicket wicket_type } } innings { innings_number batting_team total_runs total_wickets total_overs total_balls status extras { byes leg_byes wides no_balls total } overs { over_number total_runs total_balls total_wickets status balls { ball_number ball_type run_type runs byes is_wicket wicket_type } } } } }",
    "variables": {
        "matchId": "d577f3b7-c8aa-413e-8c43-021f233aaa33"
    }
}'
```

### Get Current Score Only
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query GetCurrentScore($matchId: String!) { liveScorecard(match_id: $matchId) { match_id current_score { runs wickets overs run_rate } } }",
    "variables": {
        "matchId": "d577f3b7-c8aa-413e-8c43-021f233aaa33"
    }
}'
```

### Get Current Over Only
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query GetCurrentOver($matchId: String!) { liveScorecard(match_id: $matchId) { match_id current_over { over_number total_runs total_balls total_wickets status balls { ball_number ball_type run_type runs is_wicket } } } }",
    "variables": {
        "matchId": "d577f3b7-c8aa-413e-8c43-021f233aaa33"
    }
}'
```

### Get Match Summary
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query GetMatchSummary($matchId: String!) { liveScorecard(match_id: $matchId) { match_id match_number series_name team_a team_b total_overs toss_winner current_innings match_status } }",
    "variables": {
        "matchId": "d577f3b7-c8aa-413e-8c43-021f233aaa33"
    }
}'
```

### Get Innings Summary
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query GetInningsSummary($matchId: String!) { liveScorecard(match_id: $matchId) { match_id innings { innings_number batting_team total_runs total_wickets total_overs status extras { total } } } }",
    "variables": {
        "matchId": "d577f3b7-c8aa-413e-8c43-021f233aaa33"
    }
}'
```

### GraphQL Playground
```bash
curl --location 'http://localhost:8080/api/v1/graphql/playground'
```

## GraphQL Benefits

### Field Selection
GraphQL allows you to fetch only the fields you need:

**Minimal payload** - Only current score:
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query { liveScorecard(match_id: \"d577f3b7-c8aa-413e-8c43-021f233aaa33\") { current_score { runs wickets } } }"
}'
```

**Comprehensive payload** - Full scorecard:
```bash
curl --location 'http://localhost:8080/api/v1/graphql' \
--header 'Content-Type: application/json' \
--data '{
    "query": "query { liveScorecard(match_id: \"d577f3b7-c8aa-413e-8c43-021f233aaa33\") { match_id current_score { runs wickets overs run_rate } current_over { over_number total_runs balls { ball_number ball_type runs } } innings { innings_number total_runs total_wickets } } }"
}'
```