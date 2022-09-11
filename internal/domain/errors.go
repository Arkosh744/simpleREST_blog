package domain

import "errors"

var (
	ErrPostNotFound        = errors.New("post not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
