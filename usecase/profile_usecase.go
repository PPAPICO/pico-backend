package usecase

import (
	"context"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
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

func (pu *profileUseCase) GetProfileByID(c context.Context, userID *domain.ID) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	user, err := pu.userRepository.FindByID(ctx, userID)
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
