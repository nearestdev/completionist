package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/auth"
	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)
type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
}
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body registerRequest true "Register credentials"
// @Success      201  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /register [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}
	if err := h.UserRepo.Create(user); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	httpx.JSON(w, http.StatusCreated, user)
}
// @Summary      Login user
// @Description  Authenticates a user and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login credentials"
// @Success      200  {object}  loginResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user, err := h.UserRepo.FindByEmail(req.Email)
	if err != nil {
		httpx.JSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httpx.JSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := auth.GenerateJWT(user.ID, time.Hour*24)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	httpx.JSON(w, http.StatusOK, loginResponse{Token: token})
}
// @Summary      Get current user
// @Description  Returns the profile of the currently logged-in user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	httpx.JSON(w, http.StatusOK, user)
}