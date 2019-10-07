package post

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/blog-go/api/pkg/models"
)

type PGPostStore struct {
	DB *sqlx.DB
}

type lastInsertDetails struct {
	ID   int    `db:"id"`
	Slug string `db:"slug"`
}

func (ps *PGPostStore) Create(ctx context.Context, post *models.Post) (*models.Post, error) {
	query := `INSERT INTO posts (
		slug,
		title,
		body,
		picture,
		description,
		published,
		author_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, slug`

	var lastInsert lastInsertDetails

	err := ps.DB.QueryRowxContext(
		ctx,
		query,
		post.Slug,
		post.Title,
		post.Body,
		post.Picture,
		post.Description,
		post.Published,
		post.AuthorID,
	).StructScan(&lastInsert)
	if err != nil {
		return nil, err
	}
	err = ps.addTags(ctx, lastInsert.ID, post.Tags)
	if err != nil {
		return nil, err
	}

	createdPost, err := ps.GetBySlug(ctx, lastInsert.Slug)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (ps *PGPostStore) Update(ctx context.Context, post *models.Post) (*models.Post, error) {
	var updateRow lastInsertDetails

	err := ps.DB.QueryRowxContext(
		ctx,
		`UPDATE posts SET 
			slug = $2,
			title = $3,
			body = $4,
			picture = $5,
			description = $6,
			published = $7,
			author_id = $8,
			modified = $9
		WHERE id=$1 RETURNING id, slug`,
		post.ID,
		post.Slug,
		post.Title,
		post.Body,
		post.Picture,
		post.Description,
		post.Published,
		post.AuthorID,
		time.Now().UTC(),
	).StructScan(&updateRow)

	if err != nil {
		return nil, err
	}

	// sort out tags by first getting them and then adding/removing as needed
	existingTags, err := ps.getTags(ctx, updateRow.ID)
	if err != nil {
		return nil, err
	}

	toAdd := difference(post.Tags, existingTags)
	fmt.Println(toAdd)
	err = ps.addTags(ctx, updateRow.ID, toAdd)
	toRemove := difference(existingTags, post.Tags)
	fmt.Println(toRemove)
	err = ps.removeTags(ctx, updateRow.ID, toRemove)

	if err != nil {
		return nil, err
	}

	insertedPost, err := ps.GetBySlug(ctx, updateRow.Slug)
	if err != nil {
		return nil, err
	}

	return insertedPost, nil
}

func (ps *PGPostStore) List(ctx context.Context) ([]*models.Post, error) {
	query := "SELECT * FROM posts"
	rows, err := ps.DB.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)

	for rows.Next() {
		var p models.Post
		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		tags, err := ps.getTags(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Tags = tags
		posts = append(posts, &p)
	}

	return posts, nil
}

func (ps *PGPostStore) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	post := models.Post{}
	err := ps.DB.Get(&post, "SELECT * FROM posts WHERE slug=$1", slug)
	if err != nil {
		return nil, err
	}

	tags, err := ps.getTags(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.Tags = tags
	return &post, nil
}

func contains(c []string, m string) bool {
	for _, a := range c {
		if a == m {
			return true
		}
	}
	return false
}

func difference(f, s []string) []string {
	var diff []string
	for _, fItem := range f {
		if !contains(s, fItem) {
			diff = append(diff, fItem)
		}
	}
	return diff
}

func (ps *PGPostStore) getTags(ctx context.Context, postID int) ([]string, error) {
	query := "SELECT name FROM tags INNER JOIN posttags pt ON pt.tag_id = tags.id WHERE pt.post_id = $1"
	rows, err := ps.DB.QueryxContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)

	for rows.Next() {
		var t string
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (ps PGPostStore) addTags(ctx context.Context, postID int, tagNames []string) error {
	var tagID int

	for _, tagName := range tagNames {
		err := ps.DB.GetContext(ctx, &tagID, "SELECT id FROM tags WHERE name=$1", tagName)
		if err != nil {
			return fmt.Errorf("could not get tag %s because of %s", tagName, err)
		}

		var relationID int
		err = ps.DB.QueryRowxContext(
			ctx,
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
			return fmt.Errorf("could not relate post %d to tag %d because of %s", postID, tagID, err)
		}
	}
	return nil
}

func (ps PGPostStore) removeTags(ctx context.Context, postID int, tagNames []string) error {
	var tagID int

	for _, tagName := range tagNames {
		err := ps.DB.GetContext(ctx, &tagID, "SELECT id FROM tags WHERE name=$1", tagName)
		if err != nil {
			return fmt.Errorf("could not get tag %s because of %s", tagName, err)
		}

		_, err = ps.DB.ExecContext(
			ctx,
			"DELETE FROM posttags WHERE tag_id=$1 AND post_id=$2",
			tagID,
			postID,
		)

		if err != nil {
			return fmt.Errorf("could not delete relation between post %d and tag %d because of %s", postID, tagID, err)
		}
	}
	return nil
}
