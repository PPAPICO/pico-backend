package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type PolicyJob struct {
	policyUseCase       domain.PolicyUseCase
	userRepository      domain.UserRepository
	notificationUseCase domain.NotificationUseCase
	favoriteUseCase     domain.FavoriteUseCase
}

func NewPolicyJob(
	policyUseCase domain.PolicyUseCase,
	userRepository domain.UserRepository,
	notificationUseCase domain.NotificationUseCase,
	favoriteUseCase domain.FavoriteUseCase,
) *PolicyJob {
	return &PolicyJob{
		policyUseCase:       policyUseCase,
		userRepository:      userRepository,
		notificationUseCase: notificationUseCase,
		favoriteUseCase:     favoriteUseCase,
	}
}

func (j *PolicyJob) Run() {
	ctx := context.Background()
	policies, err := j.policyUseCase.GetFromApi(ctx)
	if err != nil {
		log.Printf("policy job failed: %v", err)
	} else {
		log.Printf("policy job success: %v", len(policies))
		users, err := j.userRepository.FindAll(ctx)
		if err != nil {
			log.Printf("user repository failed: %v", err)
		} else {
			collections.ForEach(users, func(user *domain.User) {
				if matches, err := j.policyUseCase.SavePolicyMatches(ctx, user, policies); err != nil {
					log.Printf("save policy match failed: %v", err)
				} else {
					log.Printf("user: %v : save policy match success: %v", user.ID, len(matches))
				}

				if j.notificationUseCase != nil && len(user.Interests) > 0 {
					userInterestSet := make(map[domain.Interest]bool)
					for _, interest := range user.Interests {
						userInterestSet[interest] = true
					}

					for _, policy := range policies {
						for _, pInterest := range policy.Condition.Interests {
							if userInterestSet[pInterest] {
								msg := fmt.Sprintf("관심 분야(%s)의 새 정책 '%s'이(가) 등록되었습니다.", pInterest, policy.Title)
								_, _ = j.notificationUseCase.Create(ctx, &domain.Notification{
									ReceiverID: user.ID,
									Type:       domain.NotificationTypeINTEREST,
									Message:    msg,
									Metadata: map[string]any{
										"policy_id": policy.ID,
										"interest":  pInterest,
									},
								})
								break
							}
						}
					}
				}
			})
		}
	}

	j.checkFavoriteDeadlines(ctx)
}

func (j *PolicyJob) checkFavoriteDeadlines(ctx context.Context) {
	if j.favoriteUseCase == nil || j.notificationUseCase == nil {
		return
	}

	users, err := j.userRepository.FindAll(ctx)
	if err != nil {
		return
	}

	now := time.Now()
	threeDaysLater := now.Add(3 * 24 * time.Hour)

	for _, user := range users {
		favPolicies, err := j.favoriteUseCase.ListFavoritesByUserID(ctx, &user.ID)
		if err != nil {
			continue
		}

		for _, policy := range favPolicies {
			if !policy.EndDate.IsZero() && policy.EndDate.After(now) && policy.EndDate.Before(threeDaysLater) {
				daysLeft := int(policy.EndDate.Sub(now).Hours() / 24)
				if daysLeft < 0 {
					daysLeft = 0
				}
				msg := fmt.Sprintf("즐겨찾기한 정책 '%s'의 마감기한이 %d일 남았습니다.", policy.Title, daysLeft)
				_, _ = j.notificationUseCase.Create(ctx, &domain.Notification{
					ReceiverID: user.ID,
					Type:       domain.NotificationTypeSTAR,
					Message:    msg,
					Metadata: map[string]any{
						"policy_id": policy.ID,
						"days_left": daysLeft,
					},
				})
			}
		}
	}
}
