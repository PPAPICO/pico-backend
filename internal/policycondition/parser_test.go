package policycondition_test

import (
	"testing"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/policycondition"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParseUnparsedConditions(t *testing.T) {
	parser := policycondition.NewParser()

	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "장애인 키워드 포함",
			text:     "등록 장애인을 위한 지원 사업",
			expected: []string{"장애인", "등록 장애인", "장애"},
		},
		{
			name:     "학생 키워드 포함",
			text:     "대학생 및 재학생 대상 청년 정책",
			expected: []string{"학생", "대학생", "재학생"},
		},
		{
			name:     "북한이탈주민 키워드 포함",
			text:     "북한이탈주민 지원 정책",
			expected: []string{"북한이탈주민"},
		},
		{
			name:     "소득 키워드 포함",
			text:     "중위소득 50% 이하 가구 대상",
			expected: []string{"소득", "중위소득"},
		},
		{
			name:     "복수 키워드 포함",
			text:     "저소득층 장애인 대학생 지원",
			expected: []string{"장애인", "장애", "학생", "대학생", "소득", "저소득"},
		},
		{
			name:     "해당 키워드 없음",
			text:     "만 19세 이상 청년 지원",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := parser.Parse(tt.text)
			assert.Equal(t, tt.expected, res.UnparsedConditions)
		})
	}
}

func TestParser_ParseAgeConditions(t *testing.T) {
	parser := policycondition.NewParser()

	intPtr := func(v int) *int { return &v }
	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name     string
		text     string
		expected domain.PolicyCondition
	}{
		{
			name: "범위 표기 (만 19세 ~ 34세)",
			text: "지원대상: 만 19세 ~ 34세 청년",
			expected: domain.PolicyCondition{
				MinAge:       intPtr(19),
				MaxAge:       intPtr(34),
				RequireYouth: boolPtr(true),
			},
		},
		{
			name: "범위 표기 (19~34세)",
			text: "19~34세 이하 청년 대상",
			expected: domain.PolicyCondition{
				MinAge:       intPtr(19),
				MaxAge:       intPtr(34),
				RequireYouth: boolPtr(true),
			},
		},
		{
			name: "연 나이 및 초과/미만 조건",
			text: "연 18세 초과 35세 미만 대상",
			expected: domain.PolicyCondition{
				MinAge: intPtr(19),
				MaxAge: intPtr(34),
			},
		},
		{
			name: "20대~30대 연령대 조건",
			text: "20대~30대 청년층 지원",
			expected: domain.PolicyCondition{
				MinAge:       intPtr(20),
				MaxAge:       intPtr(39),
				RequireYouth: boolPtr(true),
			},
		},
		{
			name: "날짜 및 시간 제외 후 나이 추출",
			text: "신청시간 09:00~18:00, 만 19세 이상",
			expected: domain.PolicyCondition{
				MinAge: intPtr(19),
			},
		},
		{
			name: "전화번호 무시 및 나이 조건 파싱",
			text: "문의: 010-1234-5678, 지원 대상: 만 19세 ~ 34세",
			expected: domain.PolicyCondition{
				MinAge: intPtr(19),
				MaxAge: intPtr(34),
			},
		},
		{
			name: "99세 초과 및 범위 밖 나이 무시",
			text: "문의전화 010-1234-5678 또는 150세 이상",
			expected: domain.PolicyCondition{},
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := parser.Parse(tt.text)
			assert.Equal(t, tt.expected.MinAge, res.MinAge)
			assert.Equal(t, tt.expected.MaxAge, res.MaxAge)
			if tt.expected.RequireYouth != nil {
				assert.Equal(t, tt.expected.RequireYouth, res.RequireYouth)
			}
		})
	}
}
