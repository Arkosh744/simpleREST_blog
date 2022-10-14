package app

import (
	"context"
	cache "github.com/Arkosh744/FirstCache"
	"github.com/Arkosh744/simpleREST_blog/internal/config"
	"github.com/Arkosh744/simpleREST_blog/internal/repository"
	"github.com/Arkosh744/simpleREST_blog/internal/service"
	grpc_client "github.com/Arkosh744/simpleREST_blog/internal/transport/grpc"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest"
	"github.com/Arkosh744/simpleREST_blog/pkg/database"
	"github.com/Arkosh744/simpleREST_blog/pkg/hash"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Note-taking API
// @version 1.0
// @description This is a simple crud blog for note-taking

// @host localhost:8080
// @BasePath /

// @termsOfService http://swagger.io/terms/
// @host localhost:8080
func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	cfg, err := config.New("configs")
	if err != nil {
		return err
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
		return err
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher("Salty Salt")

	postsRepo := repository.NewPosts(db)
	handlerCache := cache.NewCache()
	usersRepo := repository.NewUsers(db)
	tokensRepo := repository.NewTokens(db)

	auditClient, err := grpc_client.NewClient(9000)
	if err != nil {
		return err
	}
	postService := service.NewPosts(postsRepo, handlerCache, auditClient)
	usersService := service.NewUsers(usersRepo, tokensRepo, auditClient, hasher, []byte(cfg.JWTSecret))

	handler := rest.NewHandler(postService, usersService)

	// init & run server
	srv := &http.Server{
		Addr:         ":" + cfg.SrvPort,
		Handler:      handler.InitRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		// service connections
		log.Info("Starting Server...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()
	log.Info("SERVER STARTED")

	// GRACEFUL SHUTDOWN with 5 seconds timeout BELOW
	quit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
	return nil
}
