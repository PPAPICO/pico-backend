package usecase

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/janghanul090801/pico-backend/config"
	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepository domain.UserRepository
	policyUseCase  domain.PolicyUseCase
	contextTimeout time.Duration
}

func NewAuthUseCase(userRepository domain.UserRepository, policyUseCase domain.PolicyUseCase, timeout time.Duration) domain.AuthUseCase {
	return &authUseCase{
		userRepository: userRepository,
		policyUseCase:  policyUseCase,
		contextTimeout: timeout,
	}
}

func (u *authUseCase) Register(c context.Context, name, email, password string, age, regionCode int, gender domain.Gender, isStudent, isYouth, isPregnant, isBusiness, isDisabled, isForeign bool, interest []domain.Interest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	_, err := u.userRepository.FindByEmail(ctx, email)
	if err == nil {
		return nil, domain.NewBadRequestError(errors.New("email already taken"))
	}

	encrypted, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, domain.NewBadRequestError(err)
	}

	user, err := u.userRepository.Create(ctx, &domain.User{
		Name:       name,
		Email:      email,
		Password:   string(encrypted),
		Age:        age,
		RegionCode: regionCode,
		Gender:     gender,
		IsStudent:  isStudent,
		IsYouth:    isYouth,
		IsPregnant: isPregnant,
		IsBusiness: isBusiness,
		IsForeign:  isForeign,
		IsDisabled: isDisabled,
		Interests:  interest,
	})
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	policies, err := u.policyUseCase.List(ctx)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	_, err = u.policyUseCase.SavePolicyMatches(ctx, user, policies)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return user, nil
}

func (u *authUseCase) Login(c context.Context, email, password string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.NewUnauthorizedError(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, domain.NewBadRequestError(errors.New("invalid credentials"))
	}

	return user, nil
}

func (u *authUseCase) CreateAccessAndRefreshToken(c context.Context, user *domain.User) (string, string, error) {
	access, err := token.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
	if err != nil {
		return "", "", domain.NewInternalServerError(err)
	}

	refresh, err := token.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
	if err != nil {
		return "", "", domain.NewInternalServerError(err)
	}

	return access, refresh, nil
}

func (u *authUseCase) ExtractUserFromRefreshToken(c context.Context, requestToken string) (*domain.User, error) {
	id, err := token.ExtractIDFromToken(requestToken, config.E.RefreshTokenSecret)
	if err != nil {
		return nil, domain.Error{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	user, err := u.userRepository.FindByID(c, id)
	if err != nil {
		return nil, domain.NewUnauthorizedError(err)
	}

	return user, nil
}
