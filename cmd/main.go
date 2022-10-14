package main

import (
	"context"
	"github.com/Arkosh744/simpleREST_blog/internal/config"
	"github.com/Arkosh744/simpleREST_blog/internal/repository"
	"github.com/Arkosh744/simpleREST_blog/internal/service"
	grpc_client "github.com/Arkosh744/simpleREST_blog/internal/transport/grpc"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest"
	"github.com/Arkosh744/simpleREST_blog/pkg/database"
	"github.com/Arkosh744/simpleREST_blog/pkg/hash"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cache "github.com/Arkosh744/FirstCache"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// @title Note-taking API
// @version 1.0
// @description This is a simple crud blog for note-taking

// @host localhost:8080
// @BasePath /

// @termsOfService http://swagger.io/terms/
// @host localhost:8080
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		quit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
		cancel()
	}()

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

	auditClient, err := grpc_client.NewClient(9000)
	if err != nil {
		log.Fatal(err)
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

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	})

	log.Info("SERVER STARTED")

	if err := g.Wait(); err != nil {
		log.WithFields(log.Fields{
			"context": "app.Run()",
			"problem": "server shutdowned",
		}).Error(err.Error())
	}
}
