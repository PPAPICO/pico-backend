package repository

import (
	"context"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/governmentpolicy"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/policymatch"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/user"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/collections"
)

type policyMatchRepository struct {
	client *ent.Client
}

func NewPolicyMatchRepository(client *ent.Client) domain.PolicyMatchRepository {
	return &policyMatchRepository{client: client}
}

func (r *policyMatchRepository) FindAllByUserID(c context.Context, userID *domain.ID) ([]*domain.PolicyMatch, error) {
	pms, err := r.client.PolicyMatch.Query().Where(
		policymatch.HasUserWith(
			user.IDEQ(*userID),
		),
	).All(c)
	if err != nil {
		return nil, err
	}
	return collections.Map(pms, toDomainPolicyMatch), nil
}

func (r *policyMatchRepository) FindAllByUserIDAndStatus(c context.Context, userID *domain.ID, status domain.Match) ([]*domain.PolicyMatch, error) {
	now := time.Now()

	pms, err := r.client.PolicyMatch.Query().Where(
		policymatch.HasUserWith(
			user.IDEQ(*userID),
		),
		policymatch.HasPolicyWith(
			governmentpolicy.StartDateLTE(now),
			governmentpolicy.EndDateGTE(now),
		),
	).All(c)

	if err != nil {
		return nil, err
	}
	return collections.Map(pms, toDomainPolicyMatch), nil
}

func (r *policyMatchRepository) Create(c context.Context, policyMatch *domain.PolicyMatch) (*domain.PolicyMatch, error) {
	p, err := r.client.PolicyMatch.Create().
		SetPolicyID(policyMatch.PolicyID).
		SetUserID(policyMatch.UserID).
		SetMatch(domainMatchToEntMatch(policyMatch.Status)).Save(c)
	if err != nil {
		return nil, err
	}
	return toDomainPolicyMatch(p), nil
}

func (r *policyMatchRepository) Update(c context.Context, id *domain.ID, status domain.Match) (*domain.PolicyMatch, error) {
	p, err := r.client.PolicyMatch.UpdateOneID(*id).SetMatch(domainMatchToEntMatch(status)).Save(c)
	if err != nil {
		return nil, err
	}
	return toDomainPolicyMatch(p), nil
}
