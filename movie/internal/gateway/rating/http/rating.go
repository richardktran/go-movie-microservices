package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/richardktran/go-movie-microservices/movie/internal/gateway"
	"github.com/richardktran/go-movie-microservices/rating/pkg/model"
)

// Gateway defines an HTTP gateway for a rating service.
type Gateway struct {
	addr string
}

// New creates a new HTTP gateway for a rating service.
func New(addr string) *Gateway {
	return &Gateway{addr}
}

// GetAggregatedRating returns the aggregated rating for a record
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	// Build the raw request
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)

	if err != nil {
		return 0, err
	}

	// Add query string to the request
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	// Receive and decode data
	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return err
	}

	// Add query string
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	// Call HTTP request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
