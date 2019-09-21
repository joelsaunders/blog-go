package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/bilbo-go/repository/user"
	"github.com/joelsaunders/bilbo-go/test_utils"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func init() {
	test_utils.SetUpTestDB("../../migrations")
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")
}

func TestGetUsers(t *testing.T) {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, err := sqlx.Open("txdb", cName)
	if err != nil {
		t.Fatal("could not open db")
	}
	defer db.Close()

	_, err = db.Exec("insert into users (email, password) values ('joel', 'mpassword')")

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
}
