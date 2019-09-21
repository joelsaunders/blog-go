package api

import (
	"net/http"

	"github.com/joelsaunders/bilbo-go/repository"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func UserRoutes(userStore repository.UserStore) *chi.Mux {
	router := chi.NewRouter()
	// router.Get("/{todoID}", GetATodo(configuration))
	// router.Delete("/{todoID}", DeleteTodo(configuration))
	// router.Post("/", CreateTodo(configuration))
	router.Get("/", NewUserHandler(userStore).getUserList())
	return router
}

type UserHandler struct {
	store repository.UserStore
}

func NewUserHandler(userStore repository.UserStore) *UserHandler {
	return &UserHandler{userStore}
}

func (uh UserHandler) getUserList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := uh.store.List(r.Context(), 10)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		render.JSON(w, r, users)
	}
}
