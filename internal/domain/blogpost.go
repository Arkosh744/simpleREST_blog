package domain

import (
	"errors"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Post struct {
	Id      int64     `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Date    time.Time `json:"date"`
	Updated time.Time `json:"lastUpdated"`
}

type UpdatePost struct {
	Id    int64   `json:"id"`
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type PostError struct {
	MsgErr string `json:"message"`
}
