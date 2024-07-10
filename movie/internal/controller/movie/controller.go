package movie

import (
	"context"
	"errors"

	metadataModel "github.com/richardktran/go-movie-microservices/metadata/pkg/model"
	"github.com/richardktran/go-movie-microservices/movie/internal/gateway"
	"github.com/richardktran/go-movie-microservices/movie/pkg/model"
	ratingModel "github.com/richardktran/go-movie-microservices/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type ratingGateway interface {
	GetAggregatedRating(context.Context, ratingModel.RecordID, ratingModel.RecordType) (float64, error)
	// PutRating(context.Context, ratingModel.RecordID, ratingModel.RecordType, *ratingModel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

// Controller defines a movie service controller.
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller.
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)

	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{
		Metadata: *metadata,
	}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(metadata.ID), ratingModel.RecordTypeMovie)

	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}

	return details, nil
}
