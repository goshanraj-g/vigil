package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	h := &handlers{pool: pool} // Create a pool of handlers which have access to the database pool
	r := chi.NewRouter() // Router

	r.Post("/monitors", h.createMonitor)
	r.Get("/monitors", h.listMonitors)
	r.Get("/monitors/{id}", h.getMonitor)
	r.Delete("/monitors/{id}", h.deleteMonitor)
	r.Get("/monitors/{id}/results", h.getResults)

	return r
}