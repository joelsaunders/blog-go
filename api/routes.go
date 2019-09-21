package api

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/joelsaunders/bilbo-go/repository"
)

type DB interface {
	Users() repository.UserStore
	// Test() repository.UserStore
}

type Server struct {
	db     DB
	Router *chi.Mux
}

func NewServer(db DB) *Server {
	server := new(Server)
	server.db = db
	router := Routes(db)
	server.Router = router
	return server
}

func Routes(db DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes(db.Users()))
	})

	return router
}
