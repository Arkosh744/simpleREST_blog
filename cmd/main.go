package main

import (
	"github.com/Arkosh744/simpleREST_blog/internal/repository"
	"github.com/Arkosh744/simpleREST_blog/internal/service"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest"
	"github.com/Arkosh744/simpleREST_blog/pkg/database"
	"log"
	"net/http"
	"time"

	cache "github.com/Arkosh744/FirstCache"
	_ "github.com/lib/pq"
)

// @title Note-taking API
// @version 1.0
// @description This is a simple crud blog for note-taking

// @host localhost:8080
// @BasePath /

// @termsOfService http://swagger.io/terms/
// @host localhost:8080
func main() {

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "docker",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	postsRepo := repository.NewPosts(db)
	handlerCache := cache.NewCache()
	postService := service.NewPosts(postsRepo, handlerCache)
	handler := rest.NewHandler(postService)

	// init & run server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler.InitRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf(err.Error())
	}
}
