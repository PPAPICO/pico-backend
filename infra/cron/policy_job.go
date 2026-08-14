package cron

import (
	"context"
	"log"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type PolicyJob struct {
	policyUseCase  domain.PolicyUseCase
	userRepository domain.UserRepository
}

func NewPolicyJob(policyUseCase domain.PolicyUseCase, userRepository domain.UserRepository) *PolicyJob {
	return &PolicyJob{
		policyUseCase:  policyUseCase,
		userRepository: userRepository,
	}
}

func (j *PolicyJob) Run() {
	ctx := context.Background()
	if policies, err := j.policyUseCase.GetFromApi(ctx); err != nil {
		log.Printf("policy job failed: %v", err)
	} else {
		log.Printf("policy job success: %v", len(policies))
		users, err := j.userRepository.FindAll(ctx)
		if err != nil {
			log.Printf("user repository failed: %v", err)
			return
		}
		collections.ForEach(users, func(user *domain.User) {
			if matches, err := j.policyUseCase.SavePolicyMatches(ctx, user, policies); err != nil {
				log.Printf("save policy match failed: %v", err)
			} else {
				log.Printf("user: %v : save policy match success: %v", user.ID, len(matches))
			}
		})
	}
}
