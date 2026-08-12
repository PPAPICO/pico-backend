package usecase

import "github.com/janghanul090801/pico-backend/domain"

type PolicyMatcher struct{}

func NewPolicyMatcher() *PolicyMatcher {
	return &PolicyMatcher{}
}

func (m *PolicyMatcher) Match(
	user *domain.User,
	policy *domain.Policy,
) domain.Match {
	results := []domain.ConditionResult{
		m.checkAge(user, &policy.Condition),
		m.checkRegion(user, &policy.Condition),
		m.checkGender(user, &policy.Condition),
		m.checkStudent(user, &policy.Condition),
		m.checkYouth(user, &policy.Condition),
		m.checkPregnant(user, &policy.Condition),
		m.checkBusiness(user, &policy.Condition),
		m.checkDisabled(user, &policy.Condition),
		m.checkForeign(user, &policy.Condition),
	}

	hasUnknown := len(policy.Condition.UnparsedConditions) > 0

	for _, result := range results {
		switch result {
		case domain.ConditionFail:
			return domain.MatchIMPOSSIBLE

		case domain.ConditionUnknown:
			hasUnknown = true
		}
	}

	if hasUnknown {
		return domain.MatchUNCERTAIN
	}

	return domain.MatchPOSSIBLE
}

func (m *PolicyMatcher) checkAge(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.MinAge == nil && condition.MaxAge == nil {
		return domain.ConditionPass
	}

	if condition.MinAge != nil && user.Age < *condition.MinAge {
		return domain.ConditionFail
	}

	if condition.MaxAge != nil && user.Age > *condition.MaxAge {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}

func (m *PolicyMatcher) checkRegion(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if len(condition.RegionCodes) == 0 {
		return domain.ConditionPass
	}

	for _, code := range condition.RegionCodes {
		if user.RegionCode == code {
			return domain.ConditionPass
		}
	}

	return domain.ConditionFail
}

func (m *PolicyMatcher) checkGender(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.Gender == nil {
		return domain.ConditionPass
	}

	if user.Gender != *condition.Gender {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}

func (m *PolicyMatcher) checkStudent(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.RequireStudent == nil {
		return domain.ConditionPass
	}

	if user.IsStudent != *condition.RequireStudent {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}

func (m *PolicyMatcher) checkYouth(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.RequireYouth == nil {
		return domain.ConditionPass
	}

	if user.IsYouth != *condition.RequireYouth {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}

func (m *PolicyMatcher) checkPregnant(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.RequirePregnant == nil {
		return domain.ConditionPass
	}

	if user.IsPregnant != *condition.RequirePregnant {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}
func (m *PolicyMatcher) checkBusiness(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.RequireBusiness == nil {
		return domain.ConditionPass
	}

	if user.IsBusiness != *condition.RequireBusiness {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}
func (m *PolicyMatcher) checkDisabled(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.RequireDisabled == nil {
		return domain.ConditionPass
	}

	if user.IsDisabled != *condition.RequireDisabled {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}
func (m *PolicyMatcher) checkForeign(
	user *domain.User,
	condition *domain.PolicyCondition,
) domain.ConditionResult {

	if condition.RequireForeign == nil {
		return domain.ConditionPass
	}

	if user.IsForeign != *condition.RequireForeign {
		return domain.ConditionFail
	}

	return domain.ConditionPass
}
