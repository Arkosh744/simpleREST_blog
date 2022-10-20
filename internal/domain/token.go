package domain

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
}

type Token struct {
	Token string `json:"token"`
}
