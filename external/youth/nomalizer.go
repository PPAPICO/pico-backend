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

	return condition
}
