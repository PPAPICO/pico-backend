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
	POSSIBLE   Match = iota
	UNCERTAIN  Match = iota
	IMPOSSIBLE Match = iota
)

type PolicyRepository interface {
	FindAll(c context.Context) ([]*Policy, error)
	FindByID(c context.Context, id *ID) (*Policy, error)
	FindAllByRegionCode(c context.Context, regionCode int) ([]*Policy, error)
	FindAllByRegionCodeAndActive(c context.Context, regionCode int) ([]*Policy, error)
	Create(c context.Context, policy *Policy) (*Policy, error)
	Delete(c context.Context, id *ID) error
}

type PolicyMatchRepository interface {
	FindAllByUserID(c context.Context, userID *ID) ([]*PolicyMatch, error)
	FindAllByUserIDAndStatus(c context.Context, userID *ID, status Match) ([]*PolicyMatch, error)
	Create(c context.Context, policyMatch *PolicyMatch) (*PolicyMatch, error)
	Update(c context.Context, id *ID, status Match) (*PolicyMatch, error)
}

type PolicyUseCase interface {
	GetByID(c context.Context, id *ID) (*Policy, error)
	ListRegionCodeAndActive(c context.Context, regionCode int) ([]*Policy, error)
	GetFromApi(c context.Context) ([]*Policy, error)
	GetMatchesByUserID(c context.Context, userID *ID) ([]*Policy, error)
}
