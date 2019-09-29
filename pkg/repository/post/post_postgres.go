package post

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/bilbo-go/pkg/models"
)

type PGPostStore struct {
	DB *sqlx.DB
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING slug`

	var lastInsertSlug string

	err := ps.DB.QueryRowx(
		query,
		post.Slug,
		post.Title,
		post.Body,
		post.Picture,
		post.Description,
		post.Published,
		post.AuthorID,
	).Scan(&lastInsertSlug)

	if err != nil {
		return nil, err
	}

	createdPost, err := ps.GetBySlug(ctx, lastInsertSlug)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (ps *PGPostStore) Update(ctx context.Context, post *models.Post) (*models.Post, error) {
	fmt.Println(post)
	row := ps.DB.QueryRowxContext(
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
		WHERE id=$1;`,
		post.ID,
		post.Slug,
		post.Title,
		post.Body,
		post.Picture,
		post.Description,
		post.Published,
		post.AuthorID,
		time.Now().UTC(),
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	insertedPost, err := ps.GetBySlug(ctx, post.Slug)
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
	return &post, nil
}
