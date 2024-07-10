package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/richardktran/go-movie-microservices/gen"
	"github.com/richardktran/go-movie-microservices/movie/internal/controller/movie"
	metadataGateway "github.com/richardktran/go-movie-microservices/movie/internal/gateway/metadata/grpc"
	ratingGateway "github.com/richardktran/go-movie-microservices/movie/internal/gateway/rating/grpc"
	movieGrpcHandler "github.com/richardktran/go-movie-microservices/movie/internal/handler/grpc"
	"github.com/richardktran/go-movie-microservices/pkg/discovery"
	"github.com/richardktran/go-movie-microservices/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	h := movieGrpcHandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)
	gen.RegisterMovieServiceServer(server, h)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// =============== This section is for HTTP handler ===============
	// h := movieHttpHandler.New(ctrl)

	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	// if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
	// 	panic(err)
	// }

}
