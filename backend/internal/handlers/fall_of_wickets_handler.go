package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// FallOfWicketsHandler handles HTTP requests for fall of wickets
type FallOfWicketsHandler struct {
	fallOfWicketsService *services.FallOfWicketsService
}

// NewFallOfWicketsHandler creates a new fall of wickets handler
func NewFallOfWicketsHandler(fallOfWicketsService *services.FallOfWicketsService) *FallOfWicketsHandler {
	return &FallOfWicketsHandler{
		fallOfWicketsService: fallOfWicketsService,
	}
}

// CreateFallOfWickets handles POST /fall-of-wickets
func (h *FallOfWicketsHandler) CreateFallOfWickets(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFallOfWicketsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.CreateFallOfWickets(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fallOfWickets)
}

// GetFallOfWicketsByID handles GET /fall-of-wickets/{id}
func (h *FallOfWicketsHandler) GetFallOfWicketsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.GetFallOfWicketsByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallOfWickets)
}

// ListFallOfWickets handles GET /fall-of-wickets
func (h *FallOfWicketsHandler) ListFallOfWickets(w http.ResponseWriter, r *http.Request) {
	filters := &models.FallOfWicketsFilters{}

	// Parse query parameters
	if matchID := r.URL.Query().Get("match_id"); matchID != "" {
		filters.MatchID = &matchID
	}
	if inningsID := r.URL.Query().Get("innings_id"); inningsID != "" {
		filters.InningsID = &inningsID
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	fallOfWickets, err := h.fallOfWicketsService.ListFallOfWickets(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallOfWickets)
}

// UpdateFallOfWickets handles PUT /fall-of-wickets/{id}
func (h *FallOfWicketsHandler) UpdateFallOfWickets(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateFallOfWicketsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.UpdateFallOfWickets(r.Context(), id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallOfWickets)
}

// DeleteFallOfWickets handles DELETE /fall-of-wickets/{id}
func (h *FallOfWicketsHandler) DeleteFallOfWickets(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	err := h.fallOfWicketsService.DeleteFallOfWickets(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetFallOfWicketsByMatchID handles GET /matches/{match_id}/fall-of-wickets
func (h *FallOfWicketsHandler) GetFallOfWicketsByMatchID(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.GetFallOfWicketsByMatchID(r.Context(), matchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallOfWickets)
}

// GetFallOfWicketsByInningsID handles GET /innings/{innings_id}/fall-of-wickets
func (h *FallOfWicketsHandler) GetFallOfWicketsByInningsID(w http.ResponseWriter, r *http.Request) {
	inningsID := chi.URLParam(r, "innings_id")
	if inningsID == "" {
		http.Error(w, "Innings ID is required", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.GetFallOfWicketsByInningsID(r.Context(), inningsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallOfWickets)
}

// GetFallOfWicketsByBallID handles GET /balls/{ball_id}/fall-of-wickets
func (h *FallOfWicketsHandler) GetFallOfWicketsByBallID(w http.ResponseWriter, r *http.Request) {
	ballID := chi.URLParam(r, "ball_id")
	if ballID == "" {
		http.Error(w, "Ball ID is required", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.GetFallOfWicketsByBallID(r.Context(), ballID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fallOfWickets)
}

// GetFallOfWicketsSummary handles GET /fall-of-wickets/summary
func (h *FallOfWicketsHandler) GetFallOfWicketsSummary(w http.ResponseWriter, r *http.Request) {
	matchID := r.URL.Query().Get("match_id")
	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	inningsID := r.URL.Query().Get("innings_id")
	var inningsIDPtr *string
	if inningsID != "" {
		inningsIDPtr = &inningsID
	}

	summary, err := h.fallOfWicketsService.GetFallOfWicketsSummary(r.Context(), matchID, inningsIDPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// CreateFallOfWicketsFromBall handles POST /fall-of-wickets/from-ball
func (h *FallOfWicketsHandler) CreateFallOfWicketsFromBall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BallID string `json:"ball_id"`
		Score  int    `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fallOfWickets, err := h.fallOfWicketsService.CreateFallOfWicketsFromBall(r.Context(), req.BallID, req.Score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fallOfWickets)
}
