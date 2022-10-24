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
	GetById(ctx context.Context, id int64, userId int64) (domain.UploadFile, error)
	List(ctx context.Context, userId int64) ([]domain.UploadFile, error)
	Delete(ctx context.Context, id int64, userId int64) error
}

type Files struct {
	repo        FilesRepository
	cache       *customCache.Cache
	auditClient AuditClient
}

func NewFiles(repo FilesRepository, auditClient AuditClient) *Files {
	return &Files{
		repo:        repo,
		auditClient: auditClient,
	}
}

func (f *Files) Upload(ctx context.Context, file domain.UploadFile) error {
	newfile, err := f.repo.Upload(ctx, file)
	if err != nil {
		return err
	}

	if err := f.auditClient.SendLogRequest(ctx, audit.LogItem{
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

func (f *Files) GetById(ctx context.Context, id int64, userId int64) (domain.UploadFile, error) {
	return domain.UploadFile{}, nil
}

func (f *Files) List(ctx context.Context, userId int64) ([]domain.UploadFile, error) {
	return []domain.UploadFile{}, nil
}

func (f *Files) Delete(ctx context.Context, id int64, userId int64) error {
	return nil
}
