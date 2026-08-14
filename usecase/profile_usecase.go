package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"golang.org/x/crypto/bcrypt"
)

type profileUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewProfileUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.ProfileUseCase {
	return &profileUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (u *profileUseCase) GetByID(c context.Context, userID *domain.ID) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewNotFoundError(err)
		}
		return nil, domain.NewInternalServerError(err)
	}

	return &domain.Profile{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Age:        user.Age,
		RegionCode: user.RegionCode,
		Gender:     user.Gender,
		IsStudent:  user.IsStudent,
		IsYouth:    user.IsYouth,
		IsPregnant: user.IsPregnant,
		IsBusiness: user.IsBusiness,
		IsForeign:  user.IsForeign,
		Interests:  user.Interests,
		CreatedAt:  user.CreatedAt,
	}, nil
}

func (u *profileUseCase) Update(c context.Context, ID *domain.ID, params *domain.ProfileUpdateRequest) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FindByID(ctx, ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewNotFoundError(err)
		}
		return nil, domain.NewInternalServerError(err)
	}

	encrypted, err := bcrypt.GenerateFromPassword(
		[]byte(params.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, domain.NewBadRequestError(err)
	}

	if string(encrypted) != user.Password {
		return nil, domain.NewBadRequestError(errors.New("invalid password"))
	}

	user, err = u.userRepository.Update(ctx, &domain.User{
		ID:         *ID,
		Name:       user.Name,
		Age:        params.Age,
		RegionCode: params.RegionCode,
		Gender:     params.Gender,
		IsStudent:  params.IsStudent,
		IsYouth:    params.IsYouth,
		IsPregnant: params.IsPregnant,
		IsBusiness: params.IsBusiness,
		IsForeign:  params.IsForeign,
		Interests:  params.Interests,
	})
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return &domain.Profile{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Age:        user.Age,
		RegionCode: user.RegionCode,
		Gender:     user.Gender,
		IsStudent:  user.IsStudent,
		IsYouth:    user.IsYouth,
		IsPregnant: user.IsPregnant,
		IsBusiness: user.IsBusiness,
		IsForeign:  user.IsForeign,
		Interests:  user.Interests,
		CreatedAt:  user.CreatedAt,
	}, nil
}
