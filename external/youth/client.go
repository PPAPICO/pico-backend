package youth

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
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

// Fetch fetches policies from 온통청년 청년정책 API.
func (c *Client) Fetch(ctx context.Context) ([]*domain.Policy, error) {
	params := url.Values{}
	params.Set("apiKeyNm", c.apiKey)
	params.Set("pageNum", "1")
	params.Set("pageSize", "100")
	params.Set("pageType", "1")
	params.Set("rtnType", "json")

	reqURL := "https://www.youthcenter.go.kr/go/ythip/getPlcy?" + params.Encode()

	body, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("youth api request failed: %w", err)
	}

	items, err := parseResponse(body)
	if err != nil {
		return nil, fmt.Errorf("youth api response parse failed: %w", err)
	}

	result := make([]*domain.Policy, 0, len(items))

	for _, item := range items {
		title := strings.TrimSpace(item.PlcyNm)
		if title == "" {
			continue
		}

		result = append(result, toPolicy(item))
	}

	return result, nil
}
