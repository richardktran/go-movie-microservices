package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry defines the interface for service discovery
type Registry interface {
	// Register creates a service instance in the registry
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error

	// Deregister removes a service instance from the registry
	Deregister(ctx context.Context, instanceID string, serviceName string) error

	// GetServiceAddresses returns the addresses of active instances of a given service
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)

	// ReportHealthyState reports the service instance is healthy
	ReportHealthyState(instanceID string, serviceName string) error
}

var ErrNotFound = errors.New("no service addresses found")

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
