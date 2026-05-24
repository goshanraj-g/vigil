package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goshanraj-g/vigil/internal/ai"
	"github.com/goshanraj-g/vigil/internal/dedup"
	"github.com/goshanraj-g/vigil/internal/scraper"
	"github.com/goshanraj-g/vigil/internal/search"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool    *pgxpool.Pool
	search  *search.Client
	ai      *ai.Client
	dedup   *dedup.Client
	scraper *scraper.Client
	http    *http.Client
}

func New(pool *pgxpool.Pool, s *search.Client, a *ai.Client, d *dedup.Client, sc *scraper.Client) *Handler {
	return &Handler{pool: pool, search: s, ai: a, dedup: d, scraper: sc, http: &http.Client{Timeout: 10 * time.Second}}
}

func (h *Handler) Handle(ctx context.Context, t *asynq.Task) error {
	var payload SearchMonitorPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	var query, webhookURL string

	err := h.pool.QueryRow(ctx,
		`SELECT query, webhook_url FROM monitors WHERE id = $1`,
		payload.MonitorID,
	).Scan(&query, &webhookURL)
	if err != nil {
		return fmt.Errorf("fetch monitor: %w", err)
	}

	results, err := h.search.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, r := range results {
		seen, err := h.dedup.Seen(ctx, payload.MonitorID, r.URL)
		if err != nil {
			return fmt.Errorf("dedup: %w", err)
		}
		if seen {
			continue
		}
		isRelevant, err := h.ai.IsRelevant(ctx, query, r)
		if err != nil {
			return fmt.Errorf("claude: %w", err)
		}
		if !isRelevant {
			continue
		}

		content, err := h.scraper.Scrape(ctx, r.URL)
		if err != nil {
			content = ""
		}

		summary := ""
		if content != "" {
			summary, err = h.ai.Summarize(ctx, query, content)
			if err != nil {
				summary = ""
			}
		}

		_, err = h.pool.Exec(ctx,
			`INSERT INTO results (monitor_id, url, title, snippet, summary)
			VALUES ($1, $2, $3, $4, $5)`,
			payload.MonitorID, r.URL, r.Title, r.Snippet, summary,
		)
		if err != nil {
			return fmt.Errorf("insert result: %w", err)
		}

		body, err := json.Marshal(map[string]string{
			"url":     r.URL,
			"title":   r.Title,
			"snippet": r.Snippet,
			"summary": summary,
		})
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("new request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := h.http.Do(req)
		if err != nil {
			return fmt.Errorf("do: %w", err)
		}
		defer resp.Body.Close()
	}
	return nil
}
