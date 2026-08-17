package repository

import (
	"context"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"github.com/janghanul090801/pico-backend/ent/favorite"
	"github.com/janghanul090801/pico-backend/ent/governmentpolicy"
	"github.com/janghanul090801/pico-backend/ent/user"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type favoriteRepository struct {
	client *ent.Client
}

func NewFavoriteRepository(client *ent.Client) domain.FavoriteRepository {
	return &favoriteRepository{client: client}
}

func (r *favoriteRepository) Create(c context.Context, userID *domain.ID, policyID *domain.ID) (*domain.Favorite, error) {
	fav, err := r.client.Favorite.Create().
		SetUserID(*userID).
		SetPolicyID(*policyID).
		Save(c)
	if err != nil {
		return nil, err
	}
	fav, err = r.client.Favorite.Query().
		Where(favorite.IDEQ(fav.ID)).
		WithUser().
		WithPolicy().
		Only(c)
	if err != nil {
		return nil, err
	}
	return toDomainFavorite(fav), nil
}

func (r *favoriteRepository) Delete(c context.Context, userID *domain.ID, policyID *domain.ID) error {
	_, err := r.client.Favorite.Delete().
		Where(
			favorite.HasUserWith(user.IDEQ(*userID)),
			favorite.HasPolicyWith(governmentpolicy.IDEQ(*policyID)),
		).
		Exec(c)
	return err
}

func (r *favoriteRepository) FindByUserIDAndPolicyID(c context.Context, userID *domain.ID, policyID *domain.ID) (*domain.Favorite, error) {
	fav, err := r.client.Favorite.Query().
		Where(
			favorite.HasUserWith(user.IDEQ(*userID)),
			favorite.HasPolicyWith(governmentpolicy.IDEQ(*policyID)),
		).
		WithUser().
		WithPolicy().
		Only(c)
	if err != nil {
		return nil, err
	}
	return toDomainFavorite(fav), nil
}

func (r *favoriteRepository) FindAllByUserID(c context.Context, userID *domain.ID) ([]*domain.Favorite, error) {
	favs, err := r.client.Favorite.Query().
		Where(
			favorite.HasUserWith(user.IDEQ(*userID)),
		).
		WithUser().
		WithPolicy().
		All(c)
	if err != nil {
		return nil, err
	}
	return collections.Map(favs, toDomainFavorite), nil
}

func (r *favoriteRepository) FindAll(c context.Context) ([]*domain.Favorite, error) {
	favs, err := r.client.Favorite.Query().
		WithUser().
		WithPolicy().
		All(c)
	if err != nil {
		return nil, err
	}
	return collections.Map(favs, toDomainFavorite), nil
}
