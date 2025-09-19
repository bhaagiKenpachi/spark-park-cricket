import {
  apolloClient,
  GET_INNINGS_SCORE_SUMMARY,
  GET_LATEST_OVER_ONLY,
  GET_ALL_OVERS_DETAILS,
  GET_INNINGS_SCORE,
  GET_INNINGS_DETAILS,
  GET_LIVE_SCORECARD,
  GET_LATEST_OVER,
  InningsScoreSummaryResponse,
  LatestOverOnlyResponse,
  AllOversDetailsResponse,
  InningsScoreResponse,
  InningsDetailsResponse,
  LiveScorecardResponse,
  LatestOverResponse,
} from '@/lib/graphql';

export class GraphQLService {
  // ===== OPTIMIZED METHODS FOR MINIMAL DATA FETCHING =====

  /**
   * Fetch only innings score summary (runs, wickets, overs) - for scoreboard display
   * This is the most efficient query for basic scoreboard information
   */
  async getInningsScoreSummary(matchId: string, inningsNumber: number) {
    try {
      const { data } = await apolloClient.query<InningsScoreSummaryResponse>({
        query: GET_INNINGS_SCORE_SUMMARY,
        variables: {
          matchId,
          inningsNumber,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.inningsScore,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch innings score summary',
      };
    }
  }

  /**
   * Fetch only latest over details - for current over display
   * This fetches only the current over with ball details
   */
  async getLatestOverOnly(matchId: string, inningsNumber: number) {
    try {
      const { data } = await apolloClient.query<LatestOverOnlyResponse>({
        query: GET_LATEST_OVER_ONLY,
        variables: {
          matchId,
          inningsNumber,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.latestOver,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch latest over',
      };
    }
  }

  /**
   * Fetch all overs details only when needed - for expanded view
   * This is called only when user clicks to expand all overs
   */
  async getAllOversDetails(matchId: string, inningsNumber: number) {
    try {
      const { data } = await apolloClient.query<AllOversDetailsResponse>({
        query: GET_ALL_OVERS_DETAILS,
        variables: {
          matchId,
          inningsNumber,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.allOvers,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch all overs details',
      };
    }
  }

  // ===== ORIGINAL METHODS (for backward compatibility) =====

  /**
   * Fetch innings score data using GraphQL
   */
  async getInningsScore(matchId: string, inningsNumber: number) {
    try {
      const { data } = await apolloClient.query<InningsScoreResponse>({
        query: GET_INNINGS_SCORE,
        variables: {
          matchId,
          inningsNumber,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.inningsScore,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch innings score',
      };
    }
  }

  /**
   * Fetch detailed innings data using GraphQL
   */
  async getInningsDetails(matchId: string, inningsNumber: number) {
    try {
      const { data } = await apolloClient.query<InningsDetailsResponse>({
        query: GET_INNINGS_DETAILS,
        variables: {
          matchId,
          inningsNumber,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.inningsDetails,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch innings details',
      };
    }
  }

  /**
   * Fetch live scorecard data using GraphQL
   */
  async getLiveScorecard(matchId: string) {
    try {
      const { data } = await apolloClient.query<LiveScorecardResponse>({
        query: GET_LIVE_SCORECARD,
        variables: {
          matchId,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.liveScorecard,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch live scorecard',
      };
    }
  }

  /**
   * Fetch latest over data using GraphQL
   */
  async getLatestOver(matchId: string, inningsNumber: number) {
    try {
      const { data } = await apolloClient.query<LatestOverResponse>({
        query: GET_LATEST_OVER,
        variables: {
          matchId,
          inningsNumber,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.latestOver,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch latest over',
      };
    }
  }

  /**
   * Fetch updated scorecard data after ball addition
   * This method fetches the complete live scorecard to get the latest state
   */
  async getUpdatedScorecardAfterBall(matchId: string) {
    try {
      const { data } = await apolloClient.query<LiveScorecardResponse>({
        query: GET_LIVE_SCORECARD,
        variables: {
          matchId,
        },
        fetchPolicy: 'network-only', // Always fetch fresh data
      });

      return {
        success: true,
        data: data?.liveScorecard,
      };
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error
            ? error.message
            : 'Failed to fetch updated scorecard',
      };
    }
  }
}

export const graphqlService = new GraphQLService();
export default graphqlService;
