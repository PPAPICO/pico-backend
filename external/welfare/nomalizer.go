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

	condition := conditionParser.Parse(text)

	// 	DetailItem 자체엔 링크 필드가 없고, ApplmetList(신청방법 목록) 안에
	// ServSeDetailLink로 들어있음. 첫 번째 항목의 링크를 대표로 사용.
	if len(item.ApplmetList) > 0 && item.ApplmetList[0].ServSeDetailLink != "" {
		condition.SourceURL = item.ApplmetList[0].ServSeDetailLink
	}

	return condition
}