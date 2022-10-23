package rest

import (
	"context"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/Arkosh744/simpleREST_blog/docs"
)

//go:generate mockgen -source=handler.go -destination=mocks/mocks.go -package=mocks

type Posts interface {
	Create(ctx context.Context, post domain.Post) error
	GetById(ctx context.Context, id int64, userId int64) (domain.Post, error)
	List(ctx context.Context, userId int64) ([]domain.Post, error)
	Delete(ctx context.Context, id int64, userId int64) error
	Update(ctx context.Context, id int64, post domain.UpdatePost, userId int64) error
}

type Users interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	GetIdByToken(ctx context.Context, refreshToken string) (int64, error)
}

//type Files interface {
//	Upload(ctx context.Context, file domain.UploadFile) error
//	GetById(ctx context.Context, id int64, userId int64) (domain.UploadFile, error)
//	List(ctx context.Context, userId int64) ([]domain.UploadFile, error)
//	Delete(ctx context.Context, id int64, userId int64) error
//}

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

// @title Post Service REST API
// @version 1.0
// @description This is a simple crud blog for posts
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
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
	file := router.Group("/file")
	{
		file.Use(h.authMiddleware())
		file.POST("", h.UploadFile)
		//file.GET("/", h.ListFile)
		//file.GET("/:id", h.GetFile)
		//file.GET("/:id", h.DeleteFile)
	}
	router.Use(loggerMiddleware())
	return router
}
