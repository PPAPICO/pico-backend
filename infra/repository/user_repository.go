package repository

import (
	"context"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/user"
)

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) domain.UserRepository {
	return &userRepository{
		client: client,
	}
}

func (r *userRepository) Save(c context.Context, u *domain.User) (*domain.User, error) {
	builder := r.client.User.Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetAge(u.Age).
		SetRegion(u.Region).
		SetGender(user.Gender(u.Gender)).
		SetIsStudent(u.IsStudent).
		SetIsYouth(u.IsYouth).
		SetInterests(interestsToStrings(u.Interests)).
		SetIsDisabled(u.IsDisabled).
		SetNationality(user.Nationality(u.Nationality))

	if u.IsPregnant != nil {
		builder.SetIsPregnant(*u.IsPregnant)
	}

	if u.IsBusiness != nil {
		builder.SetIsBusiness(*u.IsBusiness)
	}

	userEnt, err := builder.Save(c)
	if err != nil {
		return nil, err
	}

	return toDomainUser(userEnt), nil
}

func (r *userRepository) FindAll(c context.Context) ([]*domain.User, error) {
	users, err := r.client.User.Query().All(c)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = toDomainUser(u)
	}

	return result, nil
}

func (r *userRepository) FindByEmail(c context.Context, email string) (*domain.User, error) {
	u, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(c)
	if err != nil {
		return nil, err
	}

	return toDomainUser(u), nil
}

func (r *userRepository) FindByID(c context.Context, id *domain.ID) (*domain.User, error) {
	u, err := r.client.User.Get(c, *id)
	if err != nil {
		return nil, err
	}

	return toDomainUser(u), nil
}
