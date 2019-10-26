package api

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"

	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/repository/post"
	"github.com/joelsaunders/blog-go/api/pkg/repository/tag"
	"github.com/joelsaunders/blog-go/api/pkg/repository/user"
)

// Routes is the base router mux for the server, it contains all routes available
func Routes(config *config.Config, db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	corsConfig := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(corsConfig.Handler)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes(&user.PGUserStore{DB: db}, config))
		r.Mount("/posts", PostRoutes(&post.PGPostStore{DB: db}, config))
		r.Mount("/tags", TagRoutes(&tag.PGTagStore{DB: db}, config))
		r.Get("/sitemap.xml", GetSitemap(&post.PGPostStore{DB: db}, config))
	})

	return router
}
