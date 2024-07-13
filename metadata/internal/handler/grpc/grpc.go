package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/richardktran/go-movie-microservices/gen"
	"github.com/richardktran/go-movie-microservices/metadata/internal/controller/metadata"
	"github.com/richardktran/go-movie-microservices/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	log.Printf("GetMetadata request: %v", req)
	if req == nil || req.GetMovieId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty movie id")
	}
	m, err := h.ctrl.Get(ctx, req.GetMovieId())

	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}

func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	log.Printf("PutMetadata request: %v", req)
	if req == nil || req.GetMetadata() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or nil metadata")
	}

	if err := h.ctrl.Put(ctx, req.GetMetadata().GetId(), model.MetadataFromProto(req.GetMetadata())); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PutMetadataResponse{}, nil
}
