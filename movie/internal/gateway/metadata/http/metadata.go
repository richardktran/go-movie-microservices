package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/richardktran/go-movie-microservices/metadata/pkg/model"
	"github.com/richardktran/go-movie-microservices/movie/internal/gateway"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	addr string
}

// New creates a new HTTP Gateway for a movie metadata service.
func New(addr string) *Gateway {
	return &Gateway{addr}
}

// Get gets movie metadata by a movie id
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	// Build the raw request
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)

	if err != nil {
		return nil, err
	}

	// Add query string to the request
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	// Call to the metadata service by http request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	// Error handler
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	// Receive and decode the response body
	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}
