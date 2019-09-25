package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joelsaunders/bilbo-go/api"
	config "github.com/joelsaunders/bilbo-go/config"
	"github.com/joelsaunders/bilbo-go/repository/postgres"

	"github.com/go-chi/chi"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configuration, err := config.NewConfig()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	db, err := postgres.NewDatabase(configuration.Postgres.URL, configuration.Postgres.DBPORT,
		configuration.Postgres.DBUSER, configuration.Postgres.DBPASSWORD, configuration.Postgres.DBNAME)
	defer db.Close()

	if err != nil {
		log.Panicln("Database error", err)
	}

	err = postgres.MigrateDatabase(db, "./migrations")
	if err != nil {
		log.Panicln("Migration error", err)
	}

	router := api.Routes(configuration, db)

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
