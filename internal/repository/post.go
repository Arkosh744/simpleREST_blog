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

func (b *Posts) Create(ctx context.Context, post domain.Post) error {
	_, err := b.db.ExecContext(ctx, "INSERT INTO posts (title, body) values ($1, $2)",
		post.Title, post.Body)

	return err
}

func (b *Posts) GetById(ctx context.Context, id int64) (domain.Post, error) {
	var post domain.Post
	err := b.db.QueryRowContext(ctx, "SELECT id, title, body, date, updated FROM posts WHERE id=$1", id).
		Scan(&post.Id, &post.Title, &post.Body, &post.Date, &post.Updated)
	if err == sql.ErrNoRows {
		return post, domain.ErrPostNotFound
	}

	return post, err
}

func (b *Posts) GetAll(ctx context.Context) ([]domain.Post, error) {
	rows, err := b.db.QueryContext(ctx, "SELECT id, title, body, date, updated FROM posts")
	if err != nil {
		return nil, err
	}

	books := make([]domain.Post, 0)
	for rows.Next() {
		var book domain.Post
		if err := rows.Scan(&book.Id, &book.Title, &book.Body, &book.Date, &book.Updated); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func (b *Posts) Delete(ctx context.Context, id int64) error {
	_, err := b.db.ExecContext(ctx, "DELETE FROM posts WHERE id=$1", id)
	return err
}

func (b *Posts) Update(ctx context.Context, id int64, post *domain.UpdatePost) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

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

	setValues = append(setValues, fmt.Sprintf("updated=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE posts SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)
	_, err := b.db.ExecContext(ctx, query, args...)
	return err
}
