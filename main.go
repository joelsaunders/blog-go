package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joelsaunders/bilbo-go/repository"

	"github.com/joelsaunders/bilbo-go/api"
	config "github.com/joelsaunders/bilbo-go/config"

	"github.com/go-chi/chi"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configuration, err := config.NewConfig()
	defer configuration.Database.Close()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	err = config.MigrateDatabase(configuration.Database, "./migrations")
	if err != nil {
		log.Panicln("Migration error", err)
	}

	db := repository.NewDB(configuration.Database)
	server := api.NewServer(db)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(server.Router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Println("Serving application at PORT :" + configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configuration.PORT), server.Router))
}
