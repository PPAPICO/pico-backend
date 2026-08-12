package policycondition

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/janghanul090801/pico-backend/domain"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(text string) domain.PolicyCondition {
	var condition domain.PolicyCondition

	text = strings.TrimSpace(text)

	if text == "" {
		return condition
	}

	p.parseAge(text, &condition)
	p.parseKeywords(text, &condition)

	return condition
}

var (
	minAgeRegex = regexp.MustCompile(
		`(?:만\s*)?(\d+)\s*세\s*(?:이상|부터)`,
	)

	maxAgeRegex = regexp.MustCompile(
		`(?:만\s*)?(\d+)\s*세\s*(?:이하|미만|까지)`,
	)

	rangeAgeRegex = regexp.MustCompile(
		`(?:만\s*)?(\d+)\s*세?\s*(?:~|-|부터)\s*(?:만\s*)?(\d+)\s*세?`,
	)
)

func (p *Parser) parseAge(
	text string,
	condition *domain.PolicyCondition,
) {
	if matches := rangeAgeRegex.FindStringSubmatch(text); len(matches) == 3 {
		minAge, err1 := strconv.Atoi(matches[1])
		maxAge, err2 := strconv.Atoi(matches[2])

		if err1 == nil {
			condition.MinAge = &minAge
		}

		if err2 == nil {
			condition.MaxAge = &maxAge
		}

		return
	}

	if matches := minAgeRegex.FindStringSubmatch(text); len(matches) == 2 {
		age, err := strconv.Atoi(matches[1])

		if err == nil {
			condition.MinAge = &age
		}
	}

	if matches := maxAgeRegex.FindStringSubmatch(text); len(matches) == 2 {
		age, err := strconv.Atoi(matches[1])

		if err == nil {
			condition.MaxAge = &age
		}
	}
}

func (p *Parser) parseKeywords(
	text string,
	condition *domain.PolicyCondition,
) {
	if containsAny(text,
		"청년",
		"청년층",
		"청년인",
	) {
		condition.RequireYouth = new(true)
	}

	if containsAny(text,
		"임산부",
		"임신부",
		"임신 여성",
		"임신 중",
	) {
		condition.RequirePregnant = new(true)
	}

	if containsAny(text,
		"장애인",
		"등록 장애인",
		"등록장애인",
	) {
		condition.RequireDisabled = new(true)
	}

	if containsAny(text,
		"외국인",
		"다문화",
		"다문화가족",
	) {
		condition.RequireForeign = new(true)
	}

	if containsAny(text,
		"대학생",
		"대학 재학생",
		"재학생",
	) {
		condition.RequireStudent = new(true)
	}

	if containsAny(text,
		"사업자",
		"소상공인",
		"창업자",
	) {
		condition.RequireBusiness = new(true)
	}
}

func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}
