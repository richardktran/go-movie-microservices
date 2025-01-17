package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/richardktran/go-movie-microservices/gen"
	"github.com/richardktran/go-movie-microservices/metadata/internal/controller/metadata"
	grpcHandler "github.com/richardktran/go-movie-microservices/metadata/internal/handler/grpc"
	"github.com/richardktran/go-movie-microservices/metadata/internal/repository/mysql"
	"github.com/richardktran/go-movie-microservices/pkg/discovery"
	"github.com/richardktran/go-movie-microservices/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

var serviceName = "metadata"

type serviceConfig struct {
	APIConfig apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port string `yaml:"port"`
}

func main() {
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	port := cfg.APIConfig.Port
	log.Printf("Starting the movie metadata service with port %v...", port)

	// Register the metadata service
	registry, err := consul.NewRegistry(":8500")
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
				time.Sleep(1 * time.Second)
			}
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// repo := memory.New()

	repo, err := mysql.New()

	if err != nil {
		panic(err)
	}

	ctrl := metadata.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)
	gen.RegisterMetadataServiceServer(server, h)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}

	// =============== This section is for HTTP handler ===============
	// h := httphandler.New(ctrl)

	// http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	// if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
	// 	panic(err)
	// }
}
