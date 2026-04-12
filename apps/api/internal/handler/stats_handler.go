package handler

import (
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type timeDonutEntry struct {
	ItemType string `db:"item_type" json:"itemType"`
	Count    int64  `db:"count" json:"count"`
}

type heatmapDay struct {
	Date  string `db:"date" json:"date"`
	Count int    `db:"count" json:"count"`
}

type genreEntry struct {
	Genre string `db:"genre" json:"genre"`
	Count int64  `db:"count" json:"count"`
}

func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	db := h.ListRepo.DB

	var donut []timeDonutEntry
	_ = db.Select(&donut, `
		SELECT mi.item_type, COUNT(*)::bigint AS count
		FROM user_list_items uli
		JOIN media_items mi ON mi.id = uli.media_item_id
		WHERE uli.user_id = $1 AND uli.status = 'completed'
		GROUP BY mi.item_type ORDER BY count DESC
	`, userID)

	var genres []genreEntry
	_ = db.Select(&genres, `
		SELECT unnest(mi.genres) AS genre, COUNT(*)::bigint AS count
		FROM user_list_items uli
		JOIN media_items mi ON mi.id = uli.media_item_id
		WHERE uli.user_id = $1 AND uli.status = 'completed' AND mi.genres IS NOT NULL
		GROUP BY genre ORDER BY count DESC LIMIT 12
	`, userID)

	httpx.JSON(w, http.StatusOK, map[string]interface{}{
		"timeDonut": donut,
		"genres":    genres,
	})
}

func (h *Handler) GetUserHeatmap(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var days []heatmapDay
	_ = h.ListRepo.DB.Select(&days, `
		SELECT TO_CHAR(uli.updated_at::date, 'YYYY-MM-DD') AS date, COUNT(*) AS count
		FROM user_list_items uli
		WHERE uli.user_id = $1 AND uli.status = 'completed'
		AND uli.updated_at >= NOW() - INTERVAL '1 year'
		GROUP BY uli.updated_at::date ORDER BY date
	`, userID)

	httpx.JSON(w, http.StatusOK, days)
}
