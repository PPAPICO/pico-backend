package volunteer

import (
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/policycondition"
)

var conditionParser = policycondition.NewParser()

func normalizeCondition(item DetailItem) domain.PolicyCondition {
	condition := domain.PolicyCondition{}

	if item.YngbgsPosblAt == "N" {
		condition.RequireYouth = new(false)
	}

	text := strings.Join([]string{
		item.ProgrmSj,
		item.ProgrmCn,
	}, " ")

	parsed := conditionParser.Parse(text)

	policycondition.Merge(&condition, parsed)

	return condition
}
