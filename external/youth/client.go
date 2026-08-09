package youth

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/httpclient"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/parser"
)

type Client struct {
	httpClient *httpclient.Client
	apiKey     string
}

// Fetch fetches policies from 온통청년 청년정책 API
func (c *Client) Fetch(ctx context.Context) ([]*domain.Policy, error) {
	apiKey := config.E.YouthApiKey

	reqURL := fmt.Sprintf("https://www.youthcenter.go.kr/opi/empSprtList.do?openApiVcntId=%s&pageIndex=1&display=20", url.QueryEscape(apiKey))

	body, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	var items []Item

	var xmlResp XMLResponse
	if xmlErr := xml.Unmarshal(body, &xmlResp); xmlErr == nil && len(xmlResp.PolicyList) > 0 {
		items = xmlResp.PolicyList
	} else {
		var jsonResp JSONResponse
		if jsonErr := json.Unmarshal(body, &jsonResp); jsonErr == nil && len(jsonResp.EmpsInfo.Emp) > 0 {
			items = jsonResp.EmpsInfo.Emp
		}
	}

	var result []*domain.Policy
	now := time.Now()

	for _, item := range items {
		title := strings.TrimSpace(item.PolyBizSjnNm)
		if title == "" {
			continue
		}
		desc := strings.TrimSpace(item.PolyItcnCn)
		if desc == "" {
			desc = title
		}

		sDate, eDate := parser.ParseDateRange(item.RqstPrdCn)
		if sDate.IsZero() {
			sDate = now
		}
		if eDate.IsZero() {
			eDate = now.AddDate(1, 0, 0)
		}

		regionCode := parser.ParseRegionCode(item.PolyBizSecd, item.CnsgNtiPrdCn)

		result = append(result, &domain.Policy{
			Title:       title,
			Description: desc,
			RegionCode:  regionCode,
			StartDate:   sDate,
			EndDate:     eDate,
			Address:     strings.TrimSpace(item.CnsgNtiPrdCn),
			Latitude:    0.0,
			Longitude:   0.0,
		})
	}

	return result, nil
}
