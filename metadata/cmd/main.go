package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/richardktran/go-movie-microservices/metadata/internal/controller/metadata"
	httphandler "github.com/richardktran/go-movie-microservices/metadata/internal/handler/http"
	"github.com/richardktran/go-movie-microservices/metadata/internal/repository/memory"
)

func main() {
	port := 8081
	log.Printf("Starting the movie metadata service with port %v...", port)

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		panic(err)
	}
}
