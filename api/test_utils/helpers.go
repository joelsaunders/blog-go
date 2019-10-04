package test_utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

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

func AddAuthHeader(request *http.Request, email string, secret []byte) {
	// set the correct token header
	authToken, _ := auth.GenerateToken(email, secret)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
}
