package rest

import (
	"context"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/Arkosh744/simpleREST_blog/docs"
)

type Posts interface {
	Create(ctx context.Context, post domain.Post) error
	GetById(ctx context.Context, id int64) (domain.Post, error)
	List(ctx context.Context) ([]domain.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, post *domain.UpdatePost) error
}

type Users interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type Handler struct {
	postsService Posts
	usersService Users
}

func NewHandler(posts Posts, users Users) *Handler {
	return &Handler{
		postsService: posts,
		usersService: users,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(loggerMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
	}
	post := router.Group("/post")
	{
		post.Use(h.authMiddleware())
		post.POST("", h.Create)
		post.GET("", h.List)
		post.GET("/:id", h.GetById)
		post.PUT("", h.UpdateById)
		post.DELETE("", h.DeleteById)
	}
	return router
}
