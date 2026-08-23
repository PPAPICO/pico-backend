package domain

type PolicyCondition struct {
	MinAge *int `json:"min_age,omitempty"`
	MaxAge *int `json:"max_age,omitempty"`

	RegionCodes []int `json:"region_codes,omitempty"`

	Gender *Gender `json:"gender,omitempty"`

	RequireStudent  *bool `json:"require_student,omitempty"`
	RequireYouth    *bool `json:"require_youth,omitempty"`
	RequirePregnant *bool `json:"require_pregnant,omitempty"`
	RequireBusiness *bool `json:"require_business,omitempty"`
	RequireDisabled *bool `json:"require_disabled,omitempty"`
	RequireForeign  *bool `json:"require_foreign,omitempty"`

	Interests []Interest `json:"interests,omitempty"`

	UnparsedConditions []string `json:"unparsed_conditions,omitempty"`

	// 외부 API(온통청년/복지로)가 제공하는 신청/원문 페이지 URL.
	// ent 스키마 변경 없이 JSON 컬럼(condition)에 함께 저장하기 위해 여기 추가함.
	SourceURL string `json:"source_url,omitempty"`
}

type ConditionResult int

const (
	ConditionPass ConditionResult = iota
	ConditionFail
	ConditionUnknown
)
