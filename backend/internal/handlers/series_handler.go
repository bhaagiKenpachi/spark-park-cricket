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

// SeriesHandler handles series-related HTTP requests
type SeriesHandler struct {
	service services.SeriesServiceInterface
}

// NewSeriesHandler creates a new series handler
func NewSeriesHandler(service services.SeriesServiceInterface) *SeriesHandler {
	return &SeriesHandler{
		service: service,
	}
}

// ListSeries handles GET /api/v1/series
func (h *SeriesHandler) ListSeries(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== SERIES HANDLER: ListSeries ===")
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Request method: %s", r.Method)

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	filters := &models.SeriesFilters{
		Limit:  limit,
		Offset: offset,
	}

	log.Printf("Filters: %+v", filters)
	log.Printf("Calling service.ListSeries...")

	series, err := h.service.ListSeries(r.Context(), filters)
	if err != nil {
		log.Printf("ERROR: Service call failed: %v", err)
		utils.WriteInternalError(w, err.Error())
		return
	}

	log.Printf("Service returned %d series", len(series))
	for i, s := range series {
		log.Printf("Series %d: ID=%s, Name=%s", i+1, s.ID, s.Name)
	}

	log.Printf("Writing success response...")
	utils.WriteSuccess(w, series)
}

// CreateSeries handles POST /api/v1/series
func (h *SeriesHandler) CreateSeries(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	series, err := h.service.CreateSeries(r.Context(), &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteCreated(w, series)
}

// GetSeries handles GET /api/v1/series/{id}
func (h *SeriesHandler) GetSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Series ID is required", nil)
		return
	}

	series, err := h.service.GetSeries(r.Context(), id)
	if err != nil {
		utils.WriteNotFound(w, "Series")
		return
	}

	utils.WriteSuccess(w, series)
}

// UpdateSeries handles PUT /api/v1/series/{id}
func (h *SeriesHandler) UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Series ID is required", nil)
		return
	}

	var req models.UpdateSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	series, err := h.service.UpdateSeries(r.Context(), id, &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, series)
}

// DeleteSeries handles DELETE /api/v1/series/{id}
func (h *SeriesHandler) DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Series ID is required", nil)
		return
	}

	err := h.service.DeleteSeries(r.Context(), id)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Series deleted successfully"})
}
