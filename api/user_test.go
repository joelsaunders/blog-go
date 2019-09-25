package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/bilbo-go/auth"
	"github.com/joelsaunders/bilbo-go/config"
	"github.com/joelsaunders/bilbo-go/test_utils"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joelsaunders/bilbo-go/api"
	"github.com/joelsaunders/bilbo-go/models"
	"github.com/joelsaunders/bilbo-go/repository"
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

func (fu *fakeUserDB) Create(ctx context.Context, user *models.NewUser) (*models.User, error) {
	userObj := &models.User{
		ID:       len(fu.users),
		Email:    user.Email,
		Password: user.Password,
	}
	fu.users = append(fu.users, userObj)
	return userObj, nil
}

func (fu fakeUserDB) GetByID(ctx context.Context, id int) (*models.User, error) {
	for _, user := range fu.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("no user found")
}

func (fu fakeUserDB) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range fu.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("no user found")
}

func assertResponseCode(got int, want int, t *testing.T) {
	if got != want {
		t.Fatalf("got response code %d want %d", got, want)
	}
}

func assertBody(got, want string, t *testing.T) {
	if got != want {
		t.Fatalf("expected body '%v' got '%v'", want, got)
	}
}

func assertEqualJSON(s1, s2 string, t *testing.T) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		t.Fatalf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		t.Fatalf("Error mashalling string 2 :: %s", err.Error())
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatalf("json %s and %s are not equal", s1, s2)
	}
}

func TestUsersAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	test_utils.SetUpTestDB("../migrations")
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")

	t.Run("Test Empty Response", func(t *testing.T) {
		configuration, _ := config.NewConfig()

		cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
		db, err := sqlx.Open("txdb", cName)
		if err != nil {
			t.Fatal("could not open db")
		}
		defer db.Close()

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/user", nil)
		addAuthHeader(request, "fakeemailthatdoesnotexist@gmail.com", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server := api.Routes(configuration, db)

		server.ServeHTTP(response, request)

		assertBody(response.Body.String(), "[]\n", t)
	})
}

func addAuthHeader(request *http.Request, email string, secret []byte) {
	// set the correct token header
	authToken, _ := auth.GenerateToken(email, secret)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
}

func TestUsersAPI(t *testing.T) {
	t.Run("Test users list results", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		testUser := models.User{
			ID:       1,
			Email:    "joel.st.saunders@gmail.com",
			Password: auth.HashPassword("helloooooo"),
		}

		fakeDB := fakeDB{}
		fakeDB.users = []*models.User{&testUser}
		server := api.UserRoutes(fakeDB.Users(), configuration)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		addAuthHeader(request, testUser.Email, configuration.JWTSecret)
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[{\"id\":1,\"email\":\"joel.st.saunders@gmail.com\"}]\n"
		assertBody(got, want, t)
	})

	t.Run("Test user retrieve", func(t *testing.T) {
		configuration, _ := config.NewConfig()

		testUser := models.User{
			ID:       1,
			Email:    "joel.st.saunders@gmail.com",
			Password: auth.HashPassword("helloooooo"),
		}

		fakeDB := fakeDB{}
		fakeDB.users = []*models.User{&testUser}
		server := api.UserRoutes(fakeDB.Users(), configuration)

		request, _ := http.NewRequest(http.MethodGet, "/1", nil)
		response := httptest.NewRecorder()
		addAuthHeader(request, testUser.Email, configuration.JWTSecret)
		server.ServeHTTP(response, request)

		assertResponseCode(response.Code, http.StatusOK, t)
		expectedUser, _ := json.Marshal(testUser)
		assertEqualJSON(response.Body.String(), string(expectedUser), t)
	})

	t.Run("Test create user", func(t *testing.T) {
		configuration, _ := config.NewConfig()

		fakeDB := fakeDB{}
		server := api.UserRoutes(fakeDB.Users(), configuration)

		newUser := models.NewUser{Email: "newperson@new.com", Password: "Password"}
		newUserJSON, _ := json.Marshal(newUser)

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(newUserJSON))
		request.Header.Set("Content-Type", "application/json")
		addAuthHeader(request, "pretendemail@test.com", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseCode(response.Code, http.StatusCreated, t)

		got := response.Body.String()
		want := fmt.Sprintf("{\"id\":0,\"email\":\"%s\"}\n", newUser.Email)
		assertBody(got, want, t)
	})
}

func TestLogin(t *testing.T) {
	configuration, _ := config.NewConfig()

	testUser := models.User{
		ID:       1,
		Email:    "joel.st.saunders@gmail.com",
		Password: auth.HashPassword("helloooooo"),
	}

	fakeDB := fakeDB{}
	fakeDB.users = []*models.User{
		&testUser,
	}
	server := api.UserRoutes(fakeDB.Users(), configuration)

	t.Run("login incorrect credentials", func(t *testing.T) {
		credentials := map[string]string{"email": "joel.st.saunders@gmail.com", "password": "Password"}
		credentialsJSON, _ := json.Marshal(credentials)

		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(credentialsJSON))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseCode(response.Code, http.StatusUnauthorized, t)
	})

	t.Run("login correct credentials", func(t *testing.T) {
		credentials := map[string]string{"email": "joel.st.saunders@gmail.com", "password": "helloooooo"}
		credentialsJSON, _ := json.Marshal(credentials)

		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(credentialsJSON))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		expectedToken, _ := auth.GenerateToken(testUser.Email, configuration.JWTSecret)

		assertResponseCode(response.Code, http.StatusOK, t)
		assertBody(response.Body.String(), fmt.Sprintf("{\"token\":\"%s\"}\n", expectedToken), t)
	})
}
