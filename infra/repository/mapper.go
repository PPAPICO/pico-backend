package repository

import (
	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"github.com/janghanul090801/pico-backend/ent/policymatch"
)

func toDomainUser(u *ent.User) *domain.User {
	return &domain.User{
		ID:          domain.ID(u.ID),
		Name:        u.Name,
		Email:       u.Email,
		Password:    u.Password,
		Age:         u.Age,
		RegionCode:  u.RegionCode,
		Gender:      domain.Gender(u.Gender),
		IsStudent:   u.IsStudent,
		IsYouth:     u.IsYouth,
		IsPregnant:  u.IsPregnant,
		IsBusiness:  u.IsBusiness,
		Interests:   stringsToInterests(u.Interests),
		IsDisabled:  u.IsDisabled,
		Nationality: domain.Nationality(u.Nationality),
		CreatedAt:   u.CreatedAt,
	}
}

func toDomainPolicy(p *ent.GovernmentPolicy) *domain.Policy {
	return &domain.Policy{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		RegionCode:  p.RegionCode,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		Address:     p.Address,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
	}
}

func toDomainPolicyMatch(pm *ent.PolicyMatch) *domain.PolicyMatch {
	return &domain.PolicyMatch{
		ID:       pm.ID,
		PolicyID: pm.Edges.Policy.ID,
		UserID:   pm.Edges.User.ID,
		Status:   entMatchToDomainMatch(pm.Match),
	}
}

func interestsToStrings(interests []domain.Interest) []string {
	result := make([]string, len(interests))
	for i, v := range interests {
		result[i] = string(v)
	}
	return result
}

func stringsToInterests(values []string) []domain.Interest {
	result := make([]domain.Interest, len(values))
	for i, v := range values {
		result[i] = domain.Interest(v)
	}
	return result
}

func entMatchToDomainMatch(value policymatch.Match) domain.Match {
	switch value {
	case policymatch.MatchPOSSIBLE:
		return domain.MatchIMPOSSIBLE
	case policymatch.MatchUNCERTAIN:
		return domain.MatchUNCERTAIN
	default:
		return domain.MatchIMPOSSIBLE
	}
}

func domainMatchToEntMatch(value domain.Match) policymatch.Match {
	switch value {
	case domain.MatchIMPOSSIBLE:
		return policymatch.MatchIMPOSSIBLE
	case domain.MatchUNCERTAIN:
		return policymatch.MatchUNCERTAIN
	default:
		return policymatch.MatchPOSSIBLE
	}
}
