package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var payload models.NewPost
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	post, err := h.PostsRepo.CreatePost(userID, payload)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}
	httpx.JSON(w, http.StatusCreated, post)
}

func (h *Handler) GetPostsFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	limit := int64(20)
	offset := int64(0)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n >= 0 {
			offset = n
		}
	}
	posts, err := h.PostsRepo.GetPostsWithDetails(&userID, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to load posts")
		return
	}
	httpx.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	post, err := h.PostsRepo.GetPostByIDWithDetails(postID, &userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to load post")
		return
	}
	if post == nil {
		httpx.JSONError(w, http.StatusNotFound, "Post not found")
		return
	}
	httpx.JSON(w, http.StatusOK, post)
}

func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	uidStr := chi.URLParam(r, "user_id")
	targetUserID, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}
	limit := int64(20)
	offset := int64(0)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n >= 0 {
			offset = n
		}
	}
	posts, err := h.PostsRepo.GetPostsByUserID(targetUserID, nil, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to load posts")
		return
	}
	httpx.JSON(w, http.StatusOK, posts)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	var payload models.UpdatePost
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	post, err := h.PostsRepo.UpdatePost(postID, userID, payload)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}
	httpx.JSON(w, http.StatusOK, post)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	affected, err := h.PostsRepo.DeletePost(postID, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Post not found or no permission")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	like, err := h.PostsRepo.CreatePostLike(userID, postID)
	if err != nil {
		httpx.JSONError(w, http.StatusConflict, "Failed to like post")
		return
	}
	httpx.JSON(w, http.StatusCreated, like)
}

func (h *Handler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	affected, err := h.PostsRepo.DeletePostLike(userID, postID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to unlike post")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Like not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	var payload models.NewComment
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	payload.PostID = postID
	comment, err := h.PostsRepo.CreateComment(userID, payload)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}
	httpx.JSON(w, http.StatusCreated, comment)
}

func (h *Handler) GetPostComments(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	comments, err := h.PostsRepo.GetCommentsWithDetails(postID, &userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to load comments")
		return
	}
	httpx.JSON(w, http.StatusOK, comments)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "comment_id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid comment_id")
		return
	}
	var payload models.UpdateComment
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	comment, err := h.PostsRepo.UpdateComment(commentID, userID, payload)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update comment")
		return
	}
	httpx.JSON(w, http.StatusOK, comment)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "comment_id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid comment_id")
		return
	}
	affected, err := h.PostsRepo.DeleteComment(commentID, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to delete comment")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Comment not found or no permission")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LikeComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "comment_id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid comment_id")
		return
	}
	like, err := h.PostsRepo.CreateCommentLike(userID, commentID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to like comment")
		return
	}
	httpx.JSON(w, http.StatusCreated, like)
}

func (h *Handler) UnlikeComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "comment_id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid comment_id")
		return
	}
	affected, err := h.PostsRepo.DeleteCommentLike(userID, commentID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to unlike comment")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Like not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
