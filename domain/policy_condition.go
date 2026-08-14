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
}

type ConditionResult int

const (
	ConditionPass ConditionResult = iota
	ConditionFail
	ConditionUnknown
)
