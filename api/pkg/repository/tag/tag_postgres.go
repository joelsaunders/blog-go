package tag

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/joelsaunders/blog-go/api/pkg/models"
)

type PGTagStore struct {
	DB *sqlx.DB
}

func (ts *PGTagStore) List(ctx context.Context) ([]*models.Tag, error) {
	query := `SELECT id, name FROM tags t ORDER BY t.id`
	var tags []*models.Tag

	err := ts.DB.SelectContext(ctx, &tags, query)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (ts *PGTagStore) Create(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	query := `INSERT INTO tags (name) VALUES ($1) RETURNING id, name`
	insertedTag := &models.Tag{}

	err := ts.DB.GetContext(ctx, insertedTag, query, tag.Name)

	if err != nil {
		return nil, err
	}
	return insertedTag, nil
}

func (ts *PGTagStore) Update(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	query := `UPDATE tags SET name=$2 WHERE id=$1 RETURNING id, name`
	updatedTag := &models.Tag{}

	err := ts.DB.GetContext(ctx, updatedTag, query, tag.ID, tag.Name)
	if err != nil {
		return nil, err
	}
	return updatedTag, nil
}

func handleTransactionError(err error, tx *sqlx.Tx) {
	if err != nil {
		rErr := tx.Rollback()
		log.Fatalf("error in transaction: %s \n rollback state: %s", err, rErr)
	}
}

func (ts *PGTagStore) DeleteByID(ctx context.Context, ID int) error {
	postTagsDeleteQuery := `DELETE FROM posttags pt WHERE pt.tag_id = $1`
	tagDeleteQuery := `DELETE FROM tags t WHERE t.id = $1`

	tx := ts.DB.MustBegin()
	_, err := tx.ExecContext(ctx, postTagsDeleteQuery, ID)
	handleTransactionError(err, tx)
	_, err = tx.ExecContext(ctx, tagDeleteQuery, ID)
	handleTransactionError(err, tx)
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("db error committing transaction to delete tag %d: %s", ID, err)
	}
	return nil
}
