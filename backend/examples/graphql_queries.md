# GraphQL Queries for Live Scorecard

This document provides examples of GraphQL queries for fetching live scorecard data with field-level selection.

## Basic Live Scorecard Query

Fetch only the essential live scorecard information:

```graphql
query GetLiveScorecard($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_id
    match_number
    series_name
    team_a
    team_b
    current_innings
    match_status
    current_score {
      runs
      wickets
      overs
      balls
      run_rate
    }
  }
}
```

Variables:
```json
{
  "matchId": "your-match-id-here"
}
```

## Detailed Scorecard Query

Fetch comprehensive scorecard information:

```graphql
query GetDetailedScorecard($matchId: String!) {
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

## Minimal Live Score Query

Fetch only the current score for real-time updates:

```graphql
query GetCurrentScore($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_id
    current_score {
      runs
      wickets
      overs
      run_rate
    }
  }
}
```

## Current Over Query

Fetch only the current over information:

```graphql
query GetCurrentOver($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_id
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
        is_wicket
      }
    }
  }
}
```

## Match Summary Query

Fetch only match summary information:

```graphql
query GetMatchSummary($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_id
    match_number
    series_name
    team_a
    team_b
    total_overs
    toss_winner
    current_innings
    match_status
  }
}
```

## Innings Summary Query

Fetch only innings summaries without ball details:

```graphql
query GetInningsSummary($matchId: String!) {
  liveScorecard(match_id: $matchId) {
    match_id
    innings {
      innings_number
      batting_team
      total_runs
      total_wickets
      total_overs
      status
      extras {
        total
      }
    }
  }
}
```

## WebSocket Subscription

For real-time updates, you can use WebSocket connections to the `/api/v1/ws/match/{match_id}` endpoint. The WebSocket will receive updates in the following format:

```json
{
  "type": "scorecard_update",
  "data": {
    "match_id": "match-id",
    "current_score": {
      "runs": 150,
      "wickets": 3,
      "overs": 25.3,
      "run_rate": 5.88
    },
    "current_over": {
      "over_number": 26,
      "total_runs": 8,
      "total_balls": 3,
      "status": "in_progress"
    }
  }
}
```

## Usage Examples

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query GetLiveScorecard($matchId: String!) { liveScorecard(match_id: $matchId) { match_id current_score { runs wickets overs run_rate } } }",
    "variables": {
      "matchId": "your-match-id-here"
    }
  }'
```

### JavaScript Example

```javascript
const query = `
  query GetLiveScorecard($matchId: String!) {
    liveScorecard(match_id: $matchId) {
      match_id
      current_score {
        runs
        wickets
        overs
        run_rate
      }
    }
  }
`;

const variables = {
  matchId: "your-match-id-here"
};

fetch('http://localhost:8080/api/v1/graphql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    query,
    variables
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

## Benefits of GraphQL

1. **Field Selection**: Only fetch the fields you need, reducing payload size
2. **Single Request**: Get all required data in one request
3. **Type Safety**: Strong typing with GraphQL schema
4. **Real-time Updates**: WebSocket integration for live updates
5. **Flexible Queries**: Adapt queries based on UI requirements
