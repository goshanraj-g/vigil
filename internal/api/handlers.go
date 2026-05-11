package api

import (
	"encoding/json"
	"net/http"

	"github.com/goshanraj-g/vigil/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type handlers struct {
	pool *pgxpool.Pool
}

type createMonitorRequest struct {
	Query      string `json:"query"`
	WebhookURL string `json:"webhood_url"`
	Schedule   string `json:"schedule"`
}

func (h *handlers) createMonitor(w http.ResponseWriter, r *http.Request) {
	var req createMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Query == "" || req.WebhookURL == "" || req.Schedule == "" {
		http.Error(w, "query, webhook_url, and schedule are required", http.StatusBadRequest)
		return
	}

	var m models.Monitor // Declare empty monitor
	err := h.pool.QueryRow(r.Context(),
		`INSERT INTO monitors (query, webhook_url, schedule)
		VALUES ($1, $2, $3)
		RETURNING id, query, webhook_url, schedule, status, created_at`,
		req.Query, req.WebhookURL, req.Schedule, // Sends the SQL to Postgres, $1,$2,$3 subs with req.Query...
	).Scan(&m.ID, &m.Query, &m.WebhookURL, &m.Schedule, &m.Status, &m.CreatedAt) // Scans for values to returning
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func (h *handlers) listMonitors(w http.ResponseWriter, r *http.Request) {
	rows, err := h.pool.Query(r.Context(),
		`SELECT id, query, webhook_url, schedule, status, created_at FROM monitors ORDER BY created_at DESC`,
	)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	monitors := []models.Monitor{}
	for rows.Next() {
		var m models.Monitor
		if err := rows.Scan(&m.ID, &m.Query, &m.WebhookURL, &m.Schedule, &m.Status, &m.CreatedAt); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		monitors = append(monitors, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monitors)
}

func (h *handlers) getMonitor(w http.ResponseWriter, r *http.Request)    {}
func (h *handlers) deleteMonitor(w http.ResponseWriter, r *http.Request) {}
func (h *handlers) getResults(w http.ResponseWriter, r *http.Request)    {}
