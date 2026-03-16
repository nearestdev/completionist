package handler

import (
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
)

// @Summary      List users (admin)
// @Description  Returns a paginated list of users for admin management
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit" default(50)
// @Param        offset query int false "Offset" default(0)
// @Success      200  {array}   models.User
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users [get]
func (h *Handler) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 200 {
			limit = v
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	users, err := h.UserRepo.ListUsers(limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	httpx.JSON(w, http.StatusOK, users)
}
