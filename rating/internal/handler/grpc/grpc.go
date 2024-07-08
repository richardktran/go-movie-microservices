package grpc

import (
	"context"
	"errors"

	"github.com/richardktran/go-movie-microservices/gen"
	"github.com/richardktran/go-movie-microservices/rating/internal/controller/rating"
	"github.com/richardktran/go-movie-microservices/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	*gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.GetRecordId() == "" || req.GetRecordType() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record id or record type")
	}

	v, err := h.ctrl.GetAggregation(ctx, model.RecordID(req.GetRecordId()), model.RecordType(req.GetRecordType()))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req != nil || req.GetRecordId() == "" || req.GetRecordType() == "" || req.GetRatingValue() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record id or record type or rating value")
	}

	if err := h.ctrl.PutRating(ctx,
		model.RecordID(req.GetRecordId()),
		model.RecordType(req.GetRecordType()),
		&model.Rating{
			UserID: model.UserID(req.GetUserId()),
			Value:  model.RatingValue(req.GetRatingValue()),
		}); err != nil {
		return nil, err
	}

	return &gen.PutRatingResponse{}, nil
}
