package main

import (
	"context"
	"fmt"

	"github.com/richardktran/go-movie-microservices/rating/internal/controller/rating"
	"github.com/richardktran/go-movie-microservices/rating/internal/ingester/kafka"
	"github.com/richardktran/go-movie-microservices/rating/internal/repository/memory"
)

func main() {
	fmt.Println("Starting ingestion...")
	repo := memory.New()

	ingester, err := kafka.NewIngester("localhost:9092", "rating", "ratings")
	if err != nil {
		panic(err)
	}

	ctrl := rating.New(repo, ingester)
	ctx := context.Background()

	if err := ctrl.StartIngestion(ctx); err != nil {
		panic(err)
	}
}
