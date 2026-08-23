package youth

import (
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/policycondition"
)

var conditionParser = policycondition.NewParser()

func normalizeCondition(item Item) domain.PolicyCondition {
	condition := domain.PolicyCondition{}

	condition.MinAge = policycondition.ParseAge(
		item.SprtTrgtMinAge,
	)

	condition.MaxAge = policycondition.ParseAge(
		item.SprtTrgtMaxAge,
	)

	text := strings.Join([]string{
		item.PlcyNm,
		item.PlcyKywdNm,
		item.PlcyExplnCn,
		item.PlcySprtCn,
	}, " ")

	parsed := conditionParser.Parse(text)

	policycondition.Merge(&condition, parsed)

	// 신청 페이지 URL 저장 (없으면 참고 링크로 대체)
	if item.AplyUrlAddr != "" {
		condition.SourceURL = item.AplyUrlAddr
	} else if item.RefUrlAddr1 != "" {
		condition.SourceURL = item.RefUrlAddr1
	}

	return condition
}
