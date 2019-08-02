package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/joelsaunders/bilbo-go/config"
)

func UserRoutes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()
	// router.Get("/{todoID}", GetATodo(configuration))
	// router.Delete("/{todoID}", DeleteTodo(configuration))
	// router.Post("/", CreateTodo(configuration))
	router.Get("/", GetAllUsers(configuration))
	return router
}

type User struct {
	Username string
}

func GetAllUsers(configuration *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := []User{
			{
				Username: "Joel",
			},
		}
		render.JSON(w, r, users) // A chi router helper for serializing and returning json
	}
}
