package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/bilbo-go/models"
)

type PGUserStore struct {
	DB *sqlx.DB
}

func (us *PGUserStore) List(ctx context.Context, num int) ([]*models.User, error) {
	query := "Select id, email, password from users limit $1"
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
