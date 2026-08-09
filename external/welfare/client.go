package welfare

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/delay"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/httpclient"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	httpClient       *httpclient.Client
	policyRepository domain.PolicyRepository
	apiKey           string
}

// Fetch fetches welfare services from 한국사회보장정보원 (중앙부처복지서비스)
func (c *Client) Fetch(ctx context.Context) ([]*domain.Policy, error) {
	items, err := c.list(ctx)
	if err != nil {
		return nil, err
	}

	if config.E.AppEnv == "test" && len(items) > 3 {
		items = items[:3]
	}

	var (
		mu     sync.Mutex
		result = make([]*domain.Policy, 0, len(items))
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(1)

	for _, item := range items {
		g.Go(func() error {
			exist, _ := c.policyRepository.FindByTitle(ctx, item.ServNm)
			if exist != nil {
				return nil
			}

			if err := delay.SleepContext(ctx, 100*time.Millisecond); err != nil {
				return err
			}
			policy, err := c.detail(ctx, item.ServID)
			if err != nil {
				log.Printf("welfare detail %s: %v", item.ServID, err)
				return nil
			}

			mu.Lock()
			result = append(result, policy)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) list(ctx context.Context) ([]Item, error) {
	params := url.Values{}
	params.Set("serviceKey", config.E.WelfareApiKey)
	params.Set("callTp", "L")
	params.Set("pageNo", "1")
	params.Set("numOfRows", "100")
	params.Set("srchKeyCode", "001")

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/B554287/NationalWelfareInformationsV001/NationalWelfarelistV001?%s",
		params.Encode(),
	)

	body, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("request welfare list: %w", err)
	}

	// XML 우선
	var xmlResp XMLResponse
	if err := xml.Unmarshal(body, &xmlResp); err == nil {
		if xmlResp.ResultCode != "" && xmlResp.ResultCode != "0" {
			return nil, fmt.Errorf(
				"api error: %s (%s)",
				xmlResp.ResultMessage,
				xmlResp.ResultCode,
			)
		}

		if len(xmlResp.ServList) > 0 {
			return xmlResp.ServList, nil
		}
	}

	// JSON fallback
	var jsonResp JSONResponse
	if err := json.Unmarshal(body, &jsonResp); err == nil {
		if jsonResp.WantedList.ResultCode != "" &&
			jsonResp.WantedList.ResultCode != "0" {
			return nil, fmt.Errorf(
				"api error: %s (%s)",
				jsonResp.WantedList.ResultMessage,
				jsonResp.WantedList.ResultCode,
			)
		}

		return jsonResp.WantedList.ServList, nil
	}

	return nil, fmt.Errorf("failed to decode welfare response")
}

func (c *Client) detail(
	ctx context.Context,
	servID string,
) (*domain.Policy, error) {

	params := url.Values{}
	params.Set("serviceKey", config.E.WelfareApiKey)
	params.Set("callTp", "D")
	params.Set("servId", servID)

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/B554287/NationalWelfareInformationsV001/NationalWelfaredetailedV001?%s",
		params.Encode(),
	)

	body, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	var resp DetailResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return toPolicy(resp.DetailItem), nil
}
