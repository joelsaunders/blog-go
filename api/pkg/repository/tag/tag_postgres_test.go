package tag_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository/tag"
	"github.com/joelsaunders/blog-go/api/test_utils"
)

func createTag(name string, db *sqlx.DB, t *testing.T) *models.Tag {
	tagID := test_utils.CreateTag(name, db, t)
	return &models.Tag{ID: tagID, Name: name}
}

func cleanupDefaultTags(db *sqlx.DB) {
	db.MustExec("DELETE FROM tags")
}

func TestTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	err := test_utils.SetUpTestDB("../../../migrations")
	if err != nil {
		panic(err)
	}
	txdb.Register("txdb", "postgres", "host=localhost port=15432 user=root password=root dbname=test sslmode=disable")

	t.Run("List tags", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()
		// migrations make some tags by default so remove them for test
		cleanupDefaultTags(db)

		// create tags for listing
		tagOne := createTag("tag 1", db, t)
		tagTwo := createTag("tag 2", db, t)
		expectedTagList := []*models.Tag{tagOne, tagTwo}

		tagStore := tag.PGTagStore{DB: db}
		tags, err := tagStore.List(context.Background())

		if err != nil {
			t.Fatalf("could not return tags: %s", err)
		}

		if !cmp.Equal(tags, expectedTagList) {
			t.Fatalf("tag list %#v not equal to expected %#v", tags, expectedTagList)
		}
	})

	t.Run("Create Tag", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()
		// migrations make some tags by default so remove them for test
		cleanupDefaultTags(db)
		tagStore := tag.PGTagStore{DB: db}

		toInsertTag := models.Tag{Name: "new tag"}
		insertedTag, err := tagStore.Create(context.Background(), &toInsertTag)
		if err != nil {
			t.Fatalf("could not insert tag: %s", err)
		}

		if insertedTag.Name != toInsertTag.Name {
			t.Fatalf("name of tag returned tag (%s) is not correct (%s)", insertedTag.Name, toInsertTag.Name)
		}

		if insertedTag.ID == 0 {
			t.Fatalf("id of inserted tag is 0")
		}

		existingTags, err := tagStore.List(context.Background())
		if err != nil {
			t.Fatalf("could not list after create")
		}

		if len(existingTags) != 1 {
			t.Fatalf("more than one tags were created: %d", len(existingTags))
		}

		if !cmp.Equal(existingTags[0], insertedTag) {
			t.Fatalf(
				"listed tag (%v) inserted tag is not equal to inserted (%v)",
				*existingTags[0],
				*insertedTag,
			)
		}
	})

	t.Run("Update tag", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()
		// migrations make some tags by default so remove them for test
		cleanupDefaultTags(db)
		tagStore := tag.PGTagStore{DB: db}

		tagOne := createTag("tag 1", db, t)
		toUpdateTagOne := &models.Tag{ID: tagOne.ID, Name: "new name"}

		returnedTag, err := tagStore.Update(context.Background(), toUpdateTagOne)
		if err != nil {
			t.Fatalf("could not update tag: %s", err)
		}

		if !cmp.Equal(returnedTag, toUpdateTagOne) {
			t.Fatalf("returned tag (%v) not equal to input (%v)", *returnedTag, *toUpdateTagOne)
		}

		existingTags, err := tagStore.List(context.Background())
		if err != nil {
			t.Fatalf("could not list after update")
		}

		if len(existingTags) != 1 {
			t.Fatalf("more than one tags exist after update: %d", len(existingTags))
		}

		if !cmp.Equal(existingTags[0], toUpdateTagOne) {
			t.Fatalf(
				"listed tag (%v) inserted tag is not equal to inserted (%v)",
				*existingTags[0],
				*toUpdateTagOne,
			)
		}
	})

	t.Run("Delete tag", func(t *testing.T) {
		db := test_utils.OpenTransaction(t)
		defer db.Close()
		// migrations make some tags by default so remove them for test
		cleanupDefaultTags(db)
		tagStore := tag.PGTagStore{DB: db}

		tagOne := createTag("tag 1", db, t)

		err := tagStore.DeleteByID(context.Background(), tagOne.ID)
		if err != nil {
			t.Fatalf("could not delete tag: %s", err)
		}

		existingTags, err := tagStore.List(context.Background())
		if err != nil {
			t.Fatalf("could not list after delete")
		}

		if len(existingTags) != 0 {
			t.Fatalf("tags still exist after delete: %d", len(existingTags))
		}
	})
}
