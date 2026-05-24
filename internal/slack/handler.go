package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) SlashCommand(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respond(w, "failed to parse request")
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	parts := strings.Fields(text)

	if len(parts) == 0 {
		respond(w, "Usage:\n• `/monitor <query>` - start monitoring\n• `/monitor list` - list active monitors\n• `/monitor stop <n>` - stop monitor #n")
		return
	}

	switch parts[0] {
	case "list":
		h.list(w, r)
	case "stop":
		if len(parts) < 2 {
			respond(w, "Usage: `/monitor stop <n>`")
			return
		}
		h.stop(w, r, parts[1])
	default:
		h.create(w, r, text)
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	rows, err := h.pool.Query(r.Context(),
		`SELECT id, query FROM monitors WHERE status = 'active' ORDER BY created_at ASC`)
	if err != nil {
		respond(w, "Database error")
		return
	}
	defer rows.Close()

	var lines []string
	i := 1
	for rows.Next() {
		var id, query string
		if err := rows.Scan(&id, &query); err != nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%d. *%s*", i, query))
		i++
	}

	if len(lines) == 0 {
		respond(w, "No active monitors")
		return
	}

	respond(w, "*Active monitors:*\n"+strings.Join(lines, "\n"))
}

func (h *Handler) stop(w http.ResponseWriter, r *http.Request, nStr string) {
	n, err := strconv.Atoi(nStr)
	if err != nil || n < 1 {
		respond(w, "Invalid number, use `/monitor list` to see monitor numbers")
		return
	}

	rows, err := h.pool.Query(r.Context(),
		`SELECT id, query FROM monitors WHERE status = 'active' ORDER BY created_at ASC`)
	if err != nil {
		respond(w, "Database error")
		return
	}
	defer rows.Close()

	type monitor struct{ id, query string }
	var monitors []monitor
	for rows.Next() {
		var m monitor
		if err := rows.Scan(&m.id, &m.query); err != nil {
			continue
		}
		monitors = append(monitors, m)
	}

	if n > len(monitors) {
		respond(w, fmt.Sprintf("No monitor #%d, use `/monitor list` to see active monitors", n))
		return
	}

	target := monitors[n-1]
	_, err = h.pool.Exec(r.Context(), `DELETE FROM monitors WHERE id = $1`, target.id)
	if err != nil {
		respond(w, "Database error")
		return
	}

	respond(w, fmt.Sprintf("Stopped: *%s*", target.query))
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request, query string) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		respond(w, "SLACK_WEBHOOK_URL not configured")
		return
	}

	_, err := h.pool.Exec(r.Context(),
		`INSERT INTO monitors (query, webhook_url, schedule) VALUES ($1, $2, $3)`,
		query, webhookURL, "@every 1m",
	)
	if err != nil {
		respond(w, "Failed to create monitor")
		return
	}

	respond(w, fmt.Sprintf("Monitoring: *%s*", query))
}

func respond(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"text": text})
}
