package kakaomap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

func (c *Client) Geocode(ctx context.Context, address string) (latitude, longitude float64, err error) {
	u, err := url.Parse("https://dapi.kakao.com/v2/local/search/address.json")
	if err != nil {
		return 0, 0, fmt.Errorf("parse kakao address url: %w", err)
	}

	q := u.Query()
	q.Set("query", address)
	q.Set("size", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, 0, fmt.Errorf("create kakao address request: %w", err)
	}

	req.Header.Set("Authorization", "KakaoAK "+c.apiKey)

	resp, err := c.httpClient.Req(req)
	if err != nil {
		return 0, 0, fmt.Errorf("request kakao address api: %w", err)
	}
	defer resp.Body.Close()

	var result addressResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, fmt.Errorf("decode kakao address response: %w", err)
	}

	if len(result.Documents) == 0 {
		return 0, 0, fmt.Errorf("address not found: %s", address)
	}

	doc := result.Documents[0]

	longitude, err = strconv.ParseFloat(doc.X, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse longitude: %w", err)
	}

	latitude, err = strconv.ParseFloat(doc.Y, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse latitude: %w", err)
	}

	return latitude, longitude, nil
}
