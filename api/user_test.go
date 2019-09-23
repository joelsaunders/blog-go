package api_test

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/bilbo-go/repository"
	"github.com/joelsaunders/bilbo-go/test_utils"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joelsaunders/bilbo-go/api"
	"github.com/joelsaunders/bilbo-go/models"
	_ "github.com/lib/pq"
)

type fakeDB struct {
	users []*models.User
}

func (f fakeDB) Users() repository.UserStore {
	return &fakeUserDB{f.users}
}

type fakeUserDB struct {
	users []*models.User
}

func (fu fakeUserDB) List(ctx context.Context, num int) ([]*models.User, error) {
	return fu.users, nil
}

func (fu fakeUserDB) Create(ctx context.Context, user *models.NewUser) (*models.User, error) {
	return &models.User{
		ID:       rand.Intn(100),
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func TestUsersAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	test_utils.SetUpTestDB("../migrations")
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")

	t.Run("Test Empty Response", func(t *testing.T) {
		cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
		db, err := sqlx.Open("txdb", cName)
		if err != nil {
			t.Fatal("could not open db")
		}
		defer db.Close()
		repo := repository.NewDB(db)

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/user", nil)
		response := httptest.NewRecorder()
		server := api.NewServer(repo)

		server.Router.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]\n"

		if got != want {
			t.Fatalf("expected '%v' got '%v'", want, got)
		}
	})
}

func TestUsersAPI(t *testing.T) {
	t.Run("Test users list empty", func(t *testing.T) {
		fakeDB := fakeDB{}
		fakeDB.users = make([]*models.User, 0)
		server := api.UserRoutes(fakeDB.Users())

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]\n"

		if got != want {
			t.Fatalf("expected '%s' got '%s'", want, got)
		}

	})

	t.Run("Test users list results", func(t *testing.T) {
		fakeDB := fakeDB{}
		fakeDB.users = []*models.User{
			&models.User{
				ID:       1,
				Email:    "joel.st.saunders@gmail.com",
				Password: "helloooooo",
			},
		}
		server := api.UserRoutes(fakeDB.Users())

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[{\"ID\":1,\"Email\":\"joel.st.saunders@gmail.com\",\"Password\":\"helloooooo\"}]\n"

		if got != want {
			t.Fatalf("expected '%s' got '%s'", want, got)
		}

	})
}
