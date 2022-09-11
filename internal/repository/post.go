package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"strings"
	"time"
)

type Posts struct {
	db *sql.DB
}

func NewPosts(db *sql.DB) *Posts {
	return &Posts{db}
}

func (r *Posts) Create(ctx context.Context, post domain.Post) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO posts (title, body) values ($1, $2)",
		post.Title, post.Body)

	return err
}

func (r *Posts) GetById(ctx context.Context, id int64) (domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRowContext(ctx, "SELECT id, title, body, \"createdAt\", \"updatedAt\" FROM posts WHERE id=$1", id).
		Scan(&post.Id, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return post, domain.ErrPostNotFound
	}

	return post, err
}

func (r *Posts) List(ctx context.Context) ([]domain.Post, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, body, \"createdAt\", \"updatedAt\" FROM posts")
	if err != nil {
		return nil, err
	}

	books := make([]domain.Post, 0)
	for rows.Next() {
		var book domain.Post
		if err := rows.Scan(&book.Id, &book.Title, &book.Body, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func (r *Posts) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id=$1", id)
	return err
}

func (r *Posts) Update(ctx context.Context, id int64, post *domain.UpdatePost) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

	_, err := r.GetById(ctx, id)
	if err != nil {
		return err
	}

	if post.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *post.Title)
		argId++
	}

	if post.Body != nil {
		setValues = append(setValues, fmt.Sprintf("body=$%d", argId))
		args = append(args, *post.Body)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updatedAt=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE posts SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
