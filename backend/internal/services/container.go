package services

import (
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/pkg/events"
	"spark-park-cricket-backend/pkg/websocket"
)

// Container holds all service instances
type Container struct {
	Series      *SeriesService
	Match       *MatchService
	Scorecard   *ScorecardService
	Hub         *websocket.Hub
	Broadcaster *events.EventBroadcaster
}

// NewContainer creates a new service container with all services
func NewContainer(repos *database.Repositories) *Container {
	// Create WebSocket hub
	hub := websocket.NewHub()

	// Create event broadcaster
	broadcaster := events.NewEventBroadcaster(hub)

	// Create container
	container := &Container{
		Series:      NewSeriesService(repos.Series, repos.Match),
		Match:       NewMatchService(repos.Match, repos.Series),
		Scorecard:   NewScorecardService(repos.Scorecard, repos.Match),
		Hub:         hub,
		Broadcaster: broadcaster,
	}

	return container
}
