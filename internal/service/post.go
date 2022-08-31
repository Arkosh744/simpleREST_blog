package service

import (
	"context"
	customCache "github.com/Arkosh744/FirstCache"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"strconv"
	"time"
)

type PostsRepository interface {
	Create(ctx context.Context, post domain.Post) error
	GetById(ctx context.Context, id int64) (domain.Post, error)
	List(ctx context.Context) ([]domain.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, post *domain.UpdatePost) error
}

type Posts struct {
	repo  PostsRepository
	cache *customCache.Cache
}

func NewPosts(repo PostsRepository, cache *customCache.Cache) *Posts {
	return &Posts{
		repo:  repo,
		cache: cache,
	}
}

func (p *Posts) Create(ctx context.Context, post domain.Post) error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt
	p.cache.Set(strconv.FormatInt(post.Id, 10), post, time.Second*360, ctx)
	return p.repo.Create(ctx, post)
}

func (p *Posts) GetById(ctx context.Context, id int64) (domain.Post, error) {
	if post, err := p.cache.Get(strconv.FormatInt(id, 10)); err == nil {
		return post.Value.(domain.Post), err
	} else {
		post, err := p.repo.GetById(ctx, id)
		return post, err
	}
}

func (p *Posts) List(ctx context.Context) ([]domain.Post, error) {
	posts, err := p.repo.List(ctx)
	for _, item := range posts {
		if _, err := p.cache.Get(strconv.FormatInt(item.Id, 10)); err != nil {
			p.cache.Set(strconv.FormatInt(item.Id, 10), item, time.Second*360, ctx)
		}
	}
	return posts, err
}

func (p *Posts) Delete(ctx context.Context, id int64) error {
	if _, err := p.cache.Get(strconv.FormatInt(id, 10)); err == nil {
		_ = p.cache.Delete(strconv.FormatInt(id, 10))
	}
	return p.repo.Delete(ctx, id)
}

func (p *Posts) Update(ctx context.Context, id int64, post *domain.UpdatePost) error {
	p.cache.Set(strconv.FormatInt(post.Id, 10), post, time.Second*360, ctx)
	return p.repo.Update(ctx, id, post)
}
