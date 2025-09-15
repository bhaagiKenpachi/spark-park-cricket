package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// MatchHandler handles HTTP requests for match operations
type MatchHandler struct {
	service services.MatchServiceInterface
}

// NewMatchHandler creates a new match handler
func NewMatchHandler(service services.MatchServiceInterface) *MatchHandler {
	return &MatchHandler{
		service: service,
	}
}

// ListMatches handles GET /api/v1/matches
func (h *MatchHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: ListMatches handler called")

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	seriesID := r.URL.Query().Get("series_id")
	status := r.URL.Query().Get("status")

	// Set default limit
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Set default offset
	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Create filters
	filters := &models.MatchFilters{
		Limit:  limit,
		Offset: offset,
	}

	if seriesID != "" {
		filters.SeriesID = &seriesID
	}

	if status != "" {
		matchStatus := models.MatchStatus(status)
		filters.Status = &matchStatus
	}

	log.Printf("DEBUG: Created filters: %+v", filters)

	// Get matches from service
	matches, err := h.service.ListMatches(r.Context(), filters)
	if err != nil {
		log.Printf("DEBUG: service.ListMatches failed: %v", err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("DEBUG: Retrieved %d matches", len(matches))
	utils.WriteSuccess(w, matches)
}

// CreateMatch handles POST /api/v1/matches
func (h *MatchHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: CreateMatch handler called")

	var req models.CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("DEBUG: Failed to decode request body: %v", err)
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	log.Printf("DEBUG: Decoded request: %+v", req)

	// Create match
	log.Printf("DEBUG: Calling service.CreateMatch")
	match, err := h.service.CreateMatch(r.Context(), &req)
	if err != nil {
		log.Printf("DEBUG: service.CreateMatch failed: %v", err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("DEBUG: Match created successfully, writing response: %+v", match)
	utils.WriteCreated(w, match)
}

// GetMatch handles GET /api/v1/matches/{id}
func (h *MatchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: GetMatch handler called")

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Printf("DEBUG: Match ID is required")
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	log.Printf("DEBUG: Getting match with ID: %s", id)

	// Get match
	match, err := h.service.GetMatch(r.Context(), id)
	if err != nil {
		log.Printf("DEBUG: service.GetMatch failed: %v", err)
		utils.WriteNotFound(w, "Match not found")
		return
	}

	log.Printf("DEBUG: Match retrieved successfully: %+v", match)
	utils.WriteSuccess(w, match)
}

// UpdateMatch handles PUT /api/v1/matches/{id}
func (h *MatchHandler) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: UpdateMatch handler called")

	id := chi.URLParam(r, "id")

	// Fallback: Extract ID manually if chi.URLParam fails (e.g., in E2E tests)
	if id == "" {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/matches/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			id = parts[0]
		}
	}

	if id == "" {
		log.Printf("DEBUG: Match ID is required")
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	var req models.UpdateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("DEBUG: Failed to decode request body: %v", err)
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	log.Printf("DEBUG: Updating match %s with request: %+v", id, req)

	// Update match
	match, err := h.service.UpdateMatch(r.Context(), id, &req)
	if err != nil {
		log.Printf("DEBUG: service.UpdateMatch failed: %v", err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("DEBUG: Match updated successfully: %+v", match)
	utils.WriteSuccess(w, match)
}

// DeleteMatch handles DELETE /api/v1/matches/{id}
func (h *MatchHandler) DeleteMatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: DeleteMatch handler called")

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Printf("DEBUG: Match ID is required")
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	log.Printf("DEBUG: Deleting match with ID: %s", id)

	// Delete match
	err := h.service.DeleteMatch(r.Context(), id)
	if err != nil {
		log.Printf("DEBUG: service.DeleteMatch failed: %v", err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("DEBUG: Match deleted successfully")
	utils.WriteSuccess(w, map[string]string{"message": "Match deleted successfully"})
}

// GetMatchesBySeries handles GET /api/v1/matches/series/{series_id}
func (h *MatchHandler) GetMatchesBySeries(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: GetMatchesBySeries handler called")

	seriesID := chi.URLParam(r, "series_id")
	if seriesID == "" {
		log.Printf("DEBUG: Series ID is required")
		utils.WriteValidationError(w, "Series ID is required", nil)
		return
	}

	log.Printf("DEBUG: Getting matches for series: %s", seriesID)

	// Get matches by series
	matches, err := h.service.GetMatchesBySeries(r.Context(), seriesID)
	if err != nil {
		log.Printf("DEBUG: service.GetMatchesBySeries failed: %v", err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("DEBUG: Retrieved %d matches for series", len(matches))
	utils.WriteSuccess(w, matches)
}
