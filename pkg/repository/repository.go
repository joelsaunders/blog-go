package repository

import (
	"context"

	"github.com/joelsaunders/bilbo-go/pkg/models"
)

type UserStore interface {
	List(ctx context.Context, num int) ([]*models.User, error)
	Create(ctx context.Context, user *models.NewUser) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
