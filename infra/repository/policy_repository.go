package repository

import (
	"context"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/governmentpolicy"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/collections"
)

type policyRepository struct {
	client *ent.Client
}

func NewPolicyRepository(client *ent.Client) domain.PolicyRepository {
	return &policyRepository{client: client}
}

func (r *policyRepository) FindAll(c context.Context) ([]*domain.Policy, error) {
	policies, err := r.client.GovernmentPolicy.Query().All(c)
	if err != nil {
		return nil, err
	}

	return collections.Map(policies, toDomainPolicy), nil
}

func (r *policyRepository) FindByID(c context.Context, id *domain.ID) (*domain.Policy, error) {
	p, err := r.client.GovernmentPolicy.Get(c, *id)
	if err != nil {
		return nil, err
	}
	return toDomainPolicy(p), nil
}

func (r *policyRepository) FindAllByRegionCode(c context.Context, regionCode int) ([]*domain.Policy, error) {
	policies, err := r.client.GovernmentPolicy.Query().Where(
		governmentpolicy.RegionCodeEQ(regionCode),
	).All(c)
	if err != nil {
		return nil, err
	}
	return collections.Map(policies, toDomainPolicy), nil
}

func (r *policyRepository) FindAllByRegionCodeAndActive(c context.Context, regionCode int) ([]*domain.Policy, error) {
	now := time.Now()

	policies, err := r.client.GovernmentPolicy.Query().Where(
		governmentpolicy.RegionCodeEQ(regionCode),
		governmentpolicy.StartDateLTE(now),
		governmentpolicy.EndDateGTE(now),
	).All(c)
	if err != nil {
		return nil, err
	}

	return collections.Map(policies, toDomainPolicy), nil
}

func (r *policyRepository) Create(c context.Context, policy *domain.Policy) (*domain.Policy, error) {
	p, err := r.client.GovernmentPolicy.Create().
		SetTitle(policy.Title).
		SetDescription(policy.Description).
		SetRegionCode(policy.RegionCode).
		SetStartDate(policy.StartDate).
		SetEndDate(policy.EndDate).
		SetAddress(policy.Address).
		SetLatitude(policy.Latitude).
		SetLongitude(policy.Longitude).
		Save(c)

	if err != nil {
		return nil, err
	}
	return toDomainPolicy(p), nil
}

func (r *policyRepository) Delete(c context.Context, id *domain.ID) error {
	return r.client.GovernmentPolicy.DeleteOneID(*id).Exec(c)
}
