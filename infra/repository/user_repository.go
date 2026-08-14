package repository

import (
	"context"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"github.com/janghanul090801/pico-backend/ent/user"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) domain.UserRepository {
	return &userRepository{
		client: client,
	}
}

func (r *userRepository) Create(c context.Context, u *domain.User) (*domain.User, error) {
	userEnt, err := r.client.User.Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetAge(u.Age).
		SetRegionCode(u.RegionCode).
		SetGender(user.Gender(u.Gender)).
		SetIsStudent(u.IsStudent).
		SetIsYouth(u.IsYouth).
		SetInterests(interestsToStrings(u.Interests)).
		SetIsDisabled(u.IsDisabled).
		SetIsForeign(u.IsForeign).
		SetIsPregnant(u.IsPregnant).
		SetIsBusiness(u.IsBusiness).Save(c)
	if err != nil {
		return nil, err
	}

	return toDomainUser(userEnt), nil
}

func (r *userRepository) Update(c context.Context, u *domain.User) (*domain.User, error) {
	user, err := r.client.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetAge(u.Age).
		SetRegionCode(u.RegionCode).
		SetGender(user.Gender(u.Gender)).
		SetIsStudent(u.IsStudent).
		SetIsYouth(u.IsYouth).
		SetInterests(interestsToStrings(u.Interests)).
		SetIsDisabled(u.IsDisabled).
		SetIsForeign(u.IsForeign).
		SetIsPregnant(u.IsPregnant).
		SetIsBusiness(u.IsBusiness).Save(c)
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

func (r *userRepository) FindAll(c context.Context) ([]*domain.User, error) {
	users, err := r.client.User.Query().All(c)
	if err != nil {
		return nil, err
	}

	return collections.Map(users, toDomainUser), nil
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
