package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/richardktran/go-movie-microservices/rating/internal/controller/rating"
	httpHandler "github.com/richardktran/go-movie-microservices/rating/internal/handler/http"
	"github.com/richardktran/go-movie-microservices/rating/internal/repository/memory"
)

func main() {
	port := 8082
	log.Println("Starting rating service with port %v...", port)

	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		panic(err)
	}
}
