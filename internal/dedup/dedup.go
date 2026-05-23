package dedup

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New() (*Client, error) {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		return nil, fmt.Errorf("key unavailable")
	}
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &Client{redis: rdb}, nil
}

func (c *Client) Seen(ctx context.Context, monitorID string, rawURL string) (bool, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false, fmt.Errorf("parse url error: %w", err)
	}

	normalized := parsed.Host + parsed.Path // make sure no url params

	hash := sha256.Sum256([]byte(normalized)) // fixed-length fingerprint

	key := fmt.Sprintf("vigil:seen:%s:%x", monitorID, hash) // formats as hex string

	// setNX sets key if it doesn't exist, and returns true if new or esle false
	set, err := c.redis.SetNX(ctx, key, 1, 7*24*time.Hour).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}
	return !set, nil
}
