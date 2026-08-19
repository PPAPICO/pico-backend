package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/janghanul090801/pico-backend/internal/httpclient"
)

type Client struct {
	httpClient *httpclient.Client
	apiKey     string
}

func NewClient(httpClient *httpclient.Client, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

func (c *Client) Chat(ctx context.Context, message string) (string, error) {
	reqBody := map[string]interface{}{
		"model":      "claude-sonnet-5",
		"max_tokens": 1024,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": message,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		// "https://panel.tapie.kr/api/ai-api/v1/chat/completions",
		"https://panel.tapie.kr/api/claude-code/v1/messages",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Req(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf(
			"AI API returned status %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var response ClaudeResponse

	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}
	return response.Content[0].Text, nil
}
