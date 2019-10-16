package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/blog-go/api/pkg/api"
	"github.com/joelsaunders/blog-go/api/pkg/auth"
	config "github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository/post"
	"github.com/joelsaunders/blog-go/api/pkg/repository/postgres"
	"github.com/joelsaunders/blog-go/api/pkg/repository/user"

	"github.com/go-chi/chi"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func testDataSetup(db *sqlx.DB) {
	ctx := context.Background()

	newUser := &models.NewUser{
		Email:    "joel.st.saunders@gmail.com",
		Password: auth.HashPassword("password"),
	}
	userStore := &user.PGUserStore{DB: db}

	user, err := userStore.Create(ctx, newUser)

	if err != nil {
		log.Println("user could not be created")
		return
	}

	testPost := models.Post{
		Created:     time.Now().Round(time.Second).UTC(),
		Modified:    time.Now().Round(time.Second).UTC(),
		Slug:        "test slug",
		Title:       "test title",
		Body:        "test body",
		Picture:     "https://i.ytimg.com/vi/Vp7nW2SP6H8/maxresdefault.jpg",
		Description: "test description",
		Published:   true,
		AuthorID:    user.ID,
	}

	postStore := post.PGPostStore{DB: db}

	_, err = postStore.Create(ctx, &testPost)

	if err != nil {
		log.Panicln("could not create default post")
	}
}

func main() {
	configuration, err := config.NewConfig()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	db, err := postgres.NewDatabase(configuration.Postgres.URL, configuration.Postgres.DBPORT,
		configuration.Postgres.DBUSER, configuration.Postgres.DBPASSWORD, configuration.Postgres.DBNAME)

	if err != nil {
		log.Panicln("Database error", err)
	}
	defer db.Close()

	err = postgres.MigrateDatabase(db, "./migrations")
	if err != nil {
		log.Panicln("Migration error", err)
	}

	router := api.Routes(configuration, db)

	testDataSetup(db)

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
