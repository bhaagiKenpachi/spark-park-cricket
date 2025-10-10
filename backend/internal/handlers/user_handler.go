package handlers

import (
	"encoding/json"
	"net/http"
	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserServiceInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UpdateUserNameRequest represents the request to update user name
type UpdateUserNameRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// UpdateUserName handles PUT /api/v1/users/me
func (h *UserHandler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User authentication required", nil)
		return
	}

	// Parse request body
	var req UpdateUserNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteValidationError(w, "Validation failed", err)
		return
	}

	// Update user name
	user, err := h.userService.UpdateUserName(r.Context(), userID, req.Name)
	if err != nil {
		utils.LogError(err, "Failed to update user name", map[string]interface{}{
			"user_id": userID,
		})
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_ERROR", "Failed to update user name", nil)
		return
	}

	utils.LogInfo("User name updated successfully", map[string]interface{}{
		"user_id": userID,
		"name":    req.Name,
	})

	utils.WriteSuccess(w, user)
}
