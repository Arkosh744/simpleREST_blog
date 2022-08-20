package service

import (
	"context"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"time"
)

type PostsRepository interface {
	Create(ctx context.Context, post domain.Post) error
	GetById(ctx context.Context, id int64) (domain.Post, error)
	GetAll(ctx context.Context) ([]domain.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, post *domain.UpdatePost) error
}

type Posts struct {
	repo PostsRepository
}

func NewPosts(repo PostsRepository) *Posts {
	return &Posts{
		repo: repo,
	}
}

func (p *Posts) Create(ctx context.Context, post domain.Post) error {
	post.Date = time.Now()
	post.Updated = time.Now()
	return p.repo.Create(ctx, post)
}

func (p *Posts) GetById(ctx context.Context, id int64) (domain.Post, error) {
	return p.repo.GetById(ctx, id)
}

func (p *Posts) GetAll(ctx context.Context) ([]domain.Post, error) {
	return p.repo.GetAll(ctx)
}

func (p *Posts) Delete(ctx context.Context, id int64) error {
	return p.repo.Delete(ctx, id)
}

func (p *Posts) Update(ctx context.Context, id int64, post *domain.UpdatePost) error {
	return p.repo.Update(ctx, id, post)
}
