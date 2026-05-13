package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goshanraj-g/vigil/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type handlers struct {
	pool *pgxpool.Pool
}

type createMonitorRequest struct {
	Query      string `json:"query"`
	WebhookURL string `json:"webhook_url"`
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

func (h *handlers) getMonitor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var m models.Monitor
	err := h.pool.QueryRow(r.Context(),
		`SELECT id, query, webhook_url, schedule, status, created_at FROM monitors WHERE id = $1`,
		id).Scan(&m.ID, &m.Query, &m.WebhookURL, &m.Schedule, &m.Status, &m.CreatedAt)
	if err == pgx.ErrNoRows {
		http.Error(w, "monitor not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
func (h *handlers) deleteMonitor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.pool.Exec(r.Context(),
		`DELETE FROM monitors WHERE id = $1`, id)

	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected() == 0 {
		http.Error(w, "monitor not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *handlers) getResults(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rows, err := h.pool.Query(r.Context(),
		`SELECT id, monitor_id, url, title, snippet, found_at FROM results WHERE monitor_id = $1 ORDER BY found_at DESC`,
		id,
	)

	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []models.Result{}
	for rows.Next() {
		var res models.Result
		if err := rows.Scan(&res.ID, &res.MonitorID, &res.URL, &res.Title, &res.Snippet, &res.FoundAt); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		results = append(results, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)

}
