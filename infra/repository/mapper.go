package repository

import (
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent"
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
