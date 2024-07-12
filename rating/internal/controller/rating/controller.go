package rating

import (
	"context"
	"errors"
	"fmt"

	"github.com/richardktran/go-movie-microservices/rating/internal/repository"
	"github.com/richardktran/go-movie-microservices/rating/pkg/model"
)

// ErrNotFound is returned when not ratings are found for a record.
var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(context.Context, model.RecordID, model.RecordType) ([]model.Rating, error)
	Put(context.Context, model.RecordID, model.RecordType, *model.Rating) error
}

type ratingIngester interface {
	Ingest(context.Context) (chan model.RatingEvent, error)
}

// Controller defines a rating service controller.
type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

// New creates a new rating service controller.
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo, ingester}
}

// GetAggregation returns the aggregation rating for a record or ErrNotFound if there are no rating for it.
func (c *Controller) GetAggregation(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)
	totalRatings := len(ratings)

	for _, r := range ratings {
		sum += float64(r.Value)
	}

	return sum / float64(totalRatings), nil
}

// PutRating writes a rating for a given record.
func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}

// StartIngestion starts the ingestion of rating events.
func (c *Controller) StartIngestion(ctx context.Context) error {
	ch, err := c.ingester.Ingest(ctx)
	if err != nil {
		return err
	}

	for e := range ch {
		fmt.Printf("Rating event received: %v\n", e)
		if err := c.PutRating(ctx, e.RecordID, e.RecordType, &model.Rating{
			Value:  e.Value,
			UserID: e.UserID,
		}); err != nil {
			return err
		}
	}
	return nil
}
