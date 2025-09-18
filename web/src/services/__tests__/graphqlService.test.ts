/* eslint-disable @typescript-eslint/no-require-imports */
// Simple test for GraphQL service functionality
// This test verifies that the service methods exist and can be called

// Mock Apollo Client to prevent GraphQL errors in tests
jest.mock('@apollo/client', () => ({
  ApolloClient: jest.fn(),
  InMemoryCache: jest.fn(),
  createHttpLink: jest.fn(),
  gql: jest.fn(query => query),
}));

// Mock the GraphQL client
jest.mock('../../lib/graphql', () => ({
  apolloClient: {
    query: jest
      .fn()
      .mockRejectedValue(new Error('GraphQL client not configured for tests')),
  },
  GET_INNINGS_SCORE_SUMMARY: 'GET_INNINGS_SCORE_SUMMARY',
  GET_LATEST_OVER_ONLY: 'GET_LATEST_OVER_ONLY',
  GET_ALL_OVERS_DETAILS: 'GET_ALL_OVERS_DETAILS',
  GET_INNINGS_SCORE: 'GET_INNINGS_SCORE',
  GET_INNINGS_DETAILS: 'GET_INNINGS_DETAILS',
  GET_LIVE_SCORECARD: 'GET_LIVE_SCORECARD',
  GET_LATEST_OVER: 'GET_LATEST_OVER',
}));

describe('GraphQLService', () => {
  it('should have all required methods', () => {
    // Import the service to verify it exists
    const { graphqlService } = require('../graphqlService');

    expect(graphqlService).toBeDefined();
    expect(typeof graphqlService.getInningsScoreSummary).toBe('function');
    expect(typeof graphqlService.getLatestOverOnly).toBe('function');
    expect(typeof graphqlService.getAllOversDetails).toBe('function');
    expect(typeof graphqlService.getLiveScorecard).toBe('function');
    expect(typeof graphqlService.getInningsDetails).toBe('function');
  });

  it('should handle method calls without errors', async () => {
    const { graphqlService } = require('../graphqlService');

    // Test that methods can be called and return proper error responses
    const result1 = await graphqlService.getInningsScoreSummary(
      'test-match',
      1
    );
    expect(result1.success).toBe(false);
    expect(result1.error).toBeDefined();

    const result2 = await graphqlService.getLatestOverOnly('test-match', 1);
    expect(result2.success).toBe(false);
    expect(result2.error).toBeDefined();

    const result3 = await graphqlService.getAllOversDetails('test-match', 1);
    expect(result3.success).toBe(false);
    expect(result3.error).toBeDefined();

    const result4 = await graphqlService.getLiveScorecard('test-match');
    expect(result4.success).toBe(false);
    expect(result4.error).toBeDefined();

    const result5 = await graphqlService.getInningsDetails('test-match', 1);
    expect(result5.success).toBe(false);
    expect(result5.error).toBeDefined();
  });
});
