package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/monitoring"
	"time"

	"github.com/go-chi/chi/v5"
)

// SSEHandler handles Server-Sent Events for real-time ball updates
type SSEHandler struct {
	redisClient *cache.RedisClient
	metrics     *monitoring.Metrics
	config      *config.Config
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(redisClient *cache.RedisClient, metrics *monitoring.Metrics, cfg *config.Config) *SSEHandler {
	return &SSEHandler{
		redisClient: redisClient,
		metrics:     metrics,
		config:      cfg,
	}
}

// StreamBallEvents streams ball events for a match using Server-Sent Events
func (h *SSEHandler) StreamBallEvents(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	// Check if Redis client is available
	if h.redisClient == nil {
		log.Printf("⚠️  Redis client not available for SSE streaming")
		if h.metrics != nil {
			h.metrics.RecordSSEError(matchID, "redis_unavailable")
		}
		http.Error(w, "Streaming service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// CORS headers are handled by the CORS middleware, don't override here

	// Get flusher for streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("⚠️  Streaming not supported by response writer")
		if h.metrics != nil {
			h.metrics.RecordSSEError(matchID, "flusher_not_supported")
		}
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Record connection metrics
	connectionStart := time.Now()
	endpoint := "balls"
	if h.metrics != nil {
		h.metrics.RecordSSEConnection(matchID, endpoint)
	}

	log.Printf("📡 SSE: Client connected for match %s", matchID)

	// Send initial connection message
	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: {\"message\":\"Connected to ball events stream for match %s\",\"timestamp\":\"%s\"}\n\n",
		matchID, time.Now().Format(time.RFC3339))
	flusher.Flush()

	// Record connection event sent
	if h.metrics != nil {
		h.metrics.RecordSSEEventSent(matchID, "connected")
	}

	// Get stream key for this match
	streamKey := h.redisClient.GetStreamKey(matchID)

	// Start reading from the latest event
	lastID := "$" // $ means start from the latest event

	// Create context with cancellation and timeout
	// Use background context to avoid request context timeout issues
	ctx, cancel := context.WithTimeout(context.Background(), h.config.SSEConnectionTimeout)
	defer cancel()

	// Create ticker for periodic heartbeat using config
	heartbeatTicker := time.NewTicker(h.config.SSEHeartbeatInterval)
	defer heartbeatTicker.Stop()

	// Stream events
	disconnectReason := "client_disconnect"
	for {
		select {
		case <-ctx.Done():
			// Check if this is a timeout or client disconnect
			if ctx.Err() == context.DeadlineExceeded {
				disconnectReason = "timeout"
				log.Printf("⏰ SSE TIMEOUT: Connection timeout after %v for match %s", h.config.SSEConnectionTimeout, matchID)

				// Send timeout notification to client before closing
				fmt.Fprintf(w, "event: timeout\n")
				fmt.Fprintf(w, "data: {\"message\":\"Connection timeout\",\"timeout_duration\":\"%v\",\"timestamp\":\"%s\"}\n\n",
					h.config.SSEConnectionTimeout, time.Now().Format(time.RFC3339))
				flusher.Flush()
			} else {
				log.Printf("📡 SSE: Client disconnected for match %s", matchID)
			}

			if h.metrics != nil {
				h.metrics.RecordSSEDisconnection(matchID, endpoint, disconnectReason, time.Since(connectionStart))
			}
			return

		case <-heartbeatTicker.C:
			// Send heartbeat to keep connection alive
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()

		default:
			// Read from stream with blocking
			streamReadStart := time.Now()
			messages, err := h.redisClient.ReadFromStream(ctx, streamKey, lastID, 10, 1*time.Second)
			if h.metrics != nil {
				h.metrics.RecordSSEStreamRead(matchID, time.Since(streamReadStart))
			}

			if err != nil {
				if ctx.Err() != nil {
					// Context cancelled, client disconnected
					disconnectReason = "context_cancelled"
					if h.metrics != nil {
						h.metrics.RecordSSEDisconnection(matchID, endpoint, disconnectReason, time.Since(connectionStart))
					}
					return
				}
				log.Printf("⚠️  Error reading from stream: %v", err)
				if h.metrics != nil {
					h.metrics.RecordSSEError(matchID, "stream_read_error")
				}
				time.Sleep(1 * time.Second)
				continue
			}

			// Process each message
			for _, msg := range messages {
				// Update lastID to the current message ID
				lastID = msg.ID

				// Convert message values to JSON
				eventData := make(map[string]interface{})
				for k, v := range msg.Values {
					eventData[k] = v
				}

				// Add message ID and timestamp
				eventData["stream_id"] = msg.ID

				// Marshal to JSON
				jsonData, err := json.Marshal(eventData)
				if err != nil {
					log.Printf("⚠️  Error marshaling event data: %v", err)
					if h.metrics != nil {
						h.metrics.RecordSSEError(matchID, "marshal_error")
						h.metrics.RecordSSEEventProcessed(matchID, "ball_added", "error")
					}
					continue
				}

				// Record event processed successfully
				if h.metrics != nil {
					h.metrics.RecordSSEEventProcessed(matchID, "ball_added", "success")
				}

				// Send event to client
				fmt.Fprintf(w, "event: ball_added\n")
				fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
				flusher.Flush()

				// Record event sent
				if h.metrics != nil {
					h.metrics.RecordSSEEventSent(matchID, "ball_added")
				}

				log.Printf("📡 SSE: Sent ball event to client for match %s (ID: %s)", matchID, msg.ID)
			}
		}
	}
}

// StreamMatchEvents streams all match events (not just balls) for a match
func (h *SSEHandler) StreamMatchEvents(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	// Check if Redis client is available
	if h.redisClient == nil {
		log.Printf("⚠️  Redis client not available for SSE streaming")
		if h.metrics != nil {
			h.metrics.RecordSSEError(matchID, "redis_unavailable")
		}
		http.Error(w, "Streaming service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// CORS headers are handled by the CORS middleware, don't override here

	// Get flusher for streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("⚠️  Streaming not supported by response writer")
		if h.metrics != nil {
			h.metrics.RecordSSEError(matchID, "flusher_not_supported")
		}
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Record connection metrics
	connectionStart := time.Now()
	endpoint := "events"
	if h.metrics != nil {
		h.metrics.RecordSSEConnection(matchID, endpoint)
	}

	log.Printf("📡 SSE: Client connected for all match events %s", matchID)

	// Send initial connection message
	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: {\"message\":\"Connected to match events stream for match %s\",\"timestamp\":\"%s\"}\n\n",
		matchID, time.Now().Format(time.RFC3339))
	flusher.Flush()

	// Record connection event sent
	if h.metrics != nil {
		h.metrics.RecordSSEEventSent(matchID, "connected")
	}

	// Get stream key for this match
	streamKey := h.redisClient.GetStreamKey(matchID)

	// Start reading from the latest event
	lastID := "$"

	// Create context with cancellation and timeout
	// Use background context to avoid request context timeout issues
	ctx, cancel := context.WithTimeout(context.Background(), h.config.SSEConnectionTimeout)
	defer cancel()

	// Create ticker for periodic heartbeat using config
	heartbeatTicker := time.NewTicker(h.config.SSEHeartbeatInterval)
	defer heartbeatTicker.Stop()

	// Stream events
	disconnectReason := "client_disconnect"
	for {
		select {
		case <-ctx.Done():
			// Check if this is a timeout or client disconnect
			if ctx.Err() == context.DeadlineExceeded {
				disconnectReason = "timeout"
				log.Printf("⏰ SSE TIMEOUT: Connection timeout after %v for match %s", h.config.SSEConnectionTimeout, matchID)

				// Send timeout notification to client before closing
				fmt.Fprintf(w, "event: timeout\n")
				fmt.Fprintf(w, "data: {\"message\":\"Connection timeout\",\"timeout_duration\":\"%v\",\"timestamp\":\"%s\"}\n\n",
					h.config.SSEConnectionTimeout, time.Now().Format(time.RFC3339))
				flusher.Flush()
			} else {
				log.Printf("📡 SSE: Client disconnected for match %s", matchID)
			}

			if h.metrics != nil {
				h.metrics.RecordSSEDisconnection(matchID, endpoint, disconnectReason, time.Since(connectionStart))
			}
			return

		case <-heartbeatTicker.C:
			// Send heartbeat to keep connection alive
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()

		default:
			// Read from stream with blocking
			streamReadStart := time.Now()
			messages, err := h.redisClient.ReadFromStream(ctx, streamKey, lastID, 10, 1*time.Second)
			if h.metrics != nil {
				h.metrics.RecordSSEStreamRead(matchID, time.Since(streamReadStart))
			}

			if err != nil {
				if ctx.Err() != nil {
					// Context cancelled, client disconnected
					disconnectReason = "context_cancelled"
					if h.metrics != nil {
						h.metrics.RecordSSEDisconnection(matchID, endpoint, disconnectReason, time.Since(connectionStart))
					}
					return
				}
				log.Printf("⚠️  Error reading from stream: %v", err)
				if h.metrics != nil {
					h.metrics.RecordSSEError(matchID, "stream_read_error")
				}
				time.Sleep(1 * time.Second)
				continue
			}

			// Process each message
			for _, msg := range messages {
				// Update lastID to the current message ID
				lastID = msg.ID

				// Get event type from message
				eventType, ok := msg.Values["event_type"].(string)
				if !ok {
					eventType = "match_event"
				}

				// Convert message values to JSON
				eventData := make(map[string]interface{})
				for k, v := range msg.Values {
					eventData[k] = v
				}

				// Add message ID
				eventData["stream_id"] = msg.ID

				// Marshal to JSON
				jsonData, err := json.Marshal(eventData)
				if err != nil {
					log.Printf("⚠️  Error marshaling event data: %v", err)
					if h.metrics != nil {
						h.metrics.RecordSSEError(matchID, "marshal_error")
						h.metrics.RecordSSEEventProcessed(matchID, eventType, "error")
					}
					continue
				}

				// Record event processed successfully
				if h.metrics != nil {
					h.metrics.RecordSSEEventProcessed(matchID, eventType, "success")
				}

				// Send event to client with dynamic event type
				fmt.Fprintf(w, "event: %s\n", eventType)
				fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
				flusher.Flush()

				// Record event sent
				if h.metrics != nil {
					h.metrics.RecordSSEEventSent(matchID, eventType)
				}

				log.Printf("📡 SSE: Sent %s event to client for match %s (ID: %s)", eventType, matchID, msg.ID)
			}
		}
	}
}
