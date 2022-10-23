package service

//go:generate mockgen -source=file.go -destination=mocks/file_mock.go -package=mocks

import (
	"context"
	customCache "github.com/Arkosh744/FirstCache"
	audit "github.com/Arkosh744/grpc-audit-log/pkg/domain"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/sirupsen/logrus"
)

type FilesRepository interface {
	Upload(ctx context.Context, file domain.UploadFile) (domain.UploadFile, error)
}

type Files struct {
	repo        FilesRepository
	cache       *customCache.Cache
	auditClient AuditClient
}

func NewFiles(repo FilesRepository, cache *customCache.Cache, auditClient AuditClient) *Files {
	return &Files{
		repo:        repo,
		cache:       cache,
		auditClient: auditClient,
	}
}

func (p *Files) Upload(ctx context.Context, file domain.UploadFile) error {
	newfile, err := p.repo.Upload(ctx, file)
	if err != nil {
		return err
	}

	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:   audit.ACTION_CREATE,
		Entity:   audit.ENTITY_POST,
		EntityID: newfile.Id,
		UserID:   newfile.AuthorId,
		//Timestamp: newfile.CreatedAt,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Post.Create",
		}).Error("failed to send log request:", err)
	}

	return nil
}
