package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository"
)

// PostRoutes returns a router mux that containes routes for modifying/retrieving post objects
func PostRoutes(postStore repository.PostStore, config *config.Config) *chi.Mux {
	router := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", config.JWTSecret, nil)

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Patch("/{postSlug}", newPostHandler(postStore).updatePost())
		router.Delete("/{postSlug}", newPostHandler(postStore).deletePost())
		router.Post("/", newPostHandler(postStore).createPost())
	})
	router.Get("/{postSlug}", newPostHandler(postStore).retrievePost())
	router.Get("/", newPostHandler(postStore).getPostList())
	return router
}

type postHandler struct {
	store repository.PostStore
}

func newPostHandler(store repository.PostStore) *postHandler {
	return &postHandler{store}
}

func (ph postHandler) deletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postSlug := chi.URLParam(r, "postSlug")

		err := ph.store.DeleteBySlug(r.Context(), postSlug)
		if err != nil {
			log.Println(err)
			render.Render(w, r, ErrDatabase(err))
		}
		render.NoContent(w, r)
	}
}

func (ph postHandler) getPostList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := ph.store.List(r.Context())
		if err != nil {
			log.Println(err)
			render.Render(w, r, ErrDatabase(err))
		}
		render.JSON(w, r, posts)
	}
}

func (ph postHandler) retrievePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postSlug := chi.URLParam(r, "postSlug")

		post, err := ph.store.GetBySlug(r.Context(), postSlug)
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
		}

		render.JSON(w, r, post)
	}
}

func (ph postHandler) createPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newPost := postPayload{}

		if err := render.Bind(r, &newPost); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// set the author id automatically for creation of posts
		_, claims, _ := jwtauth.FromContext(r.Context())
		userID := int(claims["id"].(float64))
		newPost.Post.AuthorID = userID

		post, err := ph.store.Create(r.Context(), newPost.Post)

		if err != nil {
			render.Render(w, r, ErrDatabase(err))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, post)
	}
}

func (ph postHandler) updatePost() http.HandlerFunc {
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
