package usecase

import (
	"context"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
)

type favoriteUseCase struct {
	favoriteRepository domain.FavoriteRepository
	policyRepository   domain.PolicyRepository
	contextTimeout     time.Duration
}

func NewFavoriteUseCase(
	favoriteRepository domain.FavoriteRepository,
	policyRepository domain.PolicyRepository,
	timeout time.Duration,
) domain.FavoriteUseCase {
	return &favoriteUseCase{
		favoriteRepository: favoriteRepository,
		policyRepository:   policyRepository,
		contextTimeout:     timeout,
	}
}

func (u *favoriteUseCase) AddFavorite(c context.Context, userID *domain.ID, policyID *domain.ID) (*domain.Favorite, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	_, err := u.policyRepository.FindByID(ctx, policyID)
	if err != nil {
		return nil, domain.NewNotFoundError(err)
	}

	existing, _ := u.favoriteRepository.FindByUserIDAndPolicyID(ctx, userID, policyID)
	if existing != nil {
		return existing, nil
	}

	fav, err := u.favoriteRepository.Create(ctx, userID, policyID)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return fav, nil
}

func (u *favoriteUseCase) RemoveFavorite(c context.Context, userID *domain.ID, policyID *domain.ID) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err := u.favoriteRepository.Delete(ctx, userID, policyID)
	if err != nil {
		return domain.NewInternalServerError(err)
	}

	return nil
}

func (u *favoriteUseCase) ListFavoritesByUserID(c context.Context, userID *domain.ID) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	favs, err := u.favoriteRepository.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	policies := make([]*domain.Policy, 0, len(favs))
	for _, f := range favs {
		p, err := u.policyRepository.FindByID(ctx, &f.PolicyID)
		if err == nil && p != nil {
			policies = append(policies, p)
		}
	}

	return policies, nil
}
