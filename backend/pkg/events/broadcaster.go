package events

import (
	"context"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/websocket"
	"time"
)

// EventBroadcaster handles broadcasting events to WebSocket clients
type EventBroadcaster struct {
	hub *websocket.Hub
}

// NewEventBroadcaster creates a new event broadcaster
func NewEventBroadcaster(hub *websocket.Hub) *EventBroadcaster {
	return &EventBroadcaster{
		hub: hub,
	}
}

// BroadcastBallEvent broadcasts a ball event to all clients watching the match
func (eb *EventBroadcaster) BroadcastBallEvent(ctx context.Context, matchID string, ballEvent *models.BallEvent, scoreboard *models.LiveScoreboard) {
	message := websocket.Message{
		Type:   "ball_event",
		RoomID: matchID,
		Data: map[string]interface{}{
			"ball_event": map[string]interface{}{
				"ball_type": ballEvent.BallType,
				"run_type":  ballEvent.RunType,
				"is_wicket": ballEvent.IsWicket,
				"timestamp": time.Now().Unix(),
			},
			"scoreboard": map[string]interface{}{
				"score":        scoreboard.Score,
				"wickets":      scoreboard.Wickets,
				"overs":        scoreboard.Overs,
				"balls":        scoreboard.Balls,
				"batting_team": scoreboard.BattingTeam,
				"updated_at":   scoreboard.UpdatedAt.Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted ball event for match %s: %s, %s", matchID, ballEvent.BallType, ballEvent.RunType)
}

// BroadcastScoreUpdate broadcasts a score update to all clients watching the match
func (eb *EventBroadcaster) BroadcastScoreUpdate(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) {
	message := websocket.Message{
		Type:   "score_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"scoreboard": map[string]interface{}{
				"score":        scoreboard.Score,
				"wickets":      scoreboard.Wickets,
				"overs":        scoreboard.Overs,
				"balls":        scoreboard.Balls,
				"batting_team": scoreboard.BattingTeam,
				"updated_at":   scoreboard.UpdatedAt.Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted score update for match %s: %d/%d in %.1f overs", matchID, scoreboard.Score, scoreboard.Wickets, scoreboard.Overs)
}

// BroadcastWicketUpdate broadcasts a wicket update to all clients watching the match
func (eb *EventBroadcaster) BroadcastWicketUpdate(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) {
	message := websocket.Message{
		Type:   "wicket_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"wicket": map[string]interface{}{
				"wickets":   scoreboard.Wickets,
				"timestamp": time.Now().Unix(),
			},
			"scoreboard": map[string]interface{}{
				"score":        scoreboard.Score,
				"wickets":      scoreboard.Wickets,
				"overs":        scoreboard.Overs,
				"balls":        scoreboard.Balls,
				"batting_team": scoreboard.BattingTeam,
				"updated_at":   scoreboard.UpdatedAt.Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted wicket update for match %s: %d wickets", matchID, scoreboard.Wickets)
}

// BroadcastOverCompletion broadcasts an over completion to all clients watching the match
func (eb *EventBroadcaster) BroadcastOverCompletion(ctx context.Context, matchID string, over *models.Over, scoreboard *models.LiveScoreboard) {
	message := websocket.Message{
		Type:   "over_completion",
		RoomID: matchID,
		Data: map[string]interface{}{
			"over": map[string]interface{}{
				"over_number":  over.OverNumber,
				"total_runs":   over.TotalRuns,
				"total_balls":  over.TotalBalls,
				"batting_team": over.BattingTeam,
				"completed_at": time.Now().Unix(),
			},
			"scoreboard": map[string]interface{}{
				"score":        scoreboard.Score,
				"wickets":      scoreboard.Wickets,
				"overs":        scoreboard.Overs,
				"balls":        scoreboard.Balls,
				"batting_team": scoreboard.BattingTeam,
				"updated_at":   scoreboard.UpdatedAt.Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted over completion for match %s: Over %d completed with %d runs", matchID, over.OverNumber, over.TotalRuns)
}

// BroadcastMatchStatusUpdate broadcasts a match status update to all clients watching the match
func (eb *EventBroadcaster) BroadcastMatchStatusUpdate(ctx context.Context, matchID string, match *models.Match) {
	message := websocket.Message{
		Type:   "match_status_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"match": map[string]interface{}{
				"id":                  match.ID,
				"status":              match.Status,
				"match_number":        match.MatchNumber,
				"team_a_player_count": match.TeamAPlayerCount,
				"team_b_player_count": match.TeamBPlayerCount,
				"updated_at":          match.UpdatedAt.Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted match status update for match %s: Status changed to %s", matchID, match.Status)
}

// BroadcastMatchStart broadcasts a match start event to all clients watching the match
func (eb *EventBroadcaster) BroadcastMatchStart(ctx context.Context, matchID string, match *models.Match) {
	message := websocket.Message{
		Type:   "match_start",
		RoomID: matchID,
		Data: map[string]interface{}{
			"match": map[string]interface{}{
				"id":                  match.ID,
				"status":              match.Status,
				"match_number":        match.MatchNumber,
				"team_a_player_count": match.TeamAPlayerCount,
				"team_b_player_count": match.TeamBPlayerCount,
				"started_at":          time.Now().Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted match start for match %s", matchID)
}

// BroadcastMatchEnd broadcasts a match end event to all clients watching the match
func (eb *EventBroadcaster) BroadcastMatchEnd(ctx context.Context, matchID string, match *models.Match, finalScoreboard *models.LiveScoreboard) {
	message := websocket.Message{
		Type:   "match_end",
		RoomID: matchID,
		Data: map[string]interface{}{
			"match": map[string]interface{}{
				"id":                  match.ID,
				"status":              match.Status,
				"match_number":        match.MatchNumber,
				"team_a_player_count": match.TeamAPlayerCount,
				"team_b_player_count": match.TeamBPlayerCount,
				"ended_at":            time.Now().Unix(),
			},
			"final_scoreboard": map[string]interface{}{
				"score":        finalScoreboard.Score,
				"wickets":      finalScoreboard.Wickets,
				"overs":        finalScoreboard.Overs,
				"balls":        finalScoreboard.Balls,
				"batting_team": finalScoreboard.BattingTeam,
				"updated_at":   finalScoreboard.UpdatedAt.Unix(),
			},
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted match end for match %s: Final score %d/%d", matchID, finalScoreboard.Score, finalScoreboard.Wickets)
}

// BroadcastCustomMessage broadcasts a custom message to all clients in a room
func (eb *EventBroadcaster) BroadcastCustomMessage(ctx context.Context, matchID string, messageType string, data interface{}) {
	message := websocket.Message{
		Type:   messageType,
		RoomID: matchID,
		Data: map[string]interface{}{
			"data":      data,
			"timestamp": time.Now().Unix(),
		},
	}

	eb.hub.BroadcastToRoom(matchID, message)
	log.Printf("Broadcasted custom message %s for match %s", messageType, matchID)
}
