package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/richardktran/go-movie-microservices/movie/internal/controller/movie"
	metadataGateway "github.com/richardktran/go-movie-microservices/movie/internal/gateway/metadata/http"
	ratingGateway "github.com/richardktran/go-movie-microservices/movie/internal/gateway/rating/http"
	movieHandler "github.com/richardktran/go-movie-microservices/movie/internal/handler/http"
)

func main() {
	port := 8083
	log.Printf("Starting the movie service with port %v...", port)
	metadataGateway := metadataGateway.New("localhost:8081")
	ratingGateway := ratingGateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := movieHandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		panic(err)
	}

}
