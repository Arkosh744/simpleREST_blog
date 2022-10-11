package domain

import "errors"

var (
	ErrPostNotFound        = errors.New("post not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidInput        = errors.New("invalid input body")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
