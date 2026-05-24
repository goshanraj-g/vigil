package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/goshanraj-g/vigil/internal/search"
)

type Client struct {
	ai anthropic.Client
}

func New() (*Client, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}
	return &Client{
		ai: anthropic.NewClient(option.WithAPIKey(key)), // anthropic SDK handles HTTP, retries and auth
	}, nil
}

func (c *Client) IsRelevant(ctx context.Context, query string, result search.Result) (bool, error) {
	prompt := fmt.Sprintf(`You are a relevance filter for a web monitoring system.

The user is monitoring the web for: %q

A search result was found:
Title: %s
URL: %s
Snippet: %s

Is this result relevant to what the user is monitoring for?
Reply with only YES or NO. No explanation.`, query, result.Title, result.URL, result.Snippet)

	msg, err := c.ai.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 5,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return false, fmt.Errorf("claude: %w", err)
	}

	if len(msg.Content) == 0 {
		return false, fmt.Errorf("claude: empty response")
	}

	text := strings.TrimSpace(strings.ToUpper(msg.Content[0].Text))
	return text == "YES", nil

}

func (c *Client) Summarize(ctx context.Context, query string, content string) (string, error) {
	prompt := fmt.Sprintf(`You are a summarizer for a web monitoring system.

The user is monitoring the web for: %q

The following article was found and deemed relevant. Write a 2-3 sentence summary
of the article, focusing on why it is relevant to what the user is monitoring for.
Be concise and factual.

Article content:
%s`, query, content)

	msg, err := c.ai.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 200,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", fmt.Errorf("claude: %w", err)
	}

	if len(msg.Content) == 0 {
		return "", fmt.Errorf("claude: empty respone")
	}

	return strings.TrimSpace(msg.Content[0].Text), nil
}
