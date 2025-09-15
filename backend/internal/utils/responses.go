package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If we can't write the response, we can't do much more
		// The client will receive a partial response
		GetLogger().Error("Failed to encode JSON response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	response := APIResponse{Data: data}
	WriteJSON(w, http.StatusOK, response)
}

// WriteCreated writes a created response
func WriteCreated(w http.ResponseWriter, data interface{}) {
	response := APIResponse{Data: data}
	WriteJSON(w, http.StatusCreated, response)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, statusCode int, code, message string, details interface{}) {
	response := APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	WriteJSON(w, statusCode, response)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, message string, details interface{}) {
	WriteError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, details)
}

// WriteNotFound writes a not found error response
func WriteNotFound(w http.ResponseWriter, resource string) {
	WriteError(w, http.StatusNotFound, "NOT_FOUND", resource+" not found", nil)
}

// WriteInternalError writes an internal server error response
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

// WriteJSONResponse writes a JSON response with the given status code (alias for WriteJSON)
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	WriteJSON(w, statusCode, data)
}

// WriteErrorResponse writes an error response (alias for WriteError)
func WriteErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	WriteError(w, statusCode, code, message, nil)
}

// WriteSuccessResponse writes a success response (alias for WriteSuccess)
func WriteSuccessResponse(w http.ResponseWriter, data interface{}) {
	WriteSuccess(w, data)
}
