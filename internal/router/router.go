package router

import (
	"github.com/Sacro/urlshortner/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(repo handlers.HandlerRepository) chi.Router {
	r := chi.NewRouter()

	r.Post("/", repo.CreateHandler)
	r.Get("/{id}", repo.RetrieveHandler)

	return r
}
