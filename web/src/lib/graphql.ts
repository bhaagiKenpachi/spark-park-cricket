import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  gql,
} from '@apollo/client';

// GraphQL endpoint
const GRAPHQL_ENDPOINT =
  process.env.NEXT_PUBLIC_GRAPHQL_URL || 'http://localhost:8080/api/v1/graphql';

// Create HTTP link
const httpLink = createHttpLink({
  uri: GRAPHQL_ENDPOINT,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Create Apollo Client
export const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
    },
    query: {
      fetchPolicy: 'network-only',
    },
  },
});

// Optimized GraphQL Queries for minimal data fetching

// 1. Fetch only innings score summary (runs, wickets, overs) - for scoreboard display
export const GET_INNINGS_SCORE_SUMMARY = gql`
  query GetInningsScoreSummary($matchId: String!, $inningsNumber: Int!) {
    inningsScore(match_id: $matchId, innings_number: $inningsNumber) {
      innings_number
      batting_team
      total_runs
      total_wickets
      total_overs
      total_balls
      status
      extras {
        total
      }
    }
  }
`;

// 2. Fetch only latest over details - for current over display
export const GET_LATEST_OVER_ONLY = gql`
  query GetLatestOverOnly($matchId: String!, $inningsNumber: Int!) {
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
`;

// 3. Fetch all overs details only when needed - for expanded view
export const GET_ALL_OVERS_DETAILS = gql`
  query GetAllOversDetails($matchId: String!, $inningsNumber: Int!) {
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
`;

// Keep the original detailed queries for backward compatibility
export const GET_INNINGS_SCORE = gql`
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
`;

export const GET_INNINGS_DETAILS = gql`
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
`;

export const GET_LIVE_SCORECARD = gql`
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
`;

export const GET_LATEST_OVER = gql`
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
`;

// TypeScript types for GraphQL responses

// Optimized response types for minimal data fetching
export interface InningsScoreSummaryResponse {
  inningsScore: {
    innings_number: number;
    batting_team: string;
    total_runs: number;
    total_wickets: number;
    total_overs: number;
    total_balls: number;
    status: string;
    extras: {
      total: number;
    };
  };
}

export interface LatestOverOnlyResponse {
  latestOver: {
    over_number: number;
    total_runs: number;
    total_balls: number;
    total_wickets: number;
    status: string;
    balls: Array<{
      ball_number: number;
      ball_type: string;
      run_type: string;
      runs: number;
      byes: number;
      is_wicket: boolean;
      wicket_type?: string;
    }>;
  };
}

export interface AllOversDetailsResponse {
  allOvers: Array<{
    over_number: number;
    total_runs: number;
    total_balls: number;
    total_wickets: number;
    status: string;
    balls: Array<{
      ball_number: number;
      ball_type: string;
      run_type: string;
      runs: number;
      byes: number;
      is_wicket: boolean;
      wicket_type?: string;
    }>;
  }>;
}

// Original detailed response types (for backward compatibility)
export interface InningsScoreResponse {
  inningsScore: {
    innings_number: number;
    batting_team: string;
    total_runs: number;
    total_wickets: number;
    total_overs: number;
    total_balls: number;
    status: string;
    extras: {
      byes: number;
      leg_byes: number;
      wides: number;
      no_balls: number;
      total: number;
    };
  };
}

export interface InningsDetailsResponse {
  inningsDetails: {
    innings_number: number;
    batting_team: string;
    total_runs: number;
    total_wickets: number;
    total_overs: number;
    total_balls: number;
    status: string;
    extras: {
      byes: number;
      leg_byes: number;
      wides: number;
      no_balls: number;
      total: number;
    };
    overs: Array<{
      over_number: number;
      total_runs: number;
      total_balls: number;
      total_wickets: number;
      status: string;
      balls: Array<{
        ball_number: number;
        ball_type: string;
        run_type: string;
        runs: number;
        byes: number;
        is_wicket: boolean;
        wicket_type?: string;
      }>;
    }>;
  };
}

export interface LiveScorecardResponse {
  liveScorecard: {
    match_id: string;
    match_number: number;
    series_name: string;
    team_a: string;
    team_b: string;
    total_overs: number;
    toss_winner: string;
    toss_type: string;
    current_innings: number;
    match_status: string;
    current_score: {
      runs: number;
      wickets: number;
      overs: number;
      balls: number;
      run_rate: number;
    };
    current_over?: {
      over_number: number;
      total_runs: number;
      total_balls: number;
      total_wickets: number;
      status: string;
      balls: Array<{
        ball_number: number;
        ball_type: string;
        run_type: string;
        runs: number;
        byes: number;
        is_wicket: boolean;
        wicket_type?: string;
      }>;
    };
    innings: Array<{
      innings_number: number;
      batting_team: string;
      total_runs: number;
      total_wickets: number;
      total_overs: number;
      total_balls: number;
      status: string;
      extras: {
        byes: number;
        leg_byes: number;
        wides: number;
        no_balls: number;
        total: number;
      };
      overs: Array<{
        over_number: number;
        total_runs: number;
        total_balls: number;
        total_wickets: number;
        status: string;
        balls: Array<{
          ball_number: number;
          ball_type: string;
          run_type: string;
          runs: number;
          byes: number;
          is_wicket: boolean;
          wicket_type?: string;
        }>;
      }>;
    }>;
  };
}

export interface LatestOverResponse {
  latestOver: {
    over_number: number;
    total_runs: number;
    total_balls: number;
    total_wickets: number;
    status: string;
    balls: Array<{
      ball_number: number;
      ball_type: string;
      run_type: string;
      runs: number;
      byes: number;
      is_wicket: boolean;
      wicket_type?: string;
    }>;
  };
}

// Simplified types for Redux state and components
export interface InningsSummary {
  innings_number: number;
  batting_team: string;
  total_runs: number;
  total_wickets: number;
  total_overs: number;
  total_balls: number;
  status: string;
  extras: {
    byes: number;
    leg_byes: number;
    wides: number;
    no_balls: number;
    total: number;
  };
  overs: OverSummary[];
}

export interface OverSummary {
  over_number: number;
  total_runs: number;
  total_balls: number;
  total_wickets: number;
  status: string;
  balls: BallSummary[];
}

export interface BallSummary {
  ball_number: number;
  ball_type: string;
  run_type: string;
  runs: number;
  byes: number;
  is_wicket: boolean;
  wicket_type?: string;
}
