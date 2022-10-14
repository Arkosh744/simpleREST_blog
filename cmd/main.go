package main

import (
	"github.com/Arkosh744/simpleREST_blog/internal/app"
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
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
