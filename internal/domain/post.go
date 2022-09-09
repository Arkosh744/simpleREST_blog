package domain

import (
	"time"
)

type Post struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	AuthorId  int64     `json:"AuthorId"`
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
