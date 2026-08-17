package domain

import (
	"context"
	"time"
)

type Favorite struct {
	ID        ID        `json:"id"`
	UserID    ID        `json:"user_id"`
	PolicyID  ID        `json:"policy_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FavoriteRepository interface {
	Create(c context.Context, userID *ID, policyID *ID) (*Favorite, error)
	Delete(c context.Context, userID *ID, policyID *ID) error
	FindByUserIDAndPolicyID(c context.Context, userID *ID, policyID *ID) (*Favorite, error)
	FindAllByUserID(c context.Context, userID *ID) ([]*Favorite, error)
	FindAll(c context.Context) ([]*Favorite, error)
}

type FavoriteUseCase interface {
	AddFavorite(c context.Context, userID *ID, policyID *ID) (*Favorite, error)
	RemoveFavorite(c context.Context, userID *ID, policyID *ID) error
	ListFavoritesByUserID(c context.Context, userID *ID) ([]*Policy, error)
}
