package welfare

import (
	"strings"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

func buildDescription(item DetailItem) string {
	var desc strings.Builder

	desc.WriteString(strings.TrimSpace(item.WlfareInfoOutlCn))

	if item.AlwServCn != "" {
		desc.WriteString("\n\n지원내용\n")
		desc.WriteString(item.AlwServCn)
	}

	if item.TgtrDtlCn != "" {
		desc.WriteString("\n\n지원대상\n")
		desc.WriteString(item.TgtrDtlCn)
	}

	if item.SlctCritCn != "" {
		desc.WriteString("\n\n선정기준\n")
		desc.WriteString(item.SlctCritCn)
	}

	if len(item.ApplmetList) > 0 {
		desc.WriteString("\n\n신청방법")
		for _, v := range item.ApplmetList {
			desc.WriteString("\n- ")
			desc.WriteString(v.ServSeDetailNm)
		}
	}

	return desc.String()
}

func toPolicy(item DetailItem) *domain.Policy {
	address := item.JurMnofNm
	if item.RprsCtadr != "" {
		address = item.RprsCtadr
	}
	return &domain.Policy{
		Title:       strings.TrimSpace(item.ServNm),
		Description: buildDescription(item),
		RegionCode:  0, // 중앙부처
		StartDate:   time.Now(),
		EndDate:     time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		Address:     address,
		Latitude:    0,
		Longitude:   0,
	}
}
