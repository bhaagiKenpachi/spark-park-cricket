package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"strconv"

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


	// Get matches from service
	matches, err := h.service.ListMatches(r.Context(), filters)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, matches)
}

// CreateMatch handles POST /api/v1/matches
func (h *MatchHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {

	var req models.CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}


	// Create match
	match, err := h.service.CreateMatch(r.Context(), &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteCreated(w, match)
}

// GetMatch handles GET /api/v1/matches/{id}
func (h *MatchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}


	// Get match
	match, err := h.service.GetMatch(r.Context(), id)
	if err != nil {
		utils.WriteNotFound(w, "Match not found")
		return
	}

	utils.WriteSuccess(w, match)
}

// UpdateMatch handles PUT /api/v1/matches/{id}
func (h *MatchHandler) UpdateMatch(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	var req models.UpdateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}


	// Update match
	match, err := h.service.UpdateMatch(r.Context(), id, &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, match)
}

// DeleteMatch handles DELETE /api/v1/matches/{id}
func (h *MatchHandler) DeleteMatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("SECURITY: DeleteMatch handler called from %s", r.RemoteAddr)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Printf("SECURITY: DeleteMatch attempted without ID from %s", r.RemoteAddr)
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	log.Printf("SECURITY: DeleteMatch request for ID: %s from %s (User-Agent: %s)", id, r.RemoteAddr, r.UserAgent())

	// Delete match
	err := h.service.DeleteMatch(r.Context(), id)
	if err != nil {
		log.Printf("SECURITY: DeleteMatch failed for %s: %v", id, err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("SECURITY: Match %s deleted successfully by %s", id, r.RemoteAddr)
	utils.WriteSuccess(w, map[string]string{"message": "Match deleted successfully"})
}

// GetMatchesBySeries handles GET /api/v1/matches/series/{series_id}
func (h *MatchHandler) GetMatchesBySeries(w http.ResponseWriter, r *http.Request) {

	seriesID := chi.URLParam(r, "series_id")
	if seriesID == "" {
		utils.WriteValidationError(w, "Series ID is required", nil)
		return
	}


	// Get matches by series
	matches, err := h.service.GetMatchesBySeries(r.Context(), seriesID)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, matches)
}
