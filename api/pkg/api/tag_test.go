package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joelsaunders/blog-go/api/pkg/api"
	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/test_utils"
)

type fakeTagDB struct {
	tags []*models.Tag
}

func (ft fakeTagDB) List(_ context.Context) ([]*models.Tag, error) {
	return ft.tags, nil
}
func (ft *fakeTagDB) Create(_ context.Context, tag *models.Tag) (*models.Tag, error) {
	tag.ID = len(ft.tags) + 1
	ft.tags = append(ft.tags, tag)
	return tag, nil
}
func (ft *fakeTagDB) Update(_ context.Context, tag *models.Tag) (*models.Tag, error) {
	dbTag, err := ft.getTagByID(tag.ID)
	if err != nil {
		return nil, err
	}
	dbTag.Name = tag.Name
	return dbTag, nil
}
func (ft *fakeTagDB) DeleteByID(_ context.Context, ID int) error {
	var newSlice = make([]*models.Tag, 0)
	for _, tag := range ft.tags {
		if tag.ID != ID {
			newSlice = append(newSlice, tag)
		}
	}
	ft.tags = newSlice
	return nil
}

func (ft fakeTagDB) getTagByID(ID int) (*models.Tag, error) {
	for _, tag := range ft.tags {
		if tag.ID == ID {
			return tag, nil
		}
	}
	return nil, fmt.Errorf("could not find tag with id %d", ID)
}

func TestTagAPI(t *testing.T) {
	t.Run("Test tag list", func(t *testing.T) {
		configuration, _ := config.NewConfig()

		testTag := models.Tag{ID: 1, Name: "test tag 1"}
		tagStore := fakeTagDB{[]*models.Tag{&testTag}}

		server := api.TagRoutes(&tagStore, configuration)
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusOK, t)

		expectedTag, _ := json.Marshal(tagStore.tags)
		test_utils.AssertEqualJSON(response.Body.String(), string(expectedTag), t)
	})

	t.Run("Test tag delete via api", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		testTag := models.Tag{ID: 1, Name: "test tag 1"}
		tagStore := fakeTagDB{[]*models.Tag{&testTag}}

		server := api.TagRoutes(&tagStore, configuration)
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%d", testTag.ID), nil)
		test_utils.AddAuthHeader(request, 1, "joel", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusNoContent, t)
		if len(tagStore.tags) != 0 {
			t.Fatalf("tag was not deleted")
		}
	})

	t.Run("Test tag create via api", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		tagStore := fakeTagDB{[]*models.Tag{}}

		toCreateTag := models.Tag{Name: "new tag name"}
		toCreateTagJSON, _ := json.Marshal(toCreateTag)

		server := api.TagRoutes(&tagStore, configuration)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(toCreateTagJSON))
		request.Header.Set("Content-Type", "application/json")
		test_utils.AddAuthHeader(request, 1, "joel", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusCreated, t)
		if len(tagStore.tags) != 1 {
			t.Errorf("tag store has more that one item in it, it has %d", len(tagStore.tags))
		}

		if tagStore.tags[0].Name != toCreateTag.Name {
			t.Errorf("created tag name (%s) is not correct (%s)", tagStore.tags[0].Name, toCreateTag.Name)
		}
	})

	t.Run("Test tag update via api", func(t *testing.T) {
		configuration, _ := config.NewConfig()

		testTag := models.Tag{ID: 1, Name: "test tag 1"}
		tagStore := fakeTagDB{[]*models.Tag{&testTag}}
		server := api.TagRoutes(&tagStore, configuration)
		response := httptest.NewRecorder()

		toUpdateTag := models.Tag{Name: "new tag name"}
		toUpdateTagJSON, _ := json.Marshal(toUpdateTag)
		request, _ := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/%d", testTag.ID),
			bytes.NewReader(toUpdateTagJSON),
		)
		request.Header.Set("Content-Type", "application/json")
		test_utils.AddAuthHeader(request, 1, "joel", configuration.JWTSecret)
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusOK, t)

		if tagStore.tags[0].Name != toUpdateTag.Name {
			t.Errorf("updated tag name (%s) is not correct (%s)", tagStore.tags[0].Name, toUpdateTag.Name)
		}
	})
}
