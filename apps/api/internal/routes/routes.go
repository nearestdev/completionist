package routes

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nearestdev/completionist/internal/config"
	"github.com/nearestdev/completionist/internal/handler"
	appmw "github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.Handler, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(cfg.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)

		r.Get("/search/anime", h.SearchAnimeJikan)
		r.Get("/search/manga", h.SearchMangaJikan)
		r.Get("/search/movies", h.SearchMoviesTMDB)
		r.Get("/search/tv", h.SearchTVTMDB)
		r.Get("/search/books", h.SearchBooks)
		r.Get("/games/rawg/search", h.RAWGSearchGames)

		r.Get("/media/trending", h.GetTrendingMedia)
		r.Get("/media/details", h.GetMediaDetails)
		r.Get("/media/{id}", h.GetMediaByID)

		r.Get("/auth/steam/callback", h.SteamCallback)
		r.Get("/auth/lastfm/callback", h.LastFMCallback)

		r.Get("/seasons/current", h.GetActiveSeason)
		r.Get("/seasons/leaderboard", h.GetSeasonLeaderboard)
		r.Get("/challenges", h.GetChallenges)

		r.Get("/badges", h.ListBadges)
		r.Get("/streaks/leaderboard", h.GetStreakLeaderboard)
		r.Get("/users/{id}/streak", h.GetUserStreak)
		r.Get("/users/{id}/badges", h.GetUserBadges)
		r.Get("/users/{id}/stats", h.GetUserStats)
		r.Get("/users/{id}/stats/heatmap", h.GetUserHeatmap)
		r.Get("/users/{id}/profile-custom", h.GetCustomProfile)
		r.Get("/reviews/user/{id}", h.GetUserReviews)
		r.Get("/reviews/media/{id}", h.GetMediaReviews)
		r.Get("/franchises", h.ListFranchises)
		r.Get("/franchises/{id}", h.GetFranchise)
		r.Get("/posts/{id}/og", h.GetPostOG)
		r.Get("/ads", h.GetActiveAd)
		r.Post("/ads/{id}/impression", h.RecordAdImpression)
		r.Post("/webhooks/stripe", h.HandleStripeWebhook)
		r.Get("/connections/{provider}/callback", h.OAuthCallback)

		r.Group(func(pr chi.Router) {
			pr.Use(appmw.AuthMiddleware(h.UserRepo))
			pr.Get("/me", h.GetMe)
			pr.Get("/me/xp", h.GetMyXPHistory)
			pr.Get("/me/rank", h.GetMyRank)
			pr.Get("/me/challenges", h.GetMyChallenges)
			pr.Put("/me/profile", h.UpdateProfile)
			pr.Post("/me/delete", h.RequestAccountDeletion)
			pr.Delete("/me/delete", h.CancelAccountDeletion)
			pr.Get("/me/export", h.ExportUserData)
			pr.Post("/me/appeal", h.SubmitBanAppeal)

			pr.Route("/lists", func(lr chi.Router) {
				lr.Post("/", h.CreateListItem)
				lr.Get("/", h.GetMyListItems)
				lr.Patch("/{item_id}", h.UpdateMyListItem)
				lr.Delete("/{item_id}", h.DeleteMyListItem)

				lr.Post("/jikan/anime/{mal_id}", h.AddAnimeFromJikan)
				lr.Post("/jikan/manga/{mal_id}", h.AddMangaFromJikan)
				lr.Post("/tmdb/movie/{tmdb_id}", h.AddMovieFromTMDB)
				lr.Post("/tmdb/tv/{tmdb_id}", h.AddTVFromTMDB)
				lr.Post("/google/books/{book_id}", h.AddBookFromGoogleBooks)
			})

			pr.Route("/wishlist", func(wr chi.Router) {
				wr.Post("/", h.AddToWishlist)
				wr.Get("/", h.GetMyWishlist)
				wr.Delete("/{item_id}", h.RemoveFromWishlist)
			})

			pr.Route("/users", func(ur chi.Router) {
				ur.Get("/search", h.SearchUsers)
				ur.Post("/{username}/follow", h.FollowUser)
				ur.Delete("/{username}/unfollow", h.UnfollowUser)
				ur.Get("/{username}/followers", h.GetUserFollowers)
				ur.Get("/{username}/following", h.GetUserFollowing)
				ur.Get("/{username}/profile", h.GetUserProfile)
				ur.Get("/{username}/posts", h.GetUserPosts)
				ur.Get("/suggestions", h.GetFollowSuggestions)
				ur.Get("/{username}/steam", h.GetUserSteamAccount)
			})

			pr.Route("/posts", func(prp chi.Router) {
				prp.Post("/", h.CreatePost)
				prp.Get("/", h.GetPostsFeed)
				prp.Get("/{post_id}", h.GetPostByID)
				prp.Patch("/{post_id}", h.UpdatePost)
				prp.Delete("/{post_id}", h.DeletePost)
				prp.Post("/{post_id}/like", h.LikePost)
				prp.Delete("/{post_id}/unlike", h.UnlikePost)
				prp.Post("/{post_id}/comments", h.CreateComment)
				prp.Get("/{post_id}/comments", h.GetPostComments)
			})

			pr.Patch("/comments/{comment_id}", h.UpdateComment)
			pr.Delete("/comments/{comment_id}", h.DeleteComment)
			pr.Post("/comments/{comment_id}/like", h.LikeComment)
			pr.Delete("/comments/{comment_id}/unlike", h.UnlikeComment)

			pr.Get("/auth/steam/login", h.SteamLogin)
			pr.Get("/me/steam", h.GetMySteamAccount)
			pr.Post("/me/steam/attach", h.AttachSteamToUser)
			pr.Get("/me/steam/owned", h.GetMySteamOwnedGames)
			pr.Get("/steam/achievements/{app_id}", h.GetMySteamAchievementsForApp)
			pr.Get("/steam/schema/{app_id}", h.GetSteamGameSchema)
			pr.Get("/games/rawg/{rawg_id}/achievements", h.RAWGGameAchievements)

			pr.Get("/auth/lastfm", h.LastFMAuth)
			pr.Get("/me/lastfm", h.GetMyLastFMAccount)
			pr.Get("/me/lastfm/recent", h.GetMyRecentTracks)

			pr.Route("/attachments", func(ar chi.Router) {
				ar.Post("/", h.CreateAttachment)
				ar.Post("/link", h.LinkAttachment)
				ar.Delete("/unlink", h.UnlinkAttachment)
				ar.Get("/by-entity", h.ListAttachmentsByEntity)
				ar.Post("/upload", h.UploadFile)
			})

			pr.Route("/messages", func(mr chi.Router) {
				mr.Post("/send", h.SendDirectMessage)
				mr.Get("/conversations", h.GetConversations)
				mr.Get("/conversation/{userId}", h.GetConversation)
				mr.Post("/{messageId}/react", h.ReactToDirectMessage)
				mr.Post("/{messageId}/pin", h.PinDirectMessage)
			})

			pr.Route("/rooms", func(rr chi.Router) {
				rr.Post("/", h.CreateRoom)
				rr.Get("/", h.GetRooms)
				rr.Get("/{roomId}", h.GetRoom)
				rr.Post("/{roomId}/join", h.JoinRoom)
				rr.Delete("/{roomId}/leave", h.LeaveRoom)
				rr.Get("/{roomId}/messages", h.GetRoomMessages)
				rr.Get("/{roomId}/members", h.GetRoomMembers)
				rr.Get("/{roomId}/state", h.GetRoomState)
				rr.Patch("/{roomId}/state", h.UpdateRoomState)
				rr.Get("/{roomId}/pinned", h.GetPinnedRoomMessages)
				rr.Post("/{roomId}/messages/{messageId}/react", h.ReactToRoomMessage)
				rr.Post("/{roomId}/messages/{messageId}/pin", h.PinRoomMessage)
			})

			pr.Route("/collections", func(cr chi.Router) {
				cr.Post("/", h.CreateCollection)
				cr.Get("/", h.GetMyCollections)
				cr.Get("/{id}", h.GetCollection)
				cr.Put("/{id}", h.UpdateCollection)
				cr.Delete("/{id}", h.DeleteCollection)
			})

			pr.Put("/lists/{item_id}/collection", h.AssignItemToCollection)
			pr.Delete("/lists/{item_id}/collection", h.UnassignItemFromCollection)

			pr.Route("/reviews", func(rr chi.Router) {
				rr.Post("/", h.CreateReview)
				rr.Put("/{id}", h.UpdateReview)
				rr.Delete("/{id}", h.DeleteReview)
			})

			pr.Route("/queue", func(qr chi.Router) {
				qr.Get("/", h.GetMyQueue)
				qr.Put("/reorder", h.ReorderQueue)
				qr.Post("/{item_id}", h.AddToQueue)
				qr.Delete("/{item_id}", h.RemoveFromQueue)
			})

			pr.Get("/suggestions", h.GetSuggestions)

			pr.Post("/posts/{id}/share", h.SharePost)

			pr.Post("/moderation/report", h.ReportContent)

			pr.Route("/subscriptions", func(sr chi.Router) {
				sr.Post("/checkout", h.CreateCheckoutSession)
				sr.Get("/me", h.GetMySubscription)
				sr.Post("/portal", h.CreatePortalSession)
			})

			pr.Route("/connections", func(cr chi.Router) {
				cr.Get("/", h.GetConnections)
				cr.Post("/{provider}/connect", h.ConnectProvider)
				cr.Delete("/{provider}", h.DisconnectProvider)
				cr.Post("/{provider}/sync", h.TriggerSync)
			})
			pr.Get("/imports", h.GetImportJobs)

			pr.Route("/admin", func(ar chi.Router) {
				ar.Use(appmw.RequireRoles(models.RoleAdmin))
				ar.Get("/users", h.AdminListUsers)

				ar.Route("/moderation", func(mr chi.Router) {
					mr.Get("/", h.AdminGetModerationQueue)
					mr.Put("/{id}", h.AdminReviewModerationItem)
				})
				ar.Get("/bans", h.AdminGetBanAppeals)
				ar.Post("/users/{id}/ban", h.AdminBanUser)
				ar.Delete("/users/{id}/ban", h.AdminUnbanUser)
				ar.Route("/ban-appeals", func(br chi.Router) {
					br.Get("/", h.AdminGetBanAppeals)
					br.Put("/{id}", h.AdminReviewBanAppeal)
				})
				ar.Route("/ads", func(adr chi.Router) {
					adr.Post("/", h.AdminCreateAdCampaign)
					adr.Get("/", h.AdminGetAdCampaigns)
					adr.Put("/{id}", h.AdminUpdateAdCampaign)
					adr.Get("/{id}/stats", h.AdminGetAdCampaignStats)
				})
				ar.Route("/franchises", func(fr chi.Router) {
					fr.Post("/", h.AdminCreateFranchise)
					fr.Put("/{id}", h.AdminUpdateFranchise)
					fr.Post("/{id}/items", h.AdminAddFranchiseItem)
					fr.Delete("/{id}/items/{itemId}", h.AdminRemoveFranchiseItem)
				})
			})

			pr.Get("/ws", h.HandleWebSocket)
		})
	})

	// Serve uploaded files regardless of whether the API runs from apps/api or repo root.
	uploadDir := config.ResolveAppPath("uploads")
	fs := http.FileServer(http.Dir(uploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))

	return r
}
