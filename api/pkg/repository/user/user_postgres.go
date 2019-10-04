package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/blog-go/api/pkg/models"
)

type PGUserStore struct {
	DB *sqlx.DB
}

func (us *PGUserStore) List(ctx context.Context, num int) ([]*models.User, error) {
	query := "Select id, email from users limit $1"
	rows, err := us.DB.QueryxContext(ctx, query, num)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)

	for rows.Next() {
		var u models.User
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (us *PGUserStore) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := models.User{}
	err := us.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *PGUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}
	err := us.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *PGUserStore) Create(ctx context.Context, user *models.NewUser) (*models.User, error) {
	// TODO: make the create return all the relevant rows
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var lastInsertID int
	err := us.DB.QueryRowx(query, user.Email, user.Password).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	createdUser, err := us.GetByID(ctx, lastInsertID)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
