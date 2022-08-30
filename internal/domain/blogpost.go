package domain

import (
	"errors"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Post struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostQuery struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdatePost struct {
	Id    int64   `json:"id"`
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type PostError struct {
	MsgErr string `json:"message"`
}
