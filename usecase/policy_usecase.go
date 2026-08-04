package usecase

import (
	"context"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

type policyUseCase struct {
	policyRepository      domain.PolicyRepository
	policyMatchRepository domain.PolicyMatchRepository
	contextTimeout        time.Duration
}

func NewPolicyUseCase(policyRepository domain.PolicyRepository, policyMatchRepository domain.PolicyMatchRepository, contextTimeout time.Duration) domain.PolicyUseCase {
	return &policyUseCase{
		policyRepository:      policyRepository,
		policyMatchRepository: policyMatchRepository,
		contextTimeout:        contextTimeout,
	}
}

func (u *policyUseCase) GetByID(c context.Context, id *domain.ID) (*domain.Policy, error) {
	//TODO implement me
	panic("implement me")
}

func (u *policyUseCase) ListRegionCodeAndActive(c context.Context, regionCode int) ([]*domain.Policy, error) {
	//TODO implement me
	panic("implement me")
}

func (u *policyUseCase) GetFromApi(c context.Context) ([]*domain.Policy, error) {
	//TODO implement me
	panic("implement me")
}

func (u *policyUseCase) GetMatchesByUserID(c context.Context, userID *domain.ID) ([]*domain.Policy, error) {
	//TODO implement me
	panic("implement me")
}
