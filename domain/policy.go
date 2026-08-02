package domain

import (
	"context"
	"time"
)

type Policy struct {
	ID          ID
	Title       string
	Description string
	RegionCode  int
	StartDate   time.Time
	EndDate     time.Time
	Address     string
	Latitude    float64
	Longitude   float64
}

type PolicyMatch struct {
	ID       ID
	PolicyID ID
	UserID   ID
	Status   Match
}

type Match int

const (
	Possible   Match = iota
	Uncertain  Match = iota
	Impossible Match = iota
)

type PolicyRepository interface {
	FindAll(ctx context.Context) ([]*Policy, error)
	FindByID(ctx context.Context, id *ID) (*Policy, error)
	FindAllByRegionCode(ctx context.Context, regionCode int) ([]*Policy, error)
	FindAllByRegionCodeAndActive(ctx context.Context, regionCode int) ([]*Policy, error)
	Save(ctx context.Context, policy *Policy) (*Policy, error)
	Delete(ctx context.Context, id *ID) error
}

type PolicyUseCase interface {
	GetByID(ctx context.Context, id *ID) (*Policy, error)
	ListRegionCodeAndActive(ctx context.Context, regionCode int) ([]*Policy, error)
	GetFromApi(ctx context.Context) ([]*Policy, error)
	GetMatchesByUserID(ctx context.Context, userID *ID) ([]*Policy, error)
}
