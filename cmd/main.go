package main

import (
	"github.com/Arkosh744/simpleREST_blog/internal/config"
	"github.com/Arkosh744/simpleREST_blog/internal/repository"
	"github.com/Arkosh744/simpleREST_blog/internal/service"
	grpc_client "github.com/Arkosh744/simpleREST_blog/internal/transport/grpc"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest"
	"github.com/Arkosh744/simpleREST_blog/pkg/database"
	"github.com/Arkosh744/simpleREST_blog/pkg/hash"
	"net/http"
	"os"
	"time"

	cache "github.com/Arkosh744/FirstCache"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// @title Note-taking API
// @version 1.0
// @description This is a simple crud blog for note-taking

// @host localhost:8080
// @BasePath /

// @termsOfService http://swagger.io/terms/
// @host localhost:8080
func main() {
	cfg, err := config.New("configs")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUser,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Password: cfg.DBPassword,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher("Salty Salt")

	postsRepo := repository.NewPosts(db)
	handlerCache := cache.NewCache()
	usersRepo := repository.NewUsers(db)
	tokensRepo := repository.NewTokens(db)

	postService := service.NewPosts(postsRepo, handlerCache)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte(cfg.JWTSecret))
	auditClient, err := grpc_client.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}
	postService := service.NewPosts(postsRepo, auditClient)
	usersService := service.NewUsers(usersRepo, tokensRepo, auditClient, hasher, []byte(cfg.JWTSecret))

	handler := rest.NewHandler(postService, usersService)

	// init & run server
	srv := &http.Server{
		Addr:         ":" + cfg.SrvPort,
		Handler:      handler.InitRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf(err.Error())
	}
}
