package volunteer

import (
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/policycondition"
)

var conditionParser = policycondition.NewParser()

func normalizeCondition(item DetailItem) domain.PolicyCondition {
	condition := domain.PolicyCondition{}

	text := strings.Join([]string{
		item.ProgrmSj,
		item.ProgrmCn,
	}, " ")

	parsed := conditionParser.Parse(text)

	policycondition.Merge(&condition, parsed)

	// YngbgsPosblAt 은 미성년자 참여가능 여부 (Y: 미성년자 가능, N: 성인만 가능)
	if item.YngbgsPosblAt == "N" {
		nineteen := 19
		if condition.MinAge == nil || *condition.MinAge < nineteen {
			condition.MinAge = &nineteen
		}
	}

	return condition
}
