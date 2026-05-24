package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func New() (*Client, error) {
	key := os.Getenv("PARALLEL_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("PARALLEL_API_KEY not set")
	}
	return &Client{
		apiKey: key,
		http:   &http.Client{Timeout: 30 * time.Second},
	}, nil
}

type parallelResponse struct {
	Results []struct {
		Excerpts []string `json:"excerpts"`
	} `json:"results"`
}

func (c *Client) Scrape(ctx context.Context, url string, objective string) (string, error) {
	body, err := json.Marshal(map[string]any{
		"urls":      []string{url},
		"objective": objective,
	})
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.parallel.ai/v1/extract", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("parallel status %d", resp.StatusCode)
	}

	var pr parallelResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	if len(pr.Results) == 0 || len(pr.Results[0].Excerpts) == 0 {
		return "", fmt.Errorf("no excerpts returned")
	}

	// Join all excerpts into one string for Claude to summarize
	// join exerpts into one string for summarization
	return strings.Join(pr.Results[0].Excerpts, "\n\n"), nil
}
