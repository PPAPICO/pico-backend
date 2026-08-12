package policycondition

import "github.com/janghanul090801/pico-backend/domain"

func Merge(
	dst *domain.PolicyCondition,
	src domain.PolicyCondition,
) {
	if dst.MinAge == nil {
		dst.MinAge = src.MinAge
	}

	if dst.MaxAge == nil {
		dst.MaxAge = src.MaxAge
	}

	if dst.Gender == nil {
		dst.Gender = src.Gender
	}

	if dst.RequireStudent == nil {
		dst.RequireStudent = src.RequireStudent
	}

	if dst.RequireYouth == nil {
		dst.RequireYouth = src.RequireYouth
	}

	if dst.RequirePregnant == nil {
		dst.RequirePregnant = src.RequirePregnant
	}

	if dst.RequireBusiness == nil {
		dst.RequireBusiness = src.RequireBusiness
	}

	if dst.RequireDisabled == nil {
		dst.RequireDisabled = src.RequireDisabled
	}

	if dst.RequireForeign == nil {
		dst.RequireForeign = src.RequireForeign
	}

	if len(src.Interests) > 0 {
		dst.Interests = append(
			dst.Interests,
			src.Interests...,
		)
	}
}
