package services

import (
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/graphql"
	"spark-park-cricket-backend/internal/interfaces"
	"spark-park-cricket-backend/pkg/events"
	"spark-park-cricket-backend/pkg/websocket"
)

// Container holds all service instances
type Container struct {
	Series           *SeriesService
	Match            *MatchService
	Scorecard        interfaces.ScorecardServiceInterface
	Hub              *websocket.Hub
	Broadcaster      *events.EventBroadcaster
	GraphQLWebSocket *graphql.GraphQLWebSocketService
}

// NewContainer creates a new service container with all services
func NewContainer(repos *database.Repositories) *Container {
	// Create WebSocket hub
	hub := websocket.NewHub()

	// Create event broadcaster
	broadcaster := events.NewEventBroadcaster(hub)

	// Create base scorecard service
	baseScorecardService := NewScorecardService(repos.Scorecard, repos.Match)

	// Create GraphQL WebSocket service
	graphqlWebSocketService := graphql.NewGraphQLWebSocketService(baseScorecardService, hub)

	// Create GraphQL-integrated scorecard service
	scorecardServiceWithGraphQL := NewScorecardServiceWithGraphQL(repos.Scorecard, repos.Match, hub)

	// Create container
	container := &Container{
		Series:           NewSeriesService(repos.Series),
		Match:            NewMatchService(repos.Match, repos.Series),
		Scorecard:        scorecardServiceWithGraphQL,
		Hub:              hub,
		Broadcaster:      broadcaster,
		GraphQLWebSocket: graphqlWebSocketService,
	}

	return container
}
