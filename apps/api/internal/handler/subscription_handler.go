package handler

import (
	"net/http"

	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
)

func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if !h.StripeService.IsConfigured() {
		httpx.JSONError(w, http.StatusServiceUnavailable, "Payments not configured")
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	user, err := h.UserRepo.FindByID(userID)
	if err != nil || user == nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}

	url, err := h.StripeService.CreateCheckoutSession(userID, user.Email)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create checkout session")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) GetMySubscription(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	sub, err := h.SubscriptionRepo.GetByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get subscription")
		return
	}
	if sub == nil {
		httpx.JSON(w, http.StatusOK, map[string]interface{}{"active": false})
		return
	}
	httpx.JSON(w, http.StatusOK, sub)
}

func (h *Handler) CreatePortalSession(w http.ResponseWriter, r *http.Request) {
	if !h.StripeService.IsConfigured() {
		httpx.JSONError(w, http.StatusServiceUnavailable, "Payments not configured")
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	url, err := h.StripeService.CreatePortalSession(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create portal session")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	if err := h.StripeService.HandleWebhook(r); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
