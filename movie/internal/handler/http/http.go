package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/richardktran/go-movie-microservices/movie/internal/controller/movie"
)

type Handler struct {
	ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	details, err := h.ctrl.Get(req.Context(), id)

	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}

	// Show details in the console
	log.Printf("Movie details: %v\n", details)
}
