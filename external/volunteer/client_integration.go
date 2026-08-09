package volunteer

import (
	"context"
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

// Fetch fetches volunteer activities and events from 행정안전부 봉사참여정보서비스
func (c *Client) Fetch(ctx context.Context) ([]*domain.Policy, error) {
	list, err := c.list(ctx)
	if err != nil {
		return nil, err
	}

	if config.E.AppEnv == "test" && len(list) > 3 {
		list = list[:3]
	}

	var (
		mu     sync.Mutex
		result = make([]*domain.Policy, 0, len(list))
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(1)

	for _, item := range list {
		g.Go(func() error {
			exist, _ := c.policyRepository.FindByTitle(ctx, item.ProgrmSj)
			if exist != nil {
				return nil
			}

			if err := delay.SleepContext(ctx, 100*time.Millisecond); err != nil {
				return err
			}
			policy, err := c.detail(ctx, item.ProgrmRegistNo)
			if err != nil {
				log.Printf("volunteer detail %s: %v", item.ProgrmRegistNo, err)
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

func (c *Client) list(ctx context.Context) ([]AreaItem, error) {
	params := url.Values{}
	params.Set("serviceKey", config.E.VolunteerApiKey)
	params.Set("pageNo", "1")
	params.Set("numOfRows", "50")
	params.Set("schSido", "6110000") // 서울

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/1741000/volunteerPartcptnService/getVltrAreaList?%s",
		params.Encode(),
	)

	body, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("request volunteer area list: %w", err)
	}

	var resp AreaResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode volunteer area list: %w", err)
	}

	if resp.Header.ResultCode != "00" {
		return nil, fmt.Errorf("api error: %s (%s)",
			resp.Header.ResultMsg,
			resp.Header.ResultCode,
		)
	}

	return resp.Body.Items.Item, nil
}

func (c *Client) detail(
	ctx context.Context,
	registNo string,
) (*domain.Policy, error) {

	params := url.Values{}
	params.Set("serviceKey", config.E.VolunteerApiKey)
	params.Set("pageNo", "1")
	params.Set("numOfRows", "1")
	params.Set("progrmRegistNo", registNo)

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/1741000/volunteerPartcptnService/getVltrPartcptnItem?%s",
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

	return toPolicy(resp.Body.Items.Item), nil
}
