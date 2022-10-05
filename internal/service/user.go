package service

//go:generate mockgen -source=user.go -destination=mocks/user_mock.go -package=mocks

import (
	"context"
	"errors"
	"fmt"
	audit "github.com/Arkosh744/grpc-audit-log/pkg/domain"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type TokensRepository interface {
	Create(ctx context.Context, token domain.RefreshToken) error
	Get(ctx context.Context, token string) (domain.RefreshToken, error)
}

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type Users struct {
	Repo        UsersRepository
	TokenRepo   TokensRepository
	AuditClient AuditClient
	Hasher      PasswordHasher
	HmacSecret  []byte
}

func NewUsers(repo UsersRepository, tokenRepo TokensRepository, auditClient AuditClient, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		Repo:        repo,
		TokenRepo:   tokenRepo,
		AuditClient: auditClient,
		Hasher:      hasher,
		HmacSecret:  secret,
	}
}

func (u *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	password, err := u.Hasher.Hash(inp.Password)
	if err != nil {
		return err
	}
	_, err = u.Repo.GetByCredentials(ctx, inp.Email, password)
	if err == nil {
		return errors.New("user already exists")
	}
	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	if err := u.Repo.Create(ctx, user); err != nil {
		return err
	}

	user, err = u.Repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		return err
	}

	if err := u.AuditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_REGISTER,
		Entity:    audit.ENTITY_USER,
		UserID:    user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignUp",
		}).Error("failed to send log request:", err)
	}
	return nil
}

func (u *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password, err := u.Hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}
	user, err := u.Repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		return "", "", err
	}

	if err := u.AuditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		UserID:    user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignIn",
		}).Error("failed to send log request:", err)
	}

	return u.generateTokens(ctx, user.ID)
}

func (u *Users) generateTokens(ctx context.Context, userID int64) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := token.SignedString(u.HmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	if err := u.TokenRepo.Create(ctx, domain.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (u *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return u.HmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !tok.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (u *Users) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := u.TokenRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return u.generateTokens(ctx, token.UserID)
}

func (u *Users) GetIdByToken(ctx context.Context, refreshToken string) (int64, error) {
	token, err := u.TokenRepo.Get(ctx, refreshToken)
	if err != nil {
		return 0, err
	}
	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return 0, domain.ErrRefreshTokenExpired
	}
	return token.UserID, nil
}
