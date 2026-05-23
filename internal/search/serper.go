package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Result struct {
	URL     string
	Title   string
	Snippet string
}

type Client struct {
	apiKey string
	http   *http.Client
}

type serperResponse struct {
	Organic []struct {
		Link    string `json:"link"`
		Title   string `json:"title"`
		Snippet string `json:"snippet"`
	} `json:"organic"`
}

func New() (*Client, error) {
	key := os.Getenv("SERPER_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("SERPER_API_KEY not set")
	}
	return &Client{
		apiKey: key,
		http:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *Client) Search(ctx context.Context, query string) ([]Result, error) {
	if query == "" {
		return nil, fmt.Errorf("query must not be empty")
	}

	// build JSON req body
	body, err := json.Marshal(map[string]string{"q": query})
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	// build JSON req
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://google.serper.dev/search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	// send the req, check if resp is nil since .close would panic
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	// check status + decod
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serper status %d", resp.StatusCode)
	}

	var sr serperResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	// serper's shape -> result type
	out := make([]Result, 0, len(sr.Organic))
	for _, r := range sr.Organic {
		if r.Link == "" {
			continue
		}
		out = append(out, Result{URL: r.Link, Title: r.Title, Snippet: r.Snippet})
	}
	return out, nil
}
