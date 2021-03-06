package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joelsaunders/blog-go/api/pkg/api"
	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/test_utils"
)

type fakePostDB struct {
	posts []*models.Post
}

func contains(item string, container []string) bool {
	for _, x := range container {
		if x == item {
			return true
		}
	}
	return false
}

func (fp fakePostDB) List(_ context.Context, filters map[string][]string) ([]*models.Post, error) {
	if len(filters) == 0 {
		return fp.posts, nil
	}
	filteredPosts := fp.posts
	n := 0
	for _, post := range filteredPosts {
		for k, v := range filters {
			switch k {
			case "tag_name":
				{
					for _, tag := range post.Tags {
						if contains(tag, v) {
							filteredPosts[n] = post
							n++
						}
					}
				}
			}
		}
	}
	return filteredPosts[:n], nil
}

func (fp fakePostDB) GetByID(_ context.Context, id int) (*models.Post, error) {
	for _, post := range fp.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (fp fakePostDB) GetBySlug(_ context.Context, slug string) (*models.Post, error) {
	for _, post := range fp.posts {
		if post.Slug == slug {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (fp *fakePostDB) Create(_ context.Context, post *models.Post) (*models.Post, error) {
	post.Created = time.Now().UTC()
	post.Modified = time.Now().UTC()
	post.ID = len(fp.posts) + 1

	fp.posts = append(fp.posts, post)
	return post, nil
}

func (fp *fakePostDB) Update(ctx context.Context, post *models.Post) (*models.Post, error) {
	dbPost, err := fp.GetByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	dbPost.Modified = time.Now().UTC()
	dbPost.Slug = post.Slug
	dbPost.Title = post.Title
	dbPost.Body = post.Body
	dbPost.Description = post.Description
	dbPost.Picture = post.Picture
	dbPost.Published = post.Published

	return dbPost, nil
}

func (fp *fakePostDB) DeleteBySlug(_ context.Context, postSlug string) error {
	var newSlice = make([]*models.Post, 0)
	for _, post := range fp.posts {
		if post.Slug != postSlug {
			newSlice = append(newSlice, post)
		}
	}
	fp.posts = newSlice
	return nil
}

func TestPostAPI(t *testing.T) {
	t.Run("Test post delete", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		testPost := models.Post{
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    1,
			AuthorEmail: "joel",
			Tags:        []string{"hello1"},
		}

		postStore := fakePostDB{[]*models.Post{&testPost}}
		server := api.PostRoutes(&postStore, configuration)
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", testPost.Slug), nil)
		test_utils.AddAuthHeader(request, 1, "joel", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusNoContent, t)
		if len(postStore.posts) != 0 {
			t.Fatalf("post was not deleted")
		}
	})

	t.Run("Test post list", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		testPost := models.Post{
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    1,
			AuthorEmail: "joel",
			Tags:        []string{"hello1"},
		}

		postStore := fakePostDB{[]*models.Post{&testPost}}
		server := api.PostRoutes(&postStore, configuration)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusOK, t)

		expectedPost, _ := json.Marshal(postStore.posts)
		test_utils.AssertEqualJSON(response.Body.String(), string(expectedPost), t)
	})

	t.Run("Test post list with filter", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		testPost1 := models.Post{
			ID:   1,
			Tags: []string{"hello1"},
		}
		testPost2 := models.Post{
			ID:   2,
			Tags: []string{"hello2"},
		}

		postStore := fakePostDB{[]*models.Post{&testPost1, &testPost2}}
		server := api.PostRoutes(&postStore, configuration)

		request, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/?tag_name=%s", testPost1.Tags[0]),
			nil,
		)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusOK, t)

		expectedPost, _ := json.Marshal([]*models.Post{&testPost1})
		test_utils.AssertEqualJSON(response.Body.String(), string(expectedPost), t)
	})

	t.Run("Test post retrieve", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		testPost := models.Post{
			ID:          1,
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    1,
		}

		postStore := fakePostDB{[]*models.Post{&testPost}}
		server := api.PostRoutes(&postStore, configuration)
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testPost.Slug), nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusOK, t)
		expectedPost, _ := json.Marshal(testPost)
		test_utils.AssertEqualJSON(response.Body.String(), string(expectedPost), t)
	})

	t.Run("Test post create", func(t *testing.T) {
		configuration, _ := config.NewConfig()
		postStore := fakePostDB{}
		server := api.PostRoutes(&postStore, configuration)

		newPost := models.Post{
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    1,
		}
		newPostJSON, _ := json.Marshal(newPost)

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(newPostJSON))
		request.Header.Set("Content-Type", "application/json")
		test_utils.AddAuthHeader(request, 1, "joel", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusCreated, t)
		// get created post from db
		createdPost, err := postStore.GetBySlug(context.Background(), newPost.Slug)
		if err != nil {
			t.Fatalf("created post cannot be found by slug")
		}

		// set auto fields
		newPost.ID = createdPost.ID
		newPost.Created = createdPost.Created
		newPost.Modified = createdPost.Modified
		expectedPost, _ := json.Marshal(newPost)
		test_utils.AssertEqualJSON(response.Body.String(), string(expectedPost), t)
	})

	t.Run("Test post update", func(t *testing.T) {
		configuration, _ := config.NewConfig()

		// initial post in store
		existingPost := models.Post{
			ID:          1,
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    1,
		}
		postStore := fakePostDB{[]*models.Post{&existingPost}}
		server := api.PostRoutes(&postStore, configuration)

		modifiedPostData := models.Post{
			Slug:        "new slug",
			Title:       "new title",
			Body:        "new body",
			Picture:     "new picture",
			Description: "new description",
			Published:   false,
			AuthorID:    2,
		}
		modifiedPostJSON, _ := json.Marshal(modifiedPostData)

		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/%s", existingPost.Slug), bytes.NewReader(modifiedPostJSON))
		request.Header.Set("Content-Type", "application/json")
		test_utils.AddAuthHeader(request, 1, "joel", configuration.JWTSecret)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		test_utils.AssertResponseCode(response.Code, http.StatusOK, t)

		// get modified post from db
		modifiedPostDB, err := postStore.GetByID(context.Background(), existingPost.ID)
		if err != nil {
			t.Fatalf("created post cannot be found by id")
		}

		// set auto fields
		modifiedPostData.ID = modifiedPostDB.ID
		modifiedPostData.Created = modifiedPostDB.Created
		modifiedPostData.Modified = modifiedPostDB.Modified
		expectedPost, _ := json.Marshal(modifiedPostData)
		test_utils.AssertEqualJSON(response.Body.String(), string(expectedPost), t)

		// assert that modified date has changed
		if modifiedPostData.Modified == modifiedPostData.Created {
			t.Errorf("post modified date has not changed from the created date")
		}
	})
}
