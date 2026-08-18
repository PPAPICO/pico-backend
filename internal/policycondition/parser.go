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
	// 전화번호 패턴 제외 (예: 010-1234-5678, 02-123-4567, 070-123-4567 등)
	phoneRegex = regexp.MustCompile(
		`\b\d{2,4}-\d{3,4}-\d{4}\b`,
	)

	// 날짜/시간 범위 제외 (예: 09:00~18:00, 2024-01-01)
	timeRangeRegex = regexp.MustCompile(
		`\b\d{1,2}:\d{2}\s*[-~]\s*\d{1,2}:\d{2}\b|\b\d{4}[-./]\d{1,2}[-./]\d{1,2}\b`,
	)

	// 나이 범위 패턴: 예) 만 19세 ~ 34세, 19~34세, 만19세-만34세, 19세부터 34세까지
	rangeAgeRegex = regexp.MustCompile(
		`(?:^|[^가-힣0-9])(?:만|연)?\s*(\d{1,3})\s*(?:세|대)?\s*(?:~|-|부터|\s*이상\s*~|\s*부터\s*)\s*(?:만|연)?\s*(\d{1,3})\s*(세|대)?(?:\s*(?:이하|까지|미만))?(?:[^가-힣0-9]|$)`,
	)

	// 최소 나이 조건: 예) 만 19세 이상, 19세부터, 연 18세 초과
	minAgeRegex = regexp.MustCompile(
		`(?:만|연)?\s*(\d{1,3})\s*(세|대)?\s*(이상|부터|초과)`,
	)

	// 최대 나이 조건: 예) 34세 이하, 만 34세 미만, 35세 미만
	maxAgeRegex = regexp.MustCompile(
		`(?:만|연)?\s*(\d{1,3})\s*(세|대)?\s*(이하|미만|까지)`,
	)
)

type ageCondition struct {
	min *int
	max *int
	raw string
}

func isValidAge(age int) bool {
	return age >= 0 && age <= 99
}

func (p *Parser) parseAge(
	text string,
	condition *domain.PolicyCondition,
) {
	cleanText := phoneRegex.ReplaceAllString(text, "")
	cleanText = timeRangeRegex.ReplaceAllString(cleanText, "")

	var candidates []ageCondition

	// 1. 나이 범위 조건
	rangeMatches := rangeAgeRegex.FindAllStringSubmatchIndex(cleanText, -1)
	for _, match := range rangeMatches {
		if len(match) < 6 {
			continue
		}

		minAge, err1 := strconv.Atoi(capture(cleanText, match, 1))
		maxAge, err2 := strconv.Atoi(capture(cleanText, match, 2))

		if err1 != nil || err2 != nil {
			continue
		}

		unit := capture(cleanText, match, 3)
		if unit == "대" {
			minAge = (minAge / 10) * 10
			maxAge = (maxAge/10)*10 + 9
		}

		if !isValidAge(minAge) || !isValidAge(maxAge) {
			continue
		}

		candidates = append(candidates, ageCondition{
			min: &minAge,
			max: &maxAge,
			raw: strings.TrimSpace(cleanText[match[0]:match[1]]),
		})
	}

	// 2. 최소 나이 조건
	minMatches := minAgeRegex.FindAllStringSubmatchIndex(cleanText, -1)
	for _, match := range minMatches {
		if len(match) < 8 {
			continue
		}

		age, err := strconv.Atoi(capture(cleanText, match, 1))
		if err != nil {
			continue
		}

		unit := capture(cleanText, match, 2)
		operator := capture(cleanText, match, 3)

		if unit == "대" {
			age = (age / 10) * 10
		} else if operator == "초과" {
			age += 1
		}

		if !isValidAge(age) {
			continue
		}

		candidates = append(candidates, ageCondition{
			min: &age,
			raw: strings.TrimSpace(cleanText[match[0]:match[1]]),
		})
	}

	// 3. 최대 나이 조건
	maxMatches := maxAgeRegex.FindAllStringSubmatchIndex(cleanText, -1)
	for _, match := range maxMatches {
		if len(match) < 8 {
			continue
		}

		age, err := strconv.Atoi(capture(cleanText, match, 1))
		if err != nil {
			continue
		}

		unit := capture(cleanText, match, 2)
		operator := capture(cleanText, match, 3)

		if unit == "대" {
			age = (age/10)*10 + 9
			if operator == "미만" {
				age -= 10
			}
		} else if operator == "미만" {
			age -= 1
		}

		if !isValidAge(age) {
			continue
		}

		candidates = append(candidates, ageCondition{
			max: &age,
			raw: strings.TrimSpace(cleanText[match[0]:match[1]]),
		})
	}

	if len(candidates) == 0 {
		return
	}

	// 단일 범위 후보 생성 시도 (1개의 min 조건 + 1개의 max 조건인 경우 통합)
	var mergedAge ageCondition
	if len(candidates) == 2 && ((candidates[0].min != nil && candidates[1].max != nil) || (candidates[0].max != nil && candidates[1].min != nil)) {
		if candidates[0].min != nil {
			mergedAge.min = candidates[0].min
			mergedAge.max = candidates[1].max
		} else {
			mergedAge.min = candidates[1].min
			mergedAge.max = candidates[0].max
		}
		candidates = []ageCondition{mergedAge}
	}

	first := candidates[0]

	for _, candidate := range candidates[1:] {
		if !sameAgeCondition(first, candidate) {
			condition.UnparsedConditions = append(
				condition.UnparsedConditions,
				p.ageContext(text, candidates),
			)
			return
		}
	}

	if first.min != nil && first.max != nil && *first.min > *first.max {
		condition.UnparsedConditions = append(
			condition.UnparsedConditions,
			p.ageContext(text, candidates),
		)
		return
	}

	if first.min != nil {
		condition.MinAge = first.min
	}

	if first.max != nil {
		condition.MaxAge = first.max
	}

}

func sameAgeCondition(a, b ageCondition) bool {
	if (a.min == nil) != (b.min == nil) {
		return false
	}

	if (a.max == nil) != (b.max == nil) {
		return false
	}

	if a.min != nil && *a.min != *b.min {
		return false
	}

	if a.max != nil && *a.max != *b.max {
		return false
	}

	return true
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

	unparsedKeywords := []string{
		"장애인", "등록 장애인", "등록장애인", "장애",
		"학생", "대학생", "고등학생", "중학생", "초등학생", "재학생", "휴학생",
		"북한이탈주민", "새터민",
		"소득", "중위소득", "차상위", "기초생활", "수급자", "저소득",
	}

	for _, kw := range unparsedKeywords {
		if strings.Contains(text, kw) {
			condition.UnparsedConditions = append(condition.UnparsedConditions, kw)
		}
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

func (p *Parser) ageContext(
	text string,
	candidates []ageCondition,
) string {
	if len(candidates) == 0 {
		return ""
	}

	sentences := regexp.MustCompile(`[.!?\r\n]+`).Split(text, -1)

	var contexts []string
	seen := make(map[string]struct{})

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)

		if sentence == "" {
			continue
		}

		for _, candidate := range candidates {
			if candidate.raw == "" {
				continue
			}

			if strings.Contains(sentence, candidate.raw) {
				if _, exists := seen[sentence]; exists {
					break
				}

				seen[sentence] = struct{}{}
				contexts = append(contexts, sentence)
				break
			}
		}
	}

	return strings.Join(contexts, "\n")
}

func capture(text string, match []int, group int) string {
	start := match[group*2]
	end := match[group*2+1]

	if start < 0 || end < 0 {
		return ""
	}

	return text[start:end]
}

