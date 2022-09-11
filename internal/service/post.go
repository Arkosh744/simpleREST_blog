package service

import (
	"context"
	customCache "github.com/Arkosh744/FirstCache"
	audit "github.com/Arkosh744/grpc-audit-log/pkg/domain"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type PostsRepository interface {
	Create(ctx context.Context, post domain.Post) (domain.Post, error)
	GetById(ctx context.Context, id int64) (domain.Post, error)
	List(ctx context.Context) ([]domain.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, post *domain.UpdatePost) error
}

type Posts struct {
	repo        PostsRepository
	cache       *customCache.Cache
	auditClient AuditClient
}

func NewPosts(repo PostsRepository, cache *customCache.Cache, auditClient AuditClient) *Posts {
	return &Posts{
		repo:        repo,
		cache:       cache,
		auditClient: auditClient,
	}
}

func (p *Posts) Create(ctx context.Context, post domain.Post) error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt
	p.cache.Set(strconv.FormatInt(post.Id, 10), post, time.Second*360, ctx)
	newPost, err := p.repo.Create(ctx, post)
	if err != nil {
		return err
	}

	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_POST,
		EntityID:  newPost.Id,
		UserID:    newPost.AuthorId,
		Timestamp: newPost.CreatedAt,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Post.Create",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (p *Posts) GetById(ctx context.Context, id int64, userId int64) (domain.Post, error) {
	if post, err := p.cache.Get(strconv.FormatInt(id, 10)); err == nil {
		return post.Value.(domain.Post), err
	} else {
		post, err := p.repo.GetById(ctx, id)
		if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
			Action:    audit.ACTION_GET,
			Entity:    audit.ENTITY_POST,
			EntityID:  post.Id,
			UserID:    userId,
			Timestamp: post.CreatedAt,
		}); err != nil {
			logrus.WithFields(logrus.Fields{
				"method": "Post.GetById",
			}).Error("failed to send log request:", err)
		}
		return post, err
	}
}

func (p *Posts) List(ctx context.Context, userId int64) ([]domain.Post, error) {
	posts, err := p.repo.List(ctx)
	for _, item := range posts {
		if _, err := p.cache.Get(strconv.FormatInt(item.Id, 10)); err != nil {
			p.cache.Set(strconv.FormatInt(item.Id, 10), item, time.Second*360, ctx)
		}
	}

	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LIST,
		Entity:    audit.ENTITY_POST,
		UserID:    userId,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Post.List",
		}).Error("failed to send log request:", err)
	}
	return posts, err
}

func (p *Posts) Delete(ctx context.Context, id int64, userId int64) error {
	if _, err := p.cache.Get(strconv.FormatInt(id, 10)); err == nil {
		_ = p.cache.Delete(strconv.FormatInt(id, 10))
	}
	err := p.repo.Delete(ctx, id)
	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_DELETE,
		Entity:    audit.ENTITY_POST,
		EntityID:  id,
		UserID:    userId,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Post.Delete",
		}).Error("failed to send log request:", err)
	}
	return err
}

func (p *Posts) Update(ctx context.Context, id int64, post *domain.UpdatePost, userId int64) error {
	p.cache.Set(strconv.FormatInt(post.Id, 10), post, time.Second*360, ctx)
	err := p.repo.Update(ctx, id, post)

	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_UPDATE,
		Entity:    audit.ENTITY_POST,
		EntityID:  id,
		UserID:    userId,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Post.Update",
		}).Error("failed to send log request:", err)
	}
	return err
}
