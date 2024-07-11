package grpcutil

import (
	"context"
	"math/rand"

	"github.com/richardktran/go-movie-microservices/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName)

	if err != nil {
		return nil, err
	}

	addr := addrs[rand.Intn(len(addrs))]

	return grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
