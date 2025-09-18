# GraphQL Scorecard Queries Examples

This document provides comprehensive examples of GraphQL queries for fetching cricket scorecard data.

## Available Queries

### 1. Live Scorecard
Get the complete live scorecard for a match with current score and all innings data.

```graphql
query GetLiveScorecard($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_id
    match_number
    series_name
    team_a
    team_b
    total_overs
    toss_winner
    toss_type
    current_innings
    match_status
    current_score {
      runs
      wickets
      overs
      balls
      run_rate
    }
    current_over {
      over_number
      total_runs
      total_balls
      total_wickets
      status
      balls {
        ball_number
        ball_type
        run_type
        runs
        byes
        is_wicket
        wicket_type
      }
    }
    innings {
      innings_number
      batting_team
      total_runs
      total_wickets
      total_overs
      total_balls
      status
      extras {
        byes
        leg_byes
        wides
        no_balls
        total
      }
      overs {
        over_number
        total_runs
        total_balls
        total_wickets
        status
        balls {
          ball_number
          ball_type
          run_type
          runs
          byes
          is_wicket
          wicket_type
        }
      }
    }
  }
}
```

### 2. Match Details
Get basic match information and details.

```graphql
query GetMatchDetails($matchId: String!) {
  matchDetails(match_id: $matchId) {
    match_id
    match_number
    series_name
    team_a
    team_b
    total_overs
    toss_winner
    toss_type
    current_innings
    match_status
    batting_team
    team_a_player_count
    team_b_player_count
  }
}
```

### 3. Match Statistics
Get aggregated statistics for the entire match.

```graphql
query GetMatchStatistics($matchId: String!) {
  matchStatistics(match_id: $matchId) {
    total_runs
    total_wickets
    total_overs
    total_balls
    run_rate
    extras {
      byes
      leg_byes
      wides
      no_balls
      total
    }
    innings_count
  }
}
```

### 4. Innings Score
Get basic score information for a specific innings.

```graphql
query GetInningsScore($matchId: String!, $inningsNumber: Int!) {
  inningsScore(match_id: $matchId, innings_number: $inningsNumber) {
    innings_number
    batting_team
    total_runs
    total_wickets
    total_overs
    total_balls
    status
    extras {
      byes
      leg_byes
      wides
      no_balls
      total
    }
  }
}
```

### 5. Innings Details
Get complete details for a specific innings including all overs.

```graphql
query GetInningsDetails($matchId: String!, $inningsNumber: Int!) {
  inningsDetails(match_id: $matchId, innings_number: $inningsNumber) {
    innings_number
    batting_team
    total_runs
    total_wickets
    total_overs
    total_balls
    status
    extras {
      byes
      leg_byes
      wides
      no_balls
      total
    }
    overs {
      over_number
      total_runs
      total_balls
      total_wickets
      status
      balls {
        ball_number
        ball_type
        run_type
        runs
        byes
        is_wicket
        wicket_type
      }
    }
  }
}
```

### 6. Over Details
Get details for a specific over in a specific innings.

```graphql
query GetOverDetails($matchId: String!, $inningsNumber: Int!, $overNumber: Int!) {
  overDetails(match_id: $matchId, innings_number: $inningsNumber, over_number: $overNumber) {
    over_number
    total_runs
    total_balls
    total_wickets
    status
    balls {
      ball_number
      ball_type
      run_type
      runs
      byes
      is_wicket
      wicket_type
    }
  }
}
```

### 7. Latest Over
Get the current/latest over for a specific innings.

```graphql
query GetLatestOver($matchId: String!, $inningsNumber: Int!) {
  latestOver(match_id: $matchId, innings_number: $inningsNumber) {
    over_number
    total_runs
    total_balls
    total_wickets
    status
    balls {
      ball_number
      ball_type
      run_type
      runs
      byes
      is_wicket
      wicket_type
    }
  }
}
```

### 8. All Overs
Get all overs for a specific innings.

```graphql
query GetAllOvers($matchId: String!, $inningsNumber: Int!) {
  allOvers(match_id: $matchId, innings_number: $inningsNumber) {
    over_number
    total_runs
    total_balls
    total_wickets
    status
    balls {
      ball_number
      ball_type
      run_type
      runs
      byes
      is_wicket
      wicket_type
    }
  }
}
```

### 9. Ball Details
Get all balls for a specific over.

```graphql
query GetBallDetails($matchId: String!, $inningsNumber: Int!, $overNumber: Int!) {
  ballDetails(match_id: $matchId, innings_number: $inningsNumber, over_number: $overNumber) {
    ball_number
    ball_type
    run_type
    runs
    byes
    is_wicket
    wicket_type
  }
}
```

### 10. Match Teams
Get team information for a match.

```graphql
query GetMatchTeams($matchId: String!) {
  matchTeams(match_id: $matchId) {
    id
    name
    players_count
    created_at
    updated_at
  }
}
```

### 11. Match Players
Get player information for a match.

```graphql
query GetMatchPlayers($matchId: String!) {
  matchPlayers(match_id: $matchId) {
    id
    name
    team_id
    created_at
    updated_at
  }
}
```

### 12. Player Statistics
Get player performance statistics for a match.

```graphql
query GetPlayerStatistics($matchId: String!) {
  playerStatistics(match_id: $matchId) {
    player_id
    player_name
    team_id
    runs_scored
    balls_faced
    wickets_taken
    overs_bowled
    runs_conceded
    strike_rate
    economy_rate
  }
}
```

## Sample Variables

```json
{
  "matchId": "123e4567-e89b-12d3-a456-426614174000",
  "inningsNumber": 1,
  "overNumber": 5
}
```

## Enums

### BallType
- `GOOD` - Legal delivery
- `WIDE` - Wide ball
- `NO_BALL` - No ball
- `DEAD_BALL` - Dead ball

### RunType
- `ZERO` - 0 runs
- `ONE` - 1 run
- `TWO` - 2 runs
- `THREE` - 3 runs
- `FOUR` - 4 runs
- `FIVE` - 5 runs
- `SIX` - 6 runs
- `SEVEN` - 7 runs
- `EIGHT` - 8 runs
- `NINE` - 9 runs
- `NO_BALL` - No ball runs
- `WIDE` - Wide runs
- `LEG_BYES` - Leg byes
- `WICKET` - Wicket

### TeamType
- `A` - Team A
- `B` - Team B

## Usage Examples

### Get Current Match Status
```graphql
query GetCurrentStatus($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_status
    current_innings
    current_score {
      runs
      wickets
      overs
      run_rate
    }
  }
}
```

### Get First Innings Summary
```graphql
query GetFirstInnings($matchId: String!) {
  inningsDetails(match_id: $matchId, innings_number: 1) {
    total_runs
    total_wickets
    total_overs
    status
    extras {
      total
    }
  }
}
```

### Get Recent Overs
```graphql
query GetRecentOvers($matchId: String!, $inningsNumber: Int!) {
  allOvers(match_id: $matchId, innings_number: $inningsNumber) {
    over_number
    total_runs
    total_balls
    status
  }
}
```

## Real-time Updates

Use the subscription for real-time scorecard updates:

```graphql
subscription ScorecardUpdates($matchId: String!) {
  scorecardUpdated(match_id: $matchId) {
    match_id
    current_score {
      runs
      wickets
      overs
      run_rate
    }
    current_over {
      over_number
      total_runs
      total_balls
      status
    }
  }
}
```

