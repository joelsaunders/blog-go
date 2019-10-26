package test_utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/joelsaunders/blog-go/api/pkg/auth"
)

func AssertResponseCode(got int, want int, t *testing.T) {
	if got != want {
		t.Fatalf("got response code %d want %d", got, want)
	}
}

func AssertEqualJSON(s1, s2 string, t *testing.T) {
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

func AddAuthHeader(request *http.Request, id int, email string, secret []byte) {
	// set the correct token header
	authToken, _ := auth.GenerateToken(id, email, secret)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
}

func InsertUser(email string, db *sqlx.DB, t *testing.T) (userID int) {
	err := db.QueryRowx(
		"insert into users (email, password) values ($1, 'mpassword') RETURNING id;",
		email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("could not insert user: %s", err)
	}
	return
}

func CreateTag(tagName string, db *sqlx.DB, t *testing.T) int {
	var tagID int

	err := db.QueryRowx(
		fmt.Sprintf(
			`insert into tags (
				name
			) values (
				'%s'
			) returning id;`,
			tagName,
		),
	).Scan(&tagID)

	if err != nil {
		t.Errorf("could not add tag %s because of %s", tagName, err)
	}

	return tagID
}

func AddTag(postID int, tagName string, db *sqlx.DB, t *testing.T) {
	tagID := CreateTag(tagName, db, t)

	var relationID int
	err := db.QueryRowx(
		fmt.Sprintf(
			`insert into posttags (
				tag_id,
				post_id
			) values (
				%d,
				%d 
			) returning id;`,
			tagID,
			postID,
		),
	).Scan(&relationID)

	if err != nil {
		t.Errorf("could not relate post %d to tag %d because of %s", postID, tagID, err)
	}
}
