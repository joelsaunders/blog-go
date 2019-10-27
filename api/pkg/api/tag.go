package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository"
)

func TagRoutes(store repository.TagStore, config *config.Config) *chi.Mux {
	router := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", config.JWTSecret, nil)

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Delete("/{tagID}", newTagHandler(store).deleteTag())
		router.Patch("/{tagID}", newTagHandler(store).updateTag())
		router.Post("/", newTagHandler(store).createTag())
	})

	router.Get("/", newTagHandler(store).listTags())
	return router
}

type tagHandler struct {
	store repository.TagStore
}

func newTagHandler(store repository.TagStore) *tagHandler {
	return &tagHandler{store: store}
}

type tagPayload struct {
	*models.Tag
}

func (tp *tagPayload) Bind(_ *http.Request) error {
	return nil
}

func (th tagHandler) listTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := th.store.List(r.Context())
		HandleApiErr(err, ErrDatabase, w, r)
		render.JSON(w, r, tags)
	}
}

func (th tagHandler) createTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newTagPayload := tagPayload{}
		err := render.Bind(r, &newTagPayload)
		if err != nil {
			HandleApiErr(err, ErrInvalidRequest, w, r)
			return
		}
		tag, err := th.store.Create(r.Context(), newTagPayload.Tag)
		if err != nil {
			HandleApiErr(err, ErrDatabase, w, r)
			return
		}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, tag)
	}
}

func (th tagHandler) updateTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(chi.URLParam(r, "tagID"))
		if err != nil {
			HandleApiErr(err, ErrInvalidRequest, w, r)
			return
		}
		// Should get by id here first and bind to existing obj
		newTagPayload := tagPayload{}
		err = render.Bind(r, &newTagPayload)
		if err != nil {
			HandleApiErr(err, ErrInvalidRequest, w, r)
			return
		}
		newTagPayload.Tag.ID = tagID

		tag, err := th.store.Update(r.Context(), newTagPayload.Tag)
		if err != nil {
			HandleApiErr(err, ErrDatabase, w, r)
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, tag)
	}
}

func (th tagHandler) deleteTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(chi.URLParam(r, "tagID"))
		if err != nil {
			HandleApiErr(err, ErrInvalidRequest, w, r)
			return
		}
		err = th.store.DeleteByID(r.Context(), tagID)
		HandleApiErr(err, ErrDatabase, w, r)
		render.NoContent(w, r)
	}
}
