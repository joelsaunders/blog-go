package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/bilbo-go/models"
	"github.com/joelsaunders/bilbo-go/repository/user"
)

type UserStore interface {
	List(ctx context.Context, num int) ([]*models.User, error)
	Create(ctx context.Context, user *models.NewUser) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type DB struct {
	DB *sqlx.DB
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{DB: db}
}

func (db *DB) Users() UserStore {
	return &user.PGUserStore{DB: db.DB}
}
