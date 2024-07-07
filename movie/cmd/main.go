package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/richardktran/go-movie-microservices/movie/internal/controller/movie"
	metadataGateway "github.com/richardktran/go-movie-microservices/movie/internal/gateway/metadata/http"
	ratingGateway "github.com/richardktran/go-movie-microservices/movie/internal/gateway/rating/http"
	movieHandler "github.com/richardktran/go-movie-microservices/movie/internal/handler/http"
	"github.com/richardktran/go-movie-microservices/pkg/discovery"
	"github.com/richardktran/go-movie-microservices/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "The server port")
	flag.Parse()
	log.Printf("Starting the %s service with port %v...", serviceName, port)

	// Register the movie service
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%v", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state:", err.Error())
				// Sleep for 1 second
				time.Sleep(1 * time.Second)
			}
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadataGateway.New(registry)
	ratingGateway := ratingGateway.New(registry)

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := movieHandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		panic(err)
	}

}
