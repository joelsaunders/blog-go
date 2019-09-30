package api

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/blog-go/pkg/config"
	"github.com/joelsaunders/blog-go/pkg/repository/post"
	"github.com/joelsaunders/blog-go/pkg/repository/user"
)

func Routes(config *config.Config, db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes(&user.PGUserStore{DB: db}, config))
		r.Mount("/posts", PostRoutes(&post.PGPostStore{DB: db}, config))
	})

	return router
}
