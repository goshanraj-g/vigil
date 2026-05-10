package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type handlers struct {
	pool *pgxpool.Pool
}

func (h *handlers) createMonitor(w http.ResponseWriter, r *http.Request) {}
func (h *handlers) listMonitors(w http.ResponseWriter, r *http.Request) {}
func (h *handlers) getMonitor(w http.ResponseWriter, r *http.Request) {}
func (h *handlers) deleteMonitor(w http.ResponseWriter, r *http.Request) {}
func (h *handlers) getResults(w http.ResponseWriter, r *http.Request) {}
