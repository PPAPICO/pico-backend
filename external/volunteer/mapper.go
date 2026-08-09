package volunteer

import (
	"fmt"
	"strings"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/parser"
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

func toPolicy(item DetailItem) *domain.Policy {
	return &domain.Policy{
		Title:       item.ProgrmSj,
		Description: buildDescription(item),
		RegionCode:  parser.ParseRegionCode(item.SidoCd, item.GugunCd),
		StartDate:   parser.ParseYYYYMMDD(item.ProgrmBgnde),
		EndDate:     parser.ParseYYYYMMDD(item.ProgrmEndde),
		Address: strings.TrimSpace(
			item.AreaAddress1 + " " +
				item.AreaAddress2 + " " +
				item.AreaAddress3,
		),
		Latitude:  parser.ParseFloat(item.AreaLalo1),
		Longitude: parser.ParseFloat(item.AreaLalo2),
	}
}
