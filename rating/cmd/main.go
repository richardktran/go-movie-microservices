package main

import (
	"fmt"
	"net/http"

	"github.com/richardktran/go-movie-microservices/rating/internal/controller/rating"
	httpHandler "github.com/richardktran/go-movie-microservices/rating/internal/handler/http"
	"github.com/richardktran/go-movie-microservices/rating/internal/repository/memory"
)

func main() {
	fmt.Println("Starting rating service...")

	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
