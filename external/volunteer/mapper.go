package volunteer

import (
	"fmt"
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/parser"
)

func buildDescription(item DetailItem) string {
	return fmt.Sprintf(
		`%s

주관기관 : %s
모집기관 : %s
모집상태 : %s
모집인원 : %d명
신청인원 : %d명
활동요일 : %s`,
		strings.TrimSpace(item.ProgrmCn),
		item.NanmmbyNm,
		item.MnnstNm,
		item.ProgrmSttusSe,
		item.RcritNmpr,
		item.AppTotal,
		item.ActWkdy,
	)
}

func getAddress(item DetailItem) string {
	addr := strings.TrimSpace(item.PostAdres)
	if addr != "" {
		return addr
	}
	parts := []string{}
	if a1 := strings.TrimSpace(item.AreaAddress1); a1 != "" {
		parts = append(parts, a1)
	}
	if a2 := strings.TrimSpace(item.AreaAddress2); a2 != "" {
		parts = append(parts, a2)
	}
	if a3 := strings.TrimSpace(item.AreaAddress3); a3 != "" {
		parts = append(parts, a3)
	}
	return strings.Join(parts, " ")
}

func ToPolicy(item DetailItem) *domain.Policy {
	return &domain.Policy{
		Title:       item.ProgrmSj,
		Description: buildDescription(item),
		RegionCode:  parser.ParseRegionCode(item.SidoCd, item.GugunCd),
		StartDate:   parser.ParseYYYYMMDD(item.ProgrmBgnde),
		EndDate:     parser.ParseYYYYMMDD(item.ProgrmEndde),
		Address:     getAddress(item),
		Latitude:    parser.ParseFloat(item.AreaLalo1),
		Longitude:   parser.ParseFloat(item.AreaLalo2),
		Condition:   normalizeCondition(item),
	}
}
