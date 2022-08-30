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
	GetAll(ctx context.Context) ([]domain.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, post *domain.UpdatePost) error
}

type Handler struct {
	postServices Posts
}

func NewHandler(posts Posts) *Handler {
	return &Handler{
		postServices: posts,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(loggerMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	post := router.Group("/post")
	{
		post.POST("", h.Create)
		post.GET("", h.List)
		post.GET("/:id", h.GetById)
		post.PUT("", h.UpdateById)
		post.DELETE("", h.DeleteById)
	}
	return router
}
