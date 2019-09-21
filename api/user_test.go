package api_test

import (
	"context"
	"fmt"
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

func init() {
	test_utils.SetUpTestDB("../migrations")
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")
}

func NewUserListRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/user", nil)
	return req
}

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

func TestListUsers(t *testing.T) {

	t.Run("Test Empty Response", func(t *testing.T) {
		cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
		db, err := sqlx.Open("txdb", cName)
		if err != nil {
			t.Fatal("could not open db")
		}
		defer db.Close()
		repo := repository.NewDB(db)

		request := NewUserListRequest()
		response := httptest.NewRecorder()
		server := api.NewServer(repo)

		server.Router.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]\n"

		if got != want {
			t.Fatalf("expected '%v' got '%v'", want, got)
		}
	})

	t.Run("Test with fake store", func(t *testing.T) {
		fakeDB := fakeDB{}
		fakeDB.users = make([]*models.User, 0)
		server := api.NewServer(fakeDB)

		request := NewUserListRequest()
		response := httptest.NewRecorder()

		server.Router.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]\n"

		if got != want {
			t.Fatalf("expected '%s' got '%s'", want, got)
		}

	})

}
