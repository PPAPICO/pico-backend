package repository_test

import (
	"context"
	"testing"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/enttest"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/infra/repository"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)

	isBusiness := true

	u, err := repo.Save(context.Background(), &domain.User{
		Name:       "hanul",
		Email:      "hanul@gmail.com",
		Password:   "123456",
		Age:        20,
		Region:     "서울특별시",
		Gender:     domain.GenderMale,
		IsStudent:  true,
		IsYouth:    true,
		IsPregnant: nil,
		IsBusiness: &isBusiness,
		Interests: []domain.Interest{
			domain.InterestCulture,
			domain.InterestEducation,
		},
		IsDisabled:  false,
		Nationality: domain.NationalityDomestic,
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)

	assert.Equal(t, "hanul", u.Name)
	assert.Equal(t, "hanul@gmail.com", u.Email)
	assert.Equal(t, 20, u.Age)
	assert.Equal(t, "서울특별시", u.Region)
	assert.Equal(t, domain.GenderMale, u.Gender)
	assert.True(t, u.IsStudent)
	assert.True(t, u.IsYouth)
	assert.Nil(t, u.IsPregnant)
	assert.NotNil(t, u.IsBusiness)
	assert.True(t, *u.IsBusiness)
	assert.Equal(t, []domain.Interest{
		domain.InterestCulture,
		domain.InterestEducation,
	}, u.Interests)
	assert.False(t, u.IsDisabled)
	assert.Equal(t, domain.NationalityDomestic, u.Nationality)

	assert.NotZero(t, u.ID)
	assert.False(t, u.CreatedAt.IsZero())
}
