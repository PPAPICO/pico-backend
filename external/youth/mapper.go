package youth

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/parser"
)

func toPolicy(item Item) *domain.Policy {
	title := strings.TrimSpace(item.PlcyNm)

	description := strings.TrimSpace(item.PlcyExplnCn)
	if description == "" {
		description = title
	}

	startDate, endDate := parsePolicyDate(
		item.BizPrdBgngYmd,
		item.BizPrdEndYmd,
		item.AplyYmd,
	)

	regionCodeInt, _ := strconv.Atoi(item.ZipCd)
	address := parser.RegionCodeToName(regionCodeInt)
	if address == "" {
		return nil
	}

	return &domain.Policy{
		Title:       title,
		Description: description,
		RegionCode:  parser.ParseRegionCode(item.ZipCd, ""),
		StartDate:   startDate,
		EndDate:     endDate,
		Address:     address,
		Latitude:    0,
		Longitude:   0,
		Condition:   normalizeCondition(item),
	}
}

func parseResponse(body []byte) ([]Item, error) {
	var res response

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Result.YouthPolicyList, nil
}

func parsePolicyDate(
	start string,
	end string,
	applyPeriod string,
) (time.Time, time.Time) {
	startDate := parser.ParseDate(start)
	endDate := parser.ParseDate(end)

	// 사업기간이 없으면 신청기간을 fallback으로 사용
	if startDate.IsZero() || endDate.IsZero() {
		applyStart, applyEnd := parseApplyPeriod(applyPeriod)

		if startDate.IsZero() {
			startDate = applyStart
		}

		if endDate.IsZero() {
			endDate = applyEnd
		}
	}

	// 그래도 없으면 기본값
	now := time.Now()

	if startDate.IsZero() {
		startDate = now
	}

	if endDate.IsZero() {
		endDate = now.AddDate(1, 0, 0)
	}

	return startDate, endDate
}

func parseApplyPeriod(value string) (time.Time, time.Time) {
	value = strings.TrimSpace(value)

	if value == "" {
		return time.Time{}, time.Time{}
	}

	// 예:
	// 20260101~20261231
	// 2026-01-01~2026-12-31
	parts := strings.SplitN(value, "~", 2)
	if len(parts) != 2 {
		return time.Time{}, time.Time{}
	}

	return parser.ParseDate(parts[0]), parser.ParseDate(parts[1])
}
