package repository

import (
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/policymatch"
)

func toDomainUser(u *ent.User) *domain.User {
	return &domain.User{
		ID:          domain.ID(u.ID),
		Name:        u.Name,
		Email:       u.Email,
		Password:    u.Password,
		Age:         u.Age,
		Region:      u.Region,
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
		return domain.IMPOSSIBLE
	case policymatch.MatchUNCERTAIN:
		return domain.UNCERTAIN
	default:
		return domain.IMPOSSIBLE
	}
}

func domainMatchToEntMatch(value domain.Match) policymatch.Match {
	switch value {
	case domain.IMPOSSIBLE:
		return policymatch.MatchIMPOSSIBLE
	case domain.UNCERTAIN:
		return policymatch.MatchUNCERTAIN
	default:
		return policymatch.MatchPOSSIBLE
	}
}
