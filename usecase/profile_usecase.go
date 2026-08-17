package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"golang.org/x/crypto/bcrypt"
)

type profileUseCase struct {
	userRepository      domain.UserRepository
	policyUseCase       domain.PolicyUseCase
	notificationUseCase domain.NotificationUseCase
	contextTimeout      time.Duration
}

func NewProfileUseCase(
	userRepository domain.UserRepository,
	policyUseCase domain.PolicyUseCase,
	notificationUseCase domain.NotificationUseCase,
	timeout time.Duration,
) domain.ProfileUseCase {
	return &profileUseCase{
		userRepository:      userRepository,
		policyUseCase:       policyUseCase,
		notificationUseCase: notificationUseCase,
		contextTimeout:      timeout,
	}
}

func (u *profileUseCase) GetByID(c context.Context, userID *domain.ID) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewNotFoundError(err)
		}
		return nil, domain.NewInternalServerError(err)
	}

	return &domain.Profile{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Age:        user.Age,
		RegionCode: user.RegionCode,
		Gender:     user.Gender,
		IsStudent:  user.IsStudent,
		IsYouth:    user.IsYouth,
		IsPregnant: user.IsPregnant,
		IsBusiness: user.IsBusiness,
		IsForeign:  user.IsForeign,
		Interests:  user.Interests,
		CreatedAt:  user.CreatedAt,
	}, nil
}

func (u *profileUseCase) Update(c context.Context, ID *domain.ID, params *domain.ProfileUpdateRequest) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FindByID(ctx, ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewNotFoundError(err)
		}
		return nil, domain.NewInternalServerError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return nil, domain.NewBadRequestError(errors.New("invalid password"))
	}

	user, err = u.userRepository.Update(ctx, &domain.User{
		ID:         *ID,
		Name:       user.Name,
		Age:        params.Age,
		RegionCode: params.RegionCode,
		Gender:     params.Gender,
		IsStudent:  params.IsStudent,
		IsYouth:    params.IsYouth,
		IsPregnant: params.IsPregnant,
		IsBusiness: params.IsBusiness,
		IsForeign:  params.IsForeign,
		Interests:  params.Interests,
	})
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	if u.policyUseCase != nil {
		oldMatches, _ := u.policyUseCase.ListMatchesByUserID(ctx, ID)
		oldMap := make(map[domain.ID]domain.Match)
		for _, m := range oldMatches {
			oldMap[m.PolicyID] = m.Status
		}

		updatedMatches, err := u.policyUseCase.UpdatePolicyMatchesForUser(ctx, ID)
		if err == nil && u.notificationUseCase != nil {
			var changedCount int
			for _, m := range updatedMatches {
				if oldStatus, exists := oldMap[m.PolicyID]; exists && oldStatus != m.Status {
					changedCount++
				}
			}
			if changedCount > 0 {
				_, _ = u.notificationUseCase.Create(ctx, &domain.Notification{
					ReceiverID: *ID,
					Type:       domain.NotificationTypeSTAR,
					Message:    "사용자 정보 변경에 따라 정책 적합도 정보가 업데이트되었습니다.",
					Metadata:   map[string]any{"changed_count": changedCount},
				})
			}
		}
	}

	return &domain.Profile{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Age:        user.Age,
		RegionCode: user.RegionCode,
		Gender:     user.Gender,
		IsStudent:  user.IsStudent,
		IsYouth:    user.IsYouth,
		IsPregnant: user.IsPregnant,
		IsBusiness: user.IsBusiness,
		IsForeign:  user.IsForeign,
		Interests:  user.Interests,
		CreatedAt:  user.CreatedAt,
	}, nil
}
