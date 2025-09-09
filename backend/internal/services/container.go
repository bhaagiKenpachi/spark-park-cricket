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
	Team        *TeamService
	Player      *PlayerService
	Scoreboard  *RealtimeScoreboardService
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
		Series:      NewSeriesService(repos.Series),
		Match:       NewMatchService(repos.Match, repos.Series, repos.Team),
		Team:        NewTeamService(repos.Team, repos.Player),
		Player:      NewPlayerService(repos.Player, repos.Team),
		Hub:         hub,
		Broadcaster: broadcaster,
	}

	// Create real-time scoreboard service with broadcaster
	container.Scoreboard = NewRealtimeScoreboardService(
		repos.Scoreboard, repos.Over, repos.Ball, repos.Match, repos.Team, repos.Player,
		broadcaster,
	)

	return container
}
