package handler

import (
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

// @Summary      Follow a user
// @Description  Creates a following relationship with another user
// @Tags         Social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username to follow"
// @Success      201
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /users/{username}/follow [post]
func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	targetUsername := chi.URLParam(r, "username")
	targetUser, err := h.UserRepo.FindByUsername(targetUsername)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	if authUserID == targetUser.ID {
		httpx.JSONError(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}
	if _, err := h.SocialRepo.CreateUserFollow(authUserID, targetUser.ID); err != nil {
		httpx.JSONError(w, http.StatusConflict, "Failed to follow user")
		return
	}
	w.WriteHeader(http.StatusCreated)
}
// @Summary      Unfollow a user
// @Description  Removes a following relationship
// @Tags         Social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username to unfollow"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{username}/unfollow [delete]
func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	targetUsername := chi.URLParam(r, "username")
	targetUser, err := h.UserRepo.FindByUsername(targetUsername)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	affected, err := h.SocialRepo.DeleteUserFollow(authUserID, targetUser.ID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to unfollow user")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Follow relationship not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
// @Summary      Get user followers
// @Description  Retrieves a list of users following the target user
// @Tags         Social
// @Accept       json
// @Produce      json
// @Param        username path string true "Username"
// @Success      200  {array}   models.UserFollowResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{username}/followers [get]
func (h *Handler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	targetUsername := chi.URLParam(r, "username")
	targetUser, err := h.UserRepo.FindByUsername(targetUsername)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	list, err := h.SocialRepo.GetUserFollowers(targetUser.ID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get followers")
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}
// @Summary      Get user following
// @Description  Retrieves a list of users the target user follows
// @Tags         Social
// @Accept       json
// @Produce      json
// @Param        username path string true "Username"
// @Success      200  {array}   models.UserFollowResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{username}/following [get]
func (h *Handler) GetUserFollowing(w http.ResponseWriter, r *http.Request) {
	targetUsername := chi.URLParam(r, "username")
	targetUser, err := h.UserRepo.FindByUsername(targetUsername)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	list, err := h.SocialRepo.GetUserFollowing(targetUser.ID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get following")
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}
// @Summary      Get follow suggestions
// @Description  Returns a list of users the current user might want to follow
// @Tags         Social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit"
// @Success      200  {array}   models.UserFollowResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/suggestions [get]
func (h *Handler) GetFollowSuggestions(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	limit := int64(10)
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.ParseInt(l, 10, 64); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}
	list, err := h.SocialRepo.GetFollowSuggestions(authUserID, limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get suggestions")
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}
// @Summary      Get user profile
// @Description  Retrieves detailed profile information and stats for a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        username path string true "Username"
// @Success      200  {object}  models.EnhancedUserProfile
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{username}/profile [get]
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := middleware.UserIDFromContext(r.Context())
	targetUsername := chi.URLParam(r, "username")
	user, err := h.UserRepo.FindByUsername(targetUsername)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	targetID := user.ID
	stats, err := h.SocialRepo.GetUserStats(targetID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to load stats")
		return
	}
	var isFollowing *bool
	var isFollowedBy *bool
	if authUserID != 0 && authUserID != targetID {
		f1, err1 := h.SocialRepo.CheckIsFollowing(authUserID, targetID)
		f2, err2 := h.SocialRepo.CheckIsFollowing(targetID, authUserID)
		if err1 == nil && err2 == nil {
			isFollowing = &f1
			isFollowedBy = &f2
		}
	}
	resp := models.EnhancedUserProfile{
		ID:           user.ID,
		Username:     user.Username,
		Bio:          nil,
		FavoriteGenres: nil,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Stats:        stats,
		IsFollowing:  isFollowing,
		IsFollowedBy: isFollowedBy,
	}
	httpx.JSON(w, http.StatusOK, resp)
}
// @Summary      Search users
// @Description  Searches for users by username
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        q query string true "Search query"
// @Param        limit query int false "Limit"
// @Success      200  {array}   models.UserFollowResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/search [get]
func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	users, err := h.UserRepo.SearchByUsername(query, limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to search for users")
		return
	}
	httpx.JSON(w, http.StatusOK, users)
}