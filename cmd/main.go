package main

import (
	"github.com/Arkosh744/simpleREST_blog/internal/app"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
