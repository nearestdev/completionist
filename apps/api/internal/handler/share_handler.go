package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) SharePost(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var in struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	post := &models.Post{
		UserID:       userID,
		Content:      in.Content,
		PostType:     models.PostTypeShare,
		SharedPostID: &postID,
	}

	query := `
		INSERT INTO posts (user_id, content, post_type, shared_post_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, media_item_id, title, content, post_type, rating, is_spoiler, shared_post_id, created_at, updated_at
	`
	if err := h.PostsRepo.DB.QueryRowx(query, post.UserID, post.Content, post.PostType, post.SharedPostID).StructScan(post); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create share post")
		return
	}

	httpx.JSON(w, http.StatusCreated, post)
}

func (h *Handler) GetPostOG(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.PostsRepo.GetPostByIDWithDetails(postID, nil)
	if err != nil || post == nil {
		httpx.JSONError(w, http.StatusNotFound, "Post not found")
		return
	}

	title := fmt.Sprintf("%s on Completionist", post.Username)
	description := post.Content
	if len(description) > 200 {
		description = description[:200] + "..."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta property="og:title" content="%s" />
<meta property="og:description" content="%s" />
<meta property="og:type" content="article" />
<meta property="og:url" content="%s/posts/%d" />
<meta name="twitter:card" content="summary" />
<meta name="twitter:title" content="%s" />
<meta name="twitter:description" content="%s" />
</head>
<body>
<script>window.location.href = '%s/posts/%d';</script>
</body>
</html>`, title, description, h.Config.FrontendBaseURL, postID, title, description, h.Config.FrontendBaseURL, postID)
}
