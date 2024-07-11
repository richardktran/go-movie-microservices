package grpc

import (
	"context"
	"errors"

	"github.com/richardktran/go-movie-microservices/gen"
	"github.com/richardktran/go-movie-microservices/metadata/pkg/model"
	"github.com/richardktran/go-movie-microservices/movie/internal/controller/movie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.GetMovieId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty movie id")
	}

	m, err := h.ctrl.Get(ctx, req.GetMovieId())

	if err != nil && errors.Is(err, movie.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Rating:   float32(*m.Rating),
			Metadata: model.MetadataToProto(&m.Metadata),
		},
	}, nil
}
