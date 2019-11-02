package repository

import (
	"context"

	"github.com/joelsaunders/blog-go/api/pkg/models"
)

type UserStore interface {
	List(ctx context.Context, num int) ([]*models.User, error)
	Create(ctx context.Context, user *models.NewUser) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}

type PostStore interface {
	List(ctx context.Context, filters map[string]string) ([]*models.Post, error)
	GetBySlug(ctx context.Context, slug string) (*models.Post, error)
	Create(ctx context.Context, post *models.Post) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) (*models.Post, error)
	DeleteBySlug(ctx context.Context, postSlug string) error
}

type TagStore interface {
	List(ctx context.Context) ([]*models.Tag, error)
	Create(ctx context.Context, tag *models.Tag) (*models.Tag, error)
	Update(ctx context.Context, tag *models.Tag) (*models.Tag, error)
	DeleteByID(ctx context.Context, ID int) error
}
