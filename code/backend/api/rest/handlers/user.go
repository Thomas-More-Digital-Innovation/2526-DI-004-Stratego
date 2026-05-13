package handlers

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"digital-innovation/gostrategy/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCurrentUserHandler returns the currently logged-in user
// @Summary Get current user
// @Description Retrieve information about the currently logged-in user
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /users/me [get]
func (h *Handler) GetCurrentUserHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	fullUser, err := db.GetUserByID(c.Request.Context(), user.ID)
	if err != nil {
		core.SendError(c, "User not found", http.StatusNotFound)
		return
	}

	core.SendJSON(c, fullUser, http.StatusOK)
}

// GetUserHandler retrieves user information by ID
// @Summary Get user by ID
// @Description Retrieve information about a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid or missing user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [get]
func (h *Handler) GetUserHandler(c *gin.Context) {
	userID, err := core.ParseID(c, "id")
	if err != nil || userID == 0 {
		core.SendError(c, "Invalid or missing user ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		core.SendError(c, "User not found", http.StatusNotFound)
		return
	}

	core.SendJSON(c, user, http.StatusOK)
}

// ChangePasswordHandler handles password updates for the current user
// @Summary Change password
// @Description Update the password for the currently logged-in user
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Password update details"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or passwords don't match"
// @Failure 401 {object} map[string]string "Unauthorized or invalid current password"
// @Failure 500 {object} map[string]string "Failed to update password"
// @Router /users/me/password [post]
func (h *Handler) ChangePasswordHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePassword(req.OldPassword); err != nil {
		core.SendError(c, "Current password: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		core.SendError(c, "New password: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePassword(req.ConfirmPassword); err != nil {
		core.SendError(c, "Confirm password: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		core.SendError(c, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Verify current password
	_, err := db.AuthenticateUser(c.Request.Context(), user.Username, req.OldPassword)
	if err != nil {
		logging.Debug(logging.TagWeb, "Password change failed for user %s: %v", user.Username, err)
		core.SendError(c, "Invalid current password", http.StatusUnauthorized)
		return
	}

	if !IsStrongPassword(req.NewPassword) {
		core.SendError(c, "Password must be at least 8 characters and contain at least one number, one uppercase and one lowercase letter", http.StatusBadRequest)
		return
	}

	// Update password
	if err := db.UpdateUserPassword(c.Request.Context(), user.ID, req.NewPassword); err != nil {
		logging.Debug(logging.TagWeb, "Failed to update password for user %d: %v", user.ID, err)
		core.SendError(c, "Failed to update password", http.StatusInternalServerError)
		return
	}

	logging.Debug(logging.TagWeb, "Password updated successfully for user %d", user.ID)
	core.SendJSON(c, gin.H{"message": "Password updated successfully"}, http.StatusOK)
}
