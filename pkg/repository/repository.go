package repository

import (
	"context"

	"github.com/joelsaunders/blog-go/pkg/models"
)

type UserStore interface {
	List(ctx context.Context, num int) ([]*models.User, error)
	Create(ctx context.Context, user *models.NewUser) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type PostStore interface {
	List(ctx context.Context) ([]*models.Post, error)
	GetBySlug(ctx context.Context, slug string) (*models.Post, error)
	Create(ctx context.Context, post *models.Post) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) (*models.Post, error)
}
