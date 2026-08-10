package repository_test

import (
	"context"
	"testing"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent/enttest"
	"github.com/janghanul090801/pico-backend/infra/repository"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)

	u, err := repo.Create(context.Background(), &domain.User{
		Name:       "hanul",
		Email:      "hanul@gmail.com",
		Password:   "123456",
		Age:        20,
		RegionCode: 110,
		Gender:     domain.GenderMale,
		IsStudent:  true,
		IsYouth:    true,
		IsPregnant: false,
		IsBusiness: true,
		Interests: []domain.Interest{
			domain.InterestCulture,
			domain.InterestEducation,
		},
		IsDisabled: false,
		IsForeign:  false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)

	assert.Equal(t, "hanul", u.Name)
	assert.Equal(t, "hanul@gmail.com", u.Email)
	assert.Equal(t, 20, u.Age)
	assert.Equal(t, 110, u.RegionCode)
	assert.Equal(t, domain.GenderMale, u.Gender)
	assert.True(t, u.IsStudent)
	assert.True(t, u.IsYouth)
	assert.False(t, u.IsPregnant)
	assert.True(t, u.IsBusiness)
	assert.True(t, u.IsBusiness)
	assert.Equal(t, []domain.Interest{
		domain.InterestCulture,
		domain.InterestEducation,
	}, u.Interests)
	assert.False(t, u.IsDisabled)
	assert.False(t, u.IsForeign)

	assert.NotZero(t, u.ID)
	assert.False(t, u.CreatedAt.IsZero())
}
