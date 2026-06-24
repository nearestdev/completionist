package handler

import (
	"net/http"
	"strconv"

	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

// @Summary      Get Active Season
// @Description  Retrieves the current active season information
// @Tags         Seasons
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Season
// @Failure      404  {object}  map[string]string
// @Router       /seasons/current [get]
func (h *Handler) GetActiveSeason(w http.ResponseWriter, r *http.Request) {
	season, err := h.ChallengeRepo.GetActiveSeason()
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "No active season found")
		return
	}
	httpx.JSON(w, http.StatusOK, season)
}

// @Summary      Get Challenges
// @Description  Retrieves list of challenges
// @Tags         Challenges
// @Accept       json
// @Produce      json
// @Param        frequency query string false "Filter by frequency (daily, weekly, etc)"
// @Success      200  {array}   models.Challenge
// @Failure      500  {object}  map[string]string
// @Router       /challenges [get]
func (h *Handler) GetChallenges(w http.ResponseWriter, r *http.Request) {
	freq := r.URL.Query().Get("frequency")
	if freq == "" {
		freq = string(models.FrequencyDaily)
	}

	chals, err := h.ChallengeRepo.GetChallengesByFrequency(models.ChallengeFrequency(freq))
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch challenges")
		return
	}
	httpx.JSON(w, http.StatusOK, chals)
}

// @Summary      Get My Challenge Progress
// @Description  Retrieves the current user's progress on challenges
// @Tags         Challenges
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.UserChallenge
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /me/challenges [get]
func (h *Handler) GetMyChallenges(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	season, _ := h.ChallengeRepo.GetActiveSeason()
	var seasonID *int64
	if season != nil {
		seasonID = &season.ID
	}

	ucs, err := h.ChallengeRepo.GetUserChallenges(userID, seasonID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch user challenges")
		return
	}
	httpx.JSON(w, http.StatusOK, ucs)
}

// @Summary      Get Season Leaderboard
// @Description  Retrieves the leaderboard for the current season based on XP gained
// @Tags         Seasons
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit"
// @Param        offset query int false "Offset"
// @Success      200  {array}   repository.LeaderboardEntry
// @Failure      500  {object}  map[string]string
// @Router       /seasons/leaderboard [get]
func (h *Handler) GetSeasonLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	season, err := h.ChallengeRepo.GetActiveSeason()
	if err != nil {
		entries, err := h.ChallengeRepo.GetGlobalLeaderboard(limit, offset)
		if err != nil {
			httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch leaderboard")
			return
		}
		httpx.JSON(w, http.StatusOK, entries)
		return
	}

	entries, err := h.ChallengeRepo.GetSeasonLeaderboard(season.StartAt, season.EndAt, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch season leaderboard")
		return
	}
	httpx.JSON(w, http.StatusOK, entries)
}
