package post_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository/post"
	"github.com/joelsaunders/blog-go/api/test_utils"
)

func insertUser(email string, db *sqlx.DB, t *testing.T) (userID int) {
	err := db.QueryRowx(
		"insert into users (email, password) values ($1, 'mpassword') RETURNING id;",
		email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("could not insert user: %s", err)
	}
	return
}

func insertPost(post *models.Post, db *sqlx.DB, t *testing.T) (postID int) {
	err := db.QueryRowx(
		fmt.Sprintf(
			`insert into posts (
				created,
				modified,
				slug,
				title,
				body,
				picture,
				description,
				published,
				author_id
			) values (
				'%s',
				'%s',
				'%s',
				'%s',
				'%s',
				'%s',
				'%s',
				%v,
				%v
			) RETURNING id;`,
			post.Created.Format("2006-01-02 15:04:05-0700"),
			post.Modified.Format("2006-01-02 15:04:05-0700"),
			post.Slug,
			post.Title,
			post.Body,
			post.Picture,
			post.Description,
			post.Published,
			post.AuthorID,
		),
	).Scan(&postID)

	if err != nil {
		t.Fatalf("could not insert post: %s", err)
	}
	return
}

func assertPostEqual(post, expectedPost *models.Post, t *testing.T) {
	if !cmp.Equal(post, expectedPost) {
		t.Fatalf("post \n%v\n is not equal to expected \n%v\n", post, expectedPost)
	}
}

func TestPosts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	test_utils.SetUpTestDB("../../../migrations")
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")

	t.Run("Test List Posts", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()

		userID := insertUser("joel", db, t)

		testPost := models.Post{
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    userID,
		}

		insertPost(&testPost, db, t)

		postStore := post.PGPostStore{db}
		ctx := context.Background()

		posts, err := postStore.List(ctx)

		if err != nil {
			t.Fatalf("could not return posts: %s", err)
		}

		if len(posts) != 1 {
			t.Fatalf("expected one post but got %d", len(posts))
		}

		posts[0].ID = 0
		assertPostEqual(posts[0], &testPost, t)
	})

	t.Run("test get post", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()

		userID := insertUser("joel", db, t)

		testPost := models.Post{
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    userID,
		}

		insertPost(&testPost, db, t)

		postStore := post.PGPostStore{db}
		ctx := context.Background()

		post, err := postStore.GetBySlug(ctx, testPost.Slug)
		testPost.ID = post.ID

		if err != nil {
			t.Fatalf("could not return posts: %s", err)
		}

		assertPostEqual(post, &testPost, t)
	})

	t.Run("test create post", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()

		userID := insertUser("joel", db, t)

		testPost := models.Post{
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    userID,
		}

		postStore := post.PGPostStore{db}
		ctx := context.Background()

		post, err := postStore.Create(ctx, &testPost)

		if err != nil {
			t.Fatalf("could not create post: %s", err)
		}
		// set auto fields
		testPost.Created = post.Created
		testPost.Modified = post.Modified
		testPost.ID = post.ID
		assertPostEqual(post, &testPost, t)
	})

	t.Run("test update post", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()

		userID := insertUser("joel", db, t)
		user2ID := insertUser("jo", db, t)

		testPost := models.Post{
			Created:     time.Now().Round(time.Second).UTC(),
			Modified:    time.Now().Round(time.Second).UTC(),
			Slug:        "test slug",
			Title:       "test title",
			Body:        "test body",
			Picture:     "test picture",
			Description: "test description",
			Published:   true,
			AuthorID:    userID,
		}

		postID := insertPost(&testPost, db, t)

		postStore := post.PGPostStore{db}
		ctx := context.Background()

		// modify all fields on the post
		testPost.ID = postID
		testPost.Slug = "new slug"
		testPost.Title = "new title"
		testPost.Body = "new body"
		testPost.Picture = "new picture"
		testPost.Description = "new description"
		testPost.Published = false
		testPost.AuthorID = user2ID

		updatedPost, err := postStore.Update(ctx, &testPost)
		if err != nil {
			t.Fatalf("could not update post: %s", err)
		}
		// post will have been modified
		testPost.Modified = updatedPost.Modified
		assertPostEqual(updatedPost, &testPost, t)
	})
}
