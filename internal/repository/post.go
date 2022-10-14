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

func (r *Posts) Create(ctx context.Context, post domain.Post) (domain.Post, error) {
	err := r.db.QueryRowContext(ctx, "INSERT INTO posts (title, body, author_id) values ($1, $2, $3) returning id",
		post.Title, post.Body, post.AuthorId).Scan(&post.Id)
	return post, err
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

func (r *Posts) Update(ctx context.Context, id int64, post domain.UpdatePost) (domain.Post, error) {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

	newPost, err := r.GetById(ctx, id)
	if err != nil {
		return domain.Post{}, err
	}

	if post.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, post.Title)
		newPost.Title = post.Title
		argId++
	}

	if post.Body != "" {
		setValues = append(setValues, fmt.Sprintf("body=$%d", argId))
		args = append(args, post.Body)
		newPost.Body = post.Body
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("\"updatedAt\"=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("%s %s %s%d", "UPDATE posts SET", setQuery, "WHERE id=$", argId)
	args = append(args, id)
	_, err = r.db.ExecContext(ctx, query, args...)
	return newPost, err
}
