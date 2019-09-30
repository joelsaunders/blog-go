package api

import (
	"net/http"

	"github.com/joelsaunders/bilbo-go/pkg/models"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/joelsaunders/bilbo-go/pkg/config"
	"github.com/joelsaunders/bilbo-go/pkg/repository"
)

func PostRoutes(postStore repository.PostStore, config *config.Config) *chi.Mux {
	router := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", config.JWTSecret, nil)

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Patch("/{postSlug}", NewPostHandler(postStore).updatePost())
		router.Post("/", NewPostHandler(postStore).createPost())
	})
	router.Get("/{postSlug}", NewPostHandler(postStore).retrievePost())
	router.Get("/", NewPostHandler(postStore).getPostList())
	return router
}

type PostHandler struct {
	store repository.PostStore
}

func NewPostHandler(store repository.PostStore) *PostHandler {
	return &PostHandler{store}
}

func (ph PostHandler) getPostList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := ph.store.List(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		render.JSON(w, r, posts)
	}
}

func (ph PostHandler) retrievePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postSlug := chi.URLParam(r, "postSlug")

		post, err := ph.store.GetBySlug(r.Context(), postSlug)
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
		}

		render.JSON(w, r, post)
	}
}

func (ph PostHandler) createPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newPost := postPayload{}

		if err := render.Bind(r, &newPost); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		post, err := ph.store.Create(r.Context(), newPost.Post)

		if err != nil {
			render.Render(w, r, ErrDatabase(err))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, post)
	}
}

func (ph PostHandler) updatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		postSlug := chi.URLParam(r, "postSlug")

		post, err := ph.store.GetBySlug(ctx, postSlug)
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
		}

		modifiedPost := postPayload{Post: post}

		if err := render.Bind(r, &modifiedPost); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		post, err = ph.store.Update(ctx, modifiedPost.Post)
		if err != nil {
			render.Render(w, r, ErrDatabase(err))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, post)
	}
}

type postPayload struct {
	*models.Post
}

func (pp *postPayload) Bind(r *http.Request) error {
	return nil
}
