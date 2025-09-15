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

type ScorecardHandler struct {
	scorecardService services.ScorecardServiceInterface
}

// NewScorecardHandler creates a new scorecard handler
func NewScorecardHandler(scorecardService services.ScorecardServiceInterface) *ScorecardHandler {
	return &ScorecardHandler{
		scorecardService: scorecardService,
	}
}

// StartScoring starts scoring for a match
func (h *ScorecardHandler) StartScoring(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting scoring request for match")

	var req models.ScorecardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	// Start scoring
	err := h.scorecardService.StartScoring(r.Context(), req.MatchID)
	if err != nil {
		log.Printf("Error starting scoring: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Return success response
	response := map[string]interface{}{
		"message":  "Scoring started successfully",
		"match_id": req.MatchID,
	}

	log.Printf("Successfully started scoring for match %s", req.MatchID)
	utils.WriteSuccessResponse(w, response)
}

// AddBall adds a ball to the scorecard
func (h *ScorecardHandler) AddBall(w http.ResponseWriter, r *http.Request) {

	var req models.BallEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate request using cricket-specific validation
	if err := utils.ValidateBallEventRequest(&req); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	// Add ball
	err := h.scorecardService.AddBall(r.Context(), &req)
	if err != nil {
		log.Printf("Error adding ball: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Return success response
	response := map[string]interface{}{
		"message":        "Ball added successfully",
		"match_id":       req.MatchID,
		"innings_number": req.InningsNumber,
		"ball_type":      req.BallType,
		"run_type":       req.RunType,
		"runs":           req.RunType.GetRunValue(),
		"byes":           req.Byes,
		"is_wicket":      req.IsWicket,
	}

	log.Printf("Successfully added ball for match %s", req.MatchID)
	utils.WriteSuccessResponse(w, response)
}

// GetScorecard gets the complete scorecard for a match
func (h *ScorecardHandler) GetScorecard(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		log.Printf("Missing match_id parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "match_id is required")
		return
	}

	log.Printf("Getting scorecard for match %s", matchID)

	// Get scorecard
	scorecard, err := h.scorecardService.GetScorecard(r.Context(), matchID)
	if err != nil {
		log.Printf("Error getting scorecard: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	log.Printf("Successfully retrieved scorecard for match %s", matchID)
	utils.WriteSuccessResponse(w, scorecard)
}

// GetCurrentOver gets the current over for a match
func (h *ScorecardHandler) GetCurrentOver(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		log.Printf("Missing match_id parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "match_id is required")
		return
	}

	inningsNumberStr := r.URL.Query().Get("innings")
	if inningsNumberStr == "" {
		inningsNumberStr = "1" // Default to first innings
	}

	inningsNumber, err := strconv.Atoi(inningsNumberStr)
	if err != nil || inningsNumber < 1 || inningsNumber > 2 {
		log.Printf("Invalid innings number: %s", inningsNumberStr)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", "innings must be 1 or 2")
		return
	}

	log.Printf("Getting current over for match %s, innings %d", matchID, inningsNumber)

	// Get current over
	over, err := h.scorecardService.GetCurrentOver(r.Context(), matchID, inningsNumber)
	if err != nil {
		log.Printf("Error getting current over: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	log.Printf("Successfully retrieved current over %d for match %s, innings %d", over.OverNumber, matchID, inningsNumber)
	utils.WriteSuccessResponse(w, over)
}

// GetInnings gets innings details for a match
func (h *ScorecardHandler) GetInnings(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		log.Printf("Missing match_id parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "match_id is required")
		return
	}

	inningsNumberStr := chi.URLParam(r, "innings_number")
	if inningsNumberStr == "" {
		log.Printf("Missing innings_number parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "innings_number is required")
		return
	}

	inningsNumber, err := strconv.Atoi(inningsNumberStr)
	if err != nil || inningsNumber < 1 || inningsNumber > 2 {
		log.Printf("Invalid innings number: %s", inningsNumberStr)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", "innings_number must be 1 or 2")
		return
	}

	log.Printf("Getting innings %d for match %s", inningsNumber, matchID)

	// Get innings from scorecard
	scorecard, err := h.scorecardService.GetScorecard(r.Context(), matchID)
	if err != nil {
		log.Printf("Error getting scorecard: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Find the requested innings
	var innings *models.InningsSummary
	for _, inn := range scorecard.Innings {
		if inn.InningsNumber == inningsNumber {
			innings = &inn
			break
		}
	}

	if innings == nil {
		log.Printf("Innings %d not found for match %s", inningsNumber, matchID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Innings not found")
		return
	}

	log.Printf("Successfully retrieved innings %d for match %s", inningsNumber, matchID)
	utils.WriteSuccessResponse(w, innings)
}

// GetOver gets over details for a match
func (h *ScorecardHandler) GetOver(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		log.Printf("Missing match_id parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "match_id is required")
		return
	}

	inningsNumberStr := chi.URLParam(r, "innings_number")
	if inningsNumberStr == "" {
		log.Printf("Missing innings_number parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "innings_number is required")
		return
	}

	inningsNumber, err := strconv.Atoi(inningsNumberStr)
	if err != nil || inningsNumber < 1 || inningsNumber > 2 {
		log.Printf("Invalid innings number: %s", inningsNumberStr)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", "innings_number must be 1 or 2")
		return
	}

	overNumberStr := chi.URLParam(r, "over_number")
	if overNumberStr == "" {
		log.Printf("Missing over_number parameter")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "MISSING_PARAMETER", "over_number is required")
		return
	}

	overNumber, err := strconv.Atoi(overNumberStr)
	if err != nil || overNumber < 1 {
		log.Printf("Invalid over number: %s", overNumberStr)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", "over_number must be a positive integer")
		return
	}

	log.Printf("Getting over %d for innings %d, match %s", overNumber, inningsNumber, matchID)

	// Get over from scorecard
	scorecard, err := h.scorecardService.GetScorecard(r.Context(), matchID)
	if err != nil {
		log.Printf("Error getting scorecard: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Find the requested over
	var over *models.OverSummary
	for _, inn := range scorecard.Innings {
		if inn.InningsNumber == inningsNumber {
			for _, ov := range inn.Overs {
				if ov.OverNumber == overNumber {
					over = &ov
					break
				}
			}
			break
		}
	}

	if over == nil {
		log.Printf("Over %d not found for innings %d, match %s", overNumber, inningsNumber, matchID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Over not found")
		return
	}

	log.Printf("Successfully retrieved over %d for innings %d, match %s", overNumber, inningsNumber, matchID)
	utils.WriteSuccessResponse(w, over)
}
