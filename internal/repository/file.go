package repository

import (
	"context"
	"database/sql"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
)

type Files struct {
	db *sql.DB
}

func NewFiles(db *sql.DB) *Files {
	return &Files{db}
}

func (r *Files) Upload(ctx context.Context, file domain.UploadFile) (domain.UploadFile, error) {
	err := r.db.QueryRowContext(ctx, "INSERT INTO files (name, author_id, comments) values ($1, $2, $3) returning id, createdAt",
		file.Name, file.AuthorId, file.Comment).Scan(&file.Id, &file.CreatedAt)
	return file, err
}
