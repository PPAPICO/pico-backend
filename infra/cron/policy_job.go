package cron

import (
	"context"
	"log"

	"github.com/janghanul090801/pico-backend/domain"
)

type PolicyJob struct {
	policyUseCase domain.PolicyUseCase
}

func NewPolicyJob(policyUseCase domain.PolicyUseCase) *PolicyJob {
	return &PolicyJob{
		policyUseCase: policyUseCase,
	}
}

func (j *PolicyJob) Run() {
	if policies, err := j.policyUseCase.GetFromApi(context.Background()); err != nil {
		log.Printf("policy job failed: %v", err)
	} else {
		log.Printf("policy job success: %v", len(policies))
	}
}
