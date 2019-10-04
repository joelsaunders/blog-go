package user_test

import (
	"context"
	"testing"

	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository/user"
	"github.com/joelsaunders/blog-go/api/test_utils"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	test_utils.SetUpTestDB("../../../migrations")
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")

	t.Run("Test Get Users", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()

		_, err := db.Exec("insert into users (email, password) values ('joel', 'mpassword')")

		if err != nil {
			t.Fatalf("could not insert user: %s", err)
		}

		userStore := user.PGUserStore{db}
		ctx := context.Background()

		users, err := userStore.List(ctx, 2)

		if err != nil {
			t.Fatalf("could not return users: %s", err)
		}

		if len(users) != 1 {
			t.Fatalf("expected one user but got %d", len(users))
		}
	})

	t.Run("Test create user", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()
		userStore := user.PGUserStore{db}
		ctx := context.Background()

		newUser := &models.NewUser{
			Email:    "joel.st.saunders@gmail.com",
			Password: "password",
		}

		user, err := userStore.Create(ctx, newUser)

		if err != nil {
			t.Fatalf("could not create user: %s", err)
		}

		if user.Email != newUser.Email {
			t.Fatalf("user %v does not same data as inserted user: %v", user, newUser)
		}
	})

}
