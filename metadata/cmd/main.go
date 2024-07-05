package main

import (
	"log"
	"net/http"

	"github.com/richardktran/go-movie-microservices/metadata/internal/controller/metadata"
	httphandler "github.com/richardktran/go-movie-microservices/metadata/internal/handler/http"
	"github.com/richardktran/go-movie-microservices/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting the movie metadata service")

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
