package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/richardktran/go-movie-microservices/pkg/discovery"
)

type Registry struct {
	sync.RWMutex
	// serviceName -> instanceID -> serviceInstance
	serviceAddrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: map[string]map[string]*serviceInstance{},
	}
}

// Register creates a service instance in the registry
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[string]*serviceInstance{}
	}

	r.serviceAddrs[serviceName][instanceID] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceName], instanceID)

	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	instances, ok := r.serviceAddrs[serviceName]
	if !ok || len(instances) == 0 {
		return nil, discovery.ErrNotFound
	}

	var addrs []string
	for _, instance := range instances {
		addrs = append(addrs, instance.hostPort)
	}

	return addrs, nil
}

func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}
