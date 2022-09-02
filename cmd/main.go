package main

import (
	"github.com/Arkosh744/simpleREST_blog/internal/config"
	"github.com/Arkosh744/simpleREST_blog/internal/repository"
	"github.com/Arkosh744/simpleREST_blog/internal/service"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest"
	"github.com/Arkosh744/simpleREST_blog/pkg/database"
	"github.com/Arkosh744/simpleREST_blog/pkg/hash"
	"net/http"
	"os"
	"time"

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
	usersRepo := repository.NewUsers(db)

	postService := service.NewPosts(postsRepo)
	usersService := service.NewUsers(usersRepo, hasher, []byte(cfg.JWTSecret))

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
