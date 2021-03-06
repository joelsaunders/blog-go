package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/joelsaunders/blog-go/api/pkg/auth"
	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository"
)

// UserRoutes is the router mux that conains routes for updating/modifying users and logging in
func UserRoutes(userStore repository.UserStore, config *config.Config) *chi.Mux {
	router := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", config.JWTSecret, nil)

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Get("/{userID}", newUserHandler(userStore).retrieveUser())
		router.Patch("/change-password", newUserHandler(userStore).updateUserPassword())
		router.Get("/", newUserHandler(userStore).getUserList())
		router.Post("/", newUserHandler(userStore).createUser())
	})

	router.Post("/login", newUserHandler(userStore).loginUser(config.JWTSecret))
	return router
}

type userHandler struct {
	store repository.UserStore
}

func newUserHandler(userStore repository.UserStore) *userHandler {
	return &userHandler{userStore}
}

func (uh userHandler) loginUser(jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials := userCredentialsPayload{}

		if err := render.Bind(r, &credentials); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		userID, err := auth.CheckCredentials(r.Context(), credentials.Email, credentials.Password, uh.store)

		if err != nil {
			render.Render(w, r, ErrAuthenication(err))
			return
		}

		token, err := auth.GenerateToken(userID, credentials.Email, jwtKey)
		if err != nil {
			render.Render(w, r, ErrTokenCreation(err))
			return
		}

		tokenResponse := userTokenResponse{token}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, tokenResponse)
	}
}

func (uh userHandler) retrieveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}

		user, err := uh.store.GetByID(r.Context(), userID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		render.JSON(w, r, user)
	}
}

func (uh userHandler) getUserList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := uh.store.List(r.Context(), 10)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		render.JSON(w, r, users)
	}
}

func (uh userHandler) updateUserPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, _ := jwtauth.FromContext(r.Context())
		userID := int(claims["id"].(float64))

		user, err := uh.store.GetByID(ctx, userID)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		passwordChangePayload := &PasswordChangePayload{}
		if err := render.Bind(r, passwordChangePayload); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		user.Password = passwordChangePayload.Password
		user, err = uh.store.Update(ctx, user)

		if err != nil {
			render.Render(w, r, ErrDatabase(err))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
	}
}

func (uh userHandler) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newUser := newUserPayload{}
		if err := render.Bind(r, &newUser); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		user, err := uh.store.Create(r.Context(), newUser.NewUser)

		if err != nil {
			render.Render(w, r, ErrDatabase(err))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, user)
	}
}

type userTokenResponse struct {
	Token string `json:"token"`
}

type userCredentialsPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ucp *userCredentialsPayload) Bind(r *http.Request) error {
	return nil
}

type newUserPayload struct {
	*models.NewUser
}

func (nu *newUserPayload) Bind(r *http.Request) error {
	nu.Password = auth.HashPassword(nu.Password)
	return nil
}

type PasswordChangePayload struct {
	Password string `json:"password"`
}

func (pcp *PasswordChangePayload) Bind(r *http.Request) error {
	pcp.Password = auth.HashPassword(pcp.Password)
	return nil
}
