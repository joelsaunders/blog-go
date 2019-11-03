package post

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/joelsaunders/blog-go/api/pkg/models"
	"github.com/joelsaunders/blog-go/api/pkg/repository/postgres"
)

var allowedPostFilters = map[string]string{
	"tag_name": "t.name",
}

type PGPostStore struct {
	DB *sqlx.DB
}

type lastInsertDetails struct {
	ID   int    `db:"id"`
	Slug string `db:"slug"`
}

func (ps *PGPostStore) DeleteBySlug(ctx context.Context, postSlug string) error {
	postTagsDeleteQuery := "DELETE FROM posttags pt WHERE pt.post_id = $1"
	postDeleteQuery := "DELETE FROM posts WHERE posts.slug = $1"

	post, err := ps.GetBySlug(ctx, postSlug)
	if err != nil {
		return fmt.Errorf("could not fetch post to delete: %s", err)
	}

	tx := ps.DB.MustBegin()
	_, err = tx.ExecContext(ctx, postTagsDeleteQuery, post.ID)
	if err != nil {
		rErr := tx.Rollback()
		log.Fatalf("error in delete db transaction: %s \n rollback state: %s", err, rErr)
	}
	_, err = tx.ExecContext(ctx, postDeleteQuery, post.Slug)
	if err != nil {
		rErr := tx.Rollback()
		log.Fatalf("error in delete db transaction: %s \n rollback state: %s", err, rErr)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("database error deleting post: %s", err)
	}
	return nil
}

func (ps *PGPostStore) Create(ctx context.Context, post *models.Post) (*models.Post, error) {
	query := `INSERT INTO posts (
		slug,
		title,
		body,
		picture,
		description,
		published,
		author_id,
		modified,
		created
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, slug`

	if post.Created.IsZero() {
		post.Created = time.Now().UTC()
	}

	if post.Modified.IsZero() {
		post.Modified = time.Now().UTC()
	}

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
		post.Modified,
		post.Created,
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
	err = ps.addTags(ctx, updateRow.ID, toAdd)
	toRemove := difference(existingTags, post.Tags)
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

func (ps *PGPostStore) List(ctx context.Context, filters map[string][]string) ([]*models.Post, error) {
	query := `
	SELECT 
		p.*,
		u.email as author_email,
		array_agg(t.name) FILTER (WHERE t.name IS NOT NULL) as tags
	FROM posts p
		INNER JOIN users u ON u.id = p.author_id
		LEFT JOIN posttags pt ON pt.post_id = p.id
		LEFT JOIN tags t ON t.id = pt.tag_id
	%s -- this is where the filters will go, they will still use bindvars so no injection possible
	GROUP BY (
		p.slug, p.id, p.created, p.modified, p.title, p.body, p.picture,
		p.description, p.published, p.author_id,
		author_email
	)
	ORDER BY p.created desc
	`
	query, varList, err := postgres.BuildFilterString(query, filters, allowedPostFilters)
	if err != nil {
		return nil, err
	}

	rows, err := ps.DB.QueryxContext(ctx, query, varList...)
	if err != nil {
		return nil, err
	}
	defer postgres.Close(rows)

	posts := make([]*models.Post, 0)

	for rows.Next() {
		var p models.Post
		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

func (ps *PGPostStore) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	query := `SELECT 
		p.*,
		u.email as author_email,
		array_agg(t.name) as tags
	FROM posts p INNER JOIN users u ON u.id = p.author_id
	INNER JOIN posttags pt ON pt.post_id = p.id
	INNER JOIN tags t ON t.id = pt.tag_id
	WHERE slug=$1
	GROUP BY (
	    p.slug, p.id, p.created, p.modified, p.title, p.body, p.picture,
		p.description, p.published, p.author_id,
		author_email
	);`
	post := models.Post{}
	err := ps.DB.GetContext(ctx, &post, query, slug)
	if err != nil {
		return nil, err
	}
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
	defer postgres.Close(rows)

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
