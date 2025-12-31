package routes

import (
	"net/http"
	"strings"

	"github.com/GATEOPENERZ/completionist-api/internal/config"
	"github.com/GATEOPENERZ/completionist-api/internal/handler"
	appmw "github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
		r.Get("/auth/steam/callback", h.SteamCallback)
		r.Get("/auth/lastfm/callback", h.LastFMCallback)
		r.Group(func(pr chi.Router) {
			pr.Use(appmw.AuthMiddleware(h.UserRepo))
			pr.Get("/me", h.GetMe)
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
			})
		})
	})
	return r
}