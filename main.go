package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	api "github.com/joelsaunders/bilbo-go/api"
	config "github.com/joelsaunders/bilbo-go/config"

	"github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/middleware"
	render "github.com/go-chi/render"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/user", api.UserRoutes(configuration))
	})

	return router
}

func main() {
	configuration, err := config.NewConfig()
	defer configuration.Database.Close()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	err = config.MigrateDatabase(configuration.Database)
	if err != nil {
		log.Panicln("Migration error", err)
	}

	router := routes(configuration)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Println("Serving application at PORT :" + configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configuration.PORT), router))
}
