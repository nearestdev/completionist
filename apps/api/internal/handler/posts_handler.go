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

// @Summary      Create a post
// @Description  Creates a new social post regarding a media item or general topic
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.NewPost true "Post content"
// @Success      201  {object}  models.Post
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts [post]
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

	if post.PostType == models.PostTypeReview {
		sid := strconv.FormatInt(post.ID, 10)
		desc := "Wrote a review"
		if post.Title != nil {
			desc = "Review: " + *post.Title
		}
		_ = h.UserRepo.AddXP(userID, "review", &sid, &desc)

		meta := map[string]interface{}{}
		if post.MediaItemID != nil {
			meta["media_item_id"] = post.MediaItemID.String()
		}
		_ = h.ChallengeService.NotifyAction(r.Context(), userID, "review", meta)
	}

	httpx.JSON(w, http.StatusCreated, post)
}

// @Summary      Get posts feed
// @Description  Retrieves the global or personalized activity feed
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit"
// @Param        offset query int false "Offset"
// @Success      200  {array}   models.PostWithDetails
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts [get]
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

// @Summary      Get post by ID
// @Description  Retrieves a single post by its ID
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Success      200  {object}  models.PostWithDetails
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id} [get]
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

// @Summary      Get user posts
// @Description  Retrieves posts made by a specific user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username"
// @Param        limit query int false "Limit"
// @Param        offset query int false "Offset"
// @Success      200  {array}   models.PostWithDetails
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{username}/posts [get]
func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	username := chi.URLParam(r, "username")
	targetUser, err := h.UserRepo.FindByUsername(username)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
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

	posts, err := h.PostsRepo.GetPostsByUserID(targetUser.ID, nil, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to load posts")
		return
	}
	httpx.JSON(w, http.StatusOK, posts)
}

// @Summary      Update post
// @Description  Updates an existing post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Param        request body models.UpdatePost true "Update details"
// @Success      200  {object}  models.Post
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id} [patch]
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

// @Summary      Delete post
// @Description  Deletes a post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id} [delete]
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

// @Summary      Like a post
// @Description  Adds a like to a post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Success      201  {object}  models.PostLike
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id}/like [post]
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

// @Summary      Unlike a post
// @Description  Removes a like from a post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id}/unlike [delete]
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

// @Summary      Create a comment
// @Description  Adds a comment to a post
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Param        request body models.NewComment true "Comment content"
// @Success      201  {object}  models.Comment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id}/comments [post]
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

	sid := strconv.FormatInt(comment.ID, 10)
	desc := "Commented on post " + strconv.FormatInt(postID, 10)
	_ = h.UserRepo.AddXP(userID, "comment", &sid, &desc)

	httpx.JSON(w, http.StatusCreated, comment)
}

// @Summary      Get post comments
// @Description  Retrieves comments for a specific post
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Success      200  {array}   models.CommentWithDetails
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{post_id}/comments [get]
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

// @Summary      Update comment
// @Description  Updates a comment's content
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id path int true "Comment ID"
// @Param        request body models.UpdateComment true "Update details"
// @Success      200  {object}  models.Comment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comments/{comment_id} [patch]
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

// @Summary      Delete comment
// @Description  Deletes a comment
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id path int true "Comment ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comments/{comment_id} [delete]
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

// @Summary      Like a comment
// @Description  Adds a like to a comment
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id path int true "Comment ID"
// @Success      201  {object}  models.CommentLike
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comments/{comment_id}/like [post]
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

// @Summary      Unlike a comment
// @Description  Removes a like from a comment
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id path int true "Comment ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comments/{comment_id}/unlike [delete]
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