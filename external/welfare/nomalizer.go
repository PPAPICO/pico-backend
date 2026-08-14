package welfare

import (
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/policycondition"
)

var conditionParser = policycondition.NewParser()

func normalizeCondition(item DetailItem) domain.PolicyCondition {
	text := strings.Join([]string{
		item.TgtrDtlCn,
		item.SlctCritCn,
	}, " ")

	return conditionParser.Parse(text)
}
